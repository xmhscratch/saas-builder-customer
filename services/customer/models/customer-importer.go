package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"localdomain/customer/core"

	qs "github.com/derekstavis/go-qs"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/segmentio/ksuid"
	excelize "github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// CustomerImporter comment ...
type CustomerImporter struct {
	SessionID      string
	OrganizationID string
	// Headers        map[string]string
	// Limit          int64
	Config       *core.Config
	db           *gorm.DB
	mqConnection *amqp.Connection
	mqChannel    *amqp.Channel
}

// NewCustomerImporter comment
func NewCustomerImporter(cfg *core.Config, organizationID string, db *gorm.DB) (*CustomerImporter, error) {
	var (
		ctx *CustomerImporter
		err error
	)

	newSessionID := ksuid.New()
	ctx = &CustomerImporter{
		SessionID:      newSessionID.String(),
		OrganizationID: organizationID,
		Config:         cfg,
		db:             db,
	}
	err = ctx.Initialize()
	return ctx, err
}

// Initialize comment
func (ctx *CustomerImporter) Initialize() error {
	// Connects opens an AMQP connection from the credentials in the URL.
	mqConnection, err := amqp.DialConfig(
		ctx.Config.AMQPConnectionString,
		amqp.Config{
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, 30*time.Second)
			},
			Properties: amqp.Table{
				"product":  os.Getenv("APP_NAME"),
				"platform": "customer-importer.publisher",
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

	return err
}

// AddFile comment
func (ctx *CustomerImporter) AddFile(driveID string, fileID string) error {
	var (
		err         error
		colNumber   int = 0
		rowNumber   int = 0
		infoHeaders map[string]string
		infoValues  []*CustomerDataRecord
		jobIDs      []string
	)

	filePath := filepath.Join(
		ctx.Config.DataPath, "drives",
		core.BuildString(driveID[0:8], "/", driveID[9:]),
		core.BuildString(fileID[0:8], "/", fileID[9:13], "/", fileID[14:18], "/", fileID[19:23], "/", fileID[24:36], ".bin"),
	)

	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		return err
	}

	fileHnd, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}

	// if fileHnd.GetActiveSheetIndex() == 0 {}

	infoHeaders = map[string]string{}
	infoValues = make([]*CustomerDataRecord, 0)

	cols, err := fileHnd.Cols("Sheet1")
	if err != nil {
		return err
	}

	for cols.Next() {
		rowValues, err := cols.Rows()
		if err != nil {
			log.Println(err)
			continue
		}

		colNumber = colNumber + 1
		colName, err := excelize.ColumnNumberToName(colNumber)
		if err != nil {
			log.Println(err)
			continue
		}

		attrName := rowValues[0]
		infoHeaders[colName] = attrName

		for rowNumber := range rowValues {
			if rowNumber == 0 {
				continue
			}

			if attrName == "emailAddress" {
				jobID := ksuid.New().String()
				jobIDs = append(jobIDs, jobID)

				dataRecord := &CustomerDataRecord{
					SessionID:      ctx.SessionID,
					JobID:          jobID,
					OrganizationID: ctx.OrganizationID,
					MethodName:     "import",
					Index:          rowNumber,
					Header:         map[string]string{},
					Info:           map[string]interface{}{},
					Event:          map[string]float64{},
					DeleteInfos:    []string{},
					DeleteEvents:   []string{},
				}
				infoValues = append(infoValues, dataRecord)
			}
		}
	}

	rows, err := fileHnd.Rows("Sheet1")
	if err != nil {
		return err
	}
	for rows.Next() {
		rowNumber = rowNumber + 1

		err := func() error {
			if rowNumber == 1 {
				// skip row reading
				_, err = rows.Columns()
				return err
			}

			if rowNumber > len(infoValues)+1 {
				return nil
			}

			dataRecord := infoValues[rowNumber-2]
			dataRecord.Header = infoHeaders

			// log.Println(dataRecord)

			rowValues, err := rows.Columns()
			if err != nil {
				log.Println(err)
				return err
			}

			for colNumber, colValue := range rowValues {
				colNumber = int(colNumber) + 1
				colName, err := excelize.ColumnNumberToName(colNumber)
				if err != nil {
					log.Println(err)
					continue
				}

				attrName := infoHeaders[colName]
				if attrName == "id" {
					dataRecord.Info["id"] = colValue
					continue
				}
				if attrName == "emailAddress" {
					dataRecord.Info["emailAddress"] = colValue

					customerID, err := FindCustomerID(ctx.db, colValue)
					if err != nil || customerID == "" {
						// log.Println(err)
						continue
					}

					delete(dataRecord.Info, "id")
				}
				if colValue == "" {
					dataRecord.DeleteInfos = append(dataRecord.DeleteInfos, attrName)
					continue
				}
				isValidAttr := core.IncludeInStringSlice(
					core.GetCaseName(attrName),
					[]string{"LowerFlat", "UpperFlat", "LowerCamel", "UpperCamel"},
				)
				if isValidAttr {
					dataRecord.Info[attrName] = colValue
				} else {
					// log.Println(attrName)
					continue
				}
				// todo: metric
			}

			return nil
		}()
		if err != nil {
			log.Println(err)
			continue
		}
	}

	for _, dataRecord := range infoValues {
		err := ctx.SendDumperRecord(dataRecord)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	if err = ctx.ProgressAddJobs(jobIDs); err != nil {
		return err
	}

	return err
}

// SendDumperRecord comment
func (ctx *CustomerImporter) SendDumperRecord(dataRecord *CustomerDataRecord) error {
	var (
		err     error
		msgBody []byte
	)
	errChan := make(chan error)

	if msgBody, err = json.Marshal(dataRecord); err != nil {
		return err
	}

	go func(msgBody []byte, errChan chan error) {
		err = ctx.mqChannel.Publish(
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
		// log.Println(string(msgBody[:]))
		errChan <- err
	}(msgBody, errChan)

	err = <-errChan
	close(errChan)

	return err
}

// ProgressAddJobs comment
func (ctx *CustomerImporter) ProgressAddJobs(jobIDs []string) error {
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
			"/progress/customer-import/",
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
func (ctx *CustomerImporter) Dispose() {
	if ctx.mqChannel != nil {
		ctx.mqChannel.Close()
	}
	if ctx.mqConnection != nil {
		ctx.mqConnection.Close()
	}
}
