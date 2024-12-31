package models

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	// "log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"localdomain/customer/core"

	qs "github.com/derekstavis/go-qs"
	amqp "github.com/rabbitmq/amqp091-go"
	ksuid "github.com/segmentio/ksuid"
	excelize "github.com/xuri/excelize/v2"

	// null "gopkg.in/guregu/null.v4"
	gorm "gorm.io/gorm"
)

// CustomerExporter comment ...
type CustomerExporter struct {
	SessionID      string
	OrganizationID string
	Headers        map[string]string
	Limit          int64
	Config         *core.Config
	attrFilter     *AttributeFilter
	// db             *gorm.DB
	mqConnection *amqp.Connection
	mqChannel    *amqp.Channel
}

// NewCustomerExporter comment
func NewCustomerExporter(cfg *core.Config, organizationID string, db *gorm.DB) (*CustomerExporter, error) {
	var (
		ctx *CustomerExporter
		err error
	)

	newSessionID := ksuid.New()
	ctx = &CustomerExporter{
		SessionID:      newSessionID.String(),
		OrganizationID: organizationID,
		Config:         cfg,
		Limit:          1000,
	}
	err = ctx.Initialize(db)
	return ctx, err
}

// Initialize comment
func (ctx *CustomerExporter) Initialize(db *gorm.DB) error {
	// Connects opens an AMQP connection from the credentials in the URL.
	mqConnection, err := amqp.DialConfig(
		ctx.Config.AMQPConnectionString,
		amqp.Config{
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, 30*time.Second)
			},
			Properties: amqp.Table{
				"product":  os.Getenv("APP_NAME"),
				"platform": "customer-exporter.publisher",
				"version":  "1.0.0",
			},
		},
	)
	if err != nil {
		return err
	}

	mqChannel, err := mqConnection.Channel()
	if err != nil {
		return err
	}

	err = mqChannel.ExchangeDeclare(
		exchange,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	_, err = mqChannel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	err = mqChannel.QueueBind(
		queueName,
		bindingKey,
		exchange,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	ctx.mqChannel = mqChannel
	ctx.mqConnection = mqConnection

	ctx.attrFilter = &AttributeFilter{
		EntityTypes: map[string]string{
			"0": "varchar",
			"1": "text",
			"2": "boolean",
			"3": "datetime",
			"4": "decimal",
			"5": "int",
		},
		CodeNames: map[string]string{},
		SortBy: map[string]string{
			"sort_order": "DESC",
		},
	}

	headers, err := ctx.ExportHeaders(db)
	if err != nil {
		return err
	}
	ctx.Headers = headers

	_, err = ctx.ExportPartialInfos(db, 0)
	return err
}

// ExportHeaders comment
func (ctx *CustomerExporter) ExportHeaders(db *gorm.DB) (map[string]string, error) {
	var (
		err     error
		results map[string]string
	)

	results = map[string]string{
		"A": "id",
		"B": "emailAddress",
		"C": "shopifyCustomerId",
		"D": "createdAt",
	}
	offsetColIndex := len(results) + 1

	attrs, err := FindAllAttributes(db, ctx.attrFilter, nil)
	if err != nil {
		return nil, err
	}
	for index, attr := range attrs {
		if attr.CodeName.IsZero() {
			continue
		}
		codeName := attr.CodeName.ValueOrZero()
		colName, err := excelize.ColumnNumberToName(index + offsetColIndex)
		if err != nil {
			continue
		}
		results[colName] = codeName
	}
	offsetColIndex = len(results) + 1

	customerEvents, err := FindAllCustomerEvents(db, &CustomerEventFilter{}, &PaginationOptions{Size: 0, Number: 0})
	if err != nil {
		return nil, err
	}
	for index, customerEvent := range customerEvents {
		evtName := customerEvent.EventName.ValueOrZero()
		colName, err := excelize.ColumnNumberToName(index + offsetColIndex)
		if err != nil {
			continue
		}
		results[colName] = evtName
	}

	return results, err
}

// ExportPartialInfos comment
func (ctx *CustomerExporter) ExportPartialInfos(db *gorm.DB, rowOffset int64) (int64, error) {
	var (
		err           error
		numRow        int64 = 0
		lastRowOffset int64 = rowOffset
		customerIDs   []string
	)

	sqlStr := `
SELECT
	id AS customerID
FROM
	customers
LIMIT ?
OFFSET ?;
		`
	rows, err := db.
		// Debug().
		Raw(sqlStr, ctx.Limit, rowOffset).
		Rows()
	if err != nil {
		return numRow, err
	}
	defer rows.Close()

	for rows.Next() {
		var customerID string
		if err = rows.Scan(&customerID); err != nil {
			return numRow, err
		}
		customerIDs = append(customerIDs, customerID)

		numRow++
		rowOffset++
	}

	if numRow == 0 {
		return numRow, err
	}

	if err := ctx.SendDumperRecord(db, lastRowOffset, customerIDs); err != nil {
		return numRow, err
	}

	if _, err = ctx.ExportPartialInfos(db, rowOffset); err != nil {
		return numRow, err
	}
	return numRow, err
}

// SendDumperRecord comment
func (ctx *CustomerExporter) SendDumperRecord(db *gorm.DB, rowOffset int64, customerIDs []string) error {
	var err error

	if err = ctx.ProgressAddJobs(customerIDs); err != nil {
		return err
	}

	for index, customerID := range customerIDs {
		errChan := make(chan error)
		customerIndex := int64(index) + rowOffset

		go func(customerIndex int64, customerID string, errChan chan error) {
			var (
				customerInfo          map[string]interface{}
				customerEventDataInfo map[string]float64
				msgBody               []byte
			)
			customerEventDataInfo = make(map[string]float64)
			customer, err := FindOneCustomer(db, customerID)
			if err != nil {
				errChan <- err
				return
			}
			if customer.ID.IsZero() {
				errChan <- errors.New("customer not found")
				return
			}

			customerInfo, err = customer.FetchDetailInfo(db, ctx.attrFilter)
			if err != nil {
				errChan <- err
				return
			}

			customerEventDataArray, err := FindAllCustomerEventData(db, customerID)
			if err != nil {
				errChan <- err
				return
			}
			for _, customerEventData := range customerEventDataArray {
				if customerEventData.EventName.IsZero() {
					continue
				}
				eventName := customerEventData.EventName.ValueOrZero()
				eventValue := customerEventData.EventValue.ValueOrZero()
				customerEventDataInfo[string(eventName)] = float64(eventValue)
			}

			if msgBody, err = json.Marshal(
				map[string]interface{}{
					"jobId":          customerID,
					"sessionId":      ctx.SessionID,
					"organizationId": ctx.OrganizationID,
					"methodName":     "export",
					"index":          customerIndex,
					"header":         ctx.Headers,
					"info":           customerInfo,
					"event":          customerEventDataInfo,
					"deleteInfos":    []string{},
					"deleteEvents":   []string{},
				},
			); err != nil {
				errChan <- err
				return
			}

			// log.Println(string(msgBody[:]))

			errChan <- ctx.mqChannel.Publish(
				exchange,   // exchange
				routingKey, // routing key
				false,      // mandatory
				false,      // immediate
				amqp.Publishing{
					Headers:         amqp.Table{},
					ContentType:     "application/json",
					ContentEncoding: "",
					Timestamp:       time.Now(),
					Body:            msgBody,
					DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
					Priority:        0,              // 0-9
				},
			)
		}(customerIndex, customerID, errChan)

		err = <-errChan
		close(errChan)
	}

	return err
}

// ProgressAddJobs comment
func (ctx *CustomerExporter) ProgressAddJobs(jobIDs []string) error {
	var (
		err       error
		scriptURL *url.URL
		postQs    string
	)

	respBytes := make(chan []byte, 64)
	go func(byts chan []byte) error {
		scriptURLString := core.BuildString(
			"http://",
			ctx.Config.ClusterHostNames.API,
			"/progress/customer-export/",
			ctx.SessionID,
		)
		if scriptURL, err = url.Parse(scriptURLString); err != nil {
			// log.Println(err)
			byts <- nil
			return err
		}
		jobList := make([]interface{}, 0)
		for i := 0; i < len(jobIDs); i++ {
			jobList = append(jobList, jobIDs[i])
		}
		if postQs, err = qs.Marshal(map[string]interface{}{
			"jobIds": jobList,
		}); err != nil {
			// log.Println(err)
			byts <- nil
			return err
		}
		req, err := http.NewRequest(
			"POST",
			scriptURL.String(),
			strings.NewReader(postQs),
		)
		if err != nil {
			// log.Println(err)
			byts <- nil
			return err
		}
		req.Header.Set("content-type", "application/x-www-form-urlencoded")
		req.Header.Set("x-organization-id", ctx.OrganizationID)
		res, err := HttpClient.Do(req)
		if err != nil {
			// log.Println(err)
			byts <- nil
			return err
		}
		defer res.Body.Close()

		if resBody, err := ioutil.ReadAll(res.Body); err != nil {
			// log.Println(err)
			byts <- nil
			return err
		} else {
			byts <- resBody
		}
		return nil
	}(respBytes)

	<-respBytes
	close(respBytes)

	return err
}

// Dispose comment
func (ctx *CustomerExporter) Dispose() {
	if ctx.mqChannel != nil {
		ctx.mqChannel.Close()
	}
	if ctx.mqConnection != nil {
		ctx.mqConnection.Close()
	}
}
