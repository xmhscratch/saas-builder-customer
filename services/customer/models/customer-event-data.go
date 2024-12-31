package models

import (
	// sq "github.com/Masterminds/squirrel"

	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	null "gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	// this is important because of the bug in go mod
	_ "github.com/influxdata/influxdb1-client"
	influxdb1 "github.com/influxdata/influxdb1-client/v2"

	"localdomain/customer/core"
)

// CustomerEventData model
type CustomerEventData struct {
	DefaultGormModelWithoutID
	EventName       null.String   `gorm:"column:event_name;primary_key;type:string;" sql:"type:varchar(255);not null;primary_key;" json:"eventName"`
	CustomerID      null.String   `gorm:"column:customer_id" sql:"type:char(27);not null;primary_key;" json:"customerId"`
	SnapshotDate    null.Time     `gorm:"column:snapshot_date" sql:"type:date" json:"snapshotDate,omitempty"`
	AggregateMethod null.String   `gorm:"column:aggregate_method;not null;default:count;" sql:"type:enum('mean','median','count','sum','first','last','min','max','stddev');not null;" json:"aggregateMethod"`
	CreatedAt       null.Time     `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP();" sql:"type:datetime;not null;" json:"createdAt"`
	UpdatedAt       null.Time     `gorm:"column:updated_at;" sql:"type:datetime;" json:"updatedAt,omitempty"`
	EventOccurrence null.Int      `gorm:"column:event_occurrence" sql:"type:int(11);" json:"eventOccurrence"`
	EventValue      null.Float    `gorm:"column:event_value" sql:"type:decimal(12,4);" json:"eventValue"`
	CustomerEvent   CustomerEvent `gorm:"foreignKey:event_name;references:event_name;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"eventInfo"`
	Customer        Customer      `gorm:"foreignKey:customer_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName specifies table name
func (ctx *CustomerEventData) TableName() string {
	return "customer_event_series"
}

// CustomerEventDataInfo model
type CustomerEventDataInfo struct {
	DefaultGormModelWithoutID
	EventName       null.String `gorm:"column:event_name;primary_key;type:string;" sql:"type:varchar(255);not null;primary_key;" json:"eventName"`
	AggregateMethod null.String `gorm:"column:aggregate_method;not null;default:count;" sql:"type:enum('mean','median','count','sum','first','last','min','max','stddev');not null;" json:"aggregateMethod"`
	EventOccurrence null.Int    `gorm:"column:event_occurrence" sql:"type:int(11);" json:"eventOccurrence"`
	EventValue      null.Float  `gorm:"column:event_value" sql:"type:decimal(12,4);" json:"eventValue"`
}

// AggregationEntry comment
type AggregationEntry struct {
	AggregateMethod     string
	PrevEventOccurrence int64
	PrevEventValue      float64
	NextEventOccurrence int64
	NextEventValue      float64
}

// FindAllCustomerEventData comment
func FindAllCustomerEventData(db *gorm.DB, customerID string) ([]*CustomerEventDataInfo, error) {
	var (
		err     error
		builder *gorm.DB
	)
	list := make([]*CustomerEventDataInfo, 0)
	entries := make(map[string]*CustomerEventDataInfo)

	builder = db.
		// Debug().
		Model(CustomerEvent{}).
		Select(core.BuildString(
			"`customer_events`.`event_name`,",
			"IFNULL(`customer_event_series`.`customer_id`, ?) AS `customer_id`,",
			"`customer_events`.`aggregate_method`,",
			// "`customer_event_series`.`snapshot_date`,",
			// "`customer_event_series`.`created_at`,",
			"`customer_event_series`.`updated_at`,",
			"`customer_event_series`.`event_occurrence`,",
			"`customer_event_series`.`event_value`",
		), customerID).
		Joins("LEFT JOIN `customer_event_series` ON `customer_event_series`.`event_name`=`customer_events`.`event_name` AND `customer_event_series`.`customer_id` = ?", customerID).
		Order("`customer_event_series`.`event_name`").
		Order("`customer_event_series`.`updated_at` ASC").
		Where("IFNULL(`customer_event_series`.`customer_id`, ?)", customerID)

	rows, err := builder.
		Rows()

	if err != nil {
		return list, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			eventData          CustomerEventDataInfo
			newEventOccurrence int64   = 1
			newEventValue      float64 = 0.0000
			aggregationEntry   *AggregationEntry
		)
		db.ScanRows(rows, &eventData)

		eventName := eventData.EventName.ValueOrZero()
		aggregateMethod := eventData.AggregateMethod.ValueOrZero()

		if entries[eventName] != nil {
			existEventData := entries[eventName]

			// log.Println(
			// 	eventName,
			// 	aggregateMethod,
			// 	existEventData.EventOccurrence.ValueOrZero(),
			// 	existEventData.EventValue.ValueOrZero(),
			// 	eventData.EventOccurrence.ValueOrZero(),
			// 	eventData.EventValue.ValueOrZero(),
			// )

			aggregationEntry = &AggregationEntry{
				AggregateMethod:     aggregateMethod,
				PrevEventOccurrence: existEventData.EventOccurrence.ValueOrZero(),
				PrevEventValue:      existEventData.EventValue.ValueOrZero(),
				NextEventOccurrence: eventData.EventOccurrence.ValueOrZero(),
				NextEventValue:      eventData.EventValue.ValueOrZero(),
			}
			newEventOccurrence, newEventValue = aggregationEntry.Aggregate(3)

			// log.Println(
			// 	eventName,
			// 	newEventOccurrence,
			// 	newEventValue,
			// )
		} else {
			aggregationEntry = &AggregationEntry{
				AggregateMethod:     aggregateMethod,
				PrevEventOccurrence: int64(1),
				PrevEventValue:      float64(0.0000),
				NextEventOccurrence: eventData.EventOccurrence.ValueOrZero(),
				NextEventValue:      eventData.EventValue.ValueOrZero(),
			}
			newEventOccurrence, newEventValue = aggregationEntry.Aggregate(1)
		}

		// log.Println(eventName, newEventOccurrence, newEventValue)

		entries[eventName] = &CustomerEventDataInfo{
			EventName:       null.StringFrom(eventName),
			AggregateMethod: null.StringFrom(aggregateMethod),
			EventOccurrence: null.IntFrom(newEventOccurrence),
			EventValue:      null.FloatFrom(newEventValue),
		}
	}

	// jsonString, _ := json.Marshal(entries)
	// log.Println(string(jsonString[:]))

	for _, entryItem := range entries {
		list = append(list, entryItem)
	}
	return list, err
}

// Sync comment
func (ctx *CustomerEventData) Sync(db *gorm.DB, cfg *core.Config) (err error) {
	return ctx.Transact(db, func(tx *gorm.DB) error {
		var (
			err               error
			aggregateMethod   string = "count"
			ctxEventInfo      CustomerEvent
			aggregateFuncName string = "count"
		)

		ctxEventName := ctx.EventName.ValueOrZero()
		if ctx.EventName.IsZero() {
			err = fmt.Errorf("event name is not provided")
			return err
		}

		err = db.
			Model(CustomerEvent{}).
			Where("event_name = ?", ctxEventName).
			First(&ctxEventInfo).
			Error

		if err != nil {
			return err
		}

		if ctxEventInfo.AggregateMethod.IsZero() {
			err = fmt.Errorf("invalid event aggregation method")
			return err
		}

		aggregateMethod = ctxEventInfo.AggregateMethod.ValueOrZero()
		// log.Println(aggregateMethod)

		clientTimeout, err := time.ParseDuration("30s")
		if err != nil {
			fmt.Println("Error parsing client timeout duration: ", err.Error())
			return err
		}
		influxCfg := influxdb1.HTTPConfig{
			Addr:               cfg.InfluxConnectionString,
			Timeout:            clientTimeout,
			InsecureSkipVerify: false,
			// Username string
			// Password string
		}
		c, err := influxdb1.NewHTTPClient(influxCfg)
		if err != nil {
			fmt.Println("error creating InfluxDB Client: ", err.Error())
			return err
		}
		defer c.Close()

		// if aggregateMethod != "" {
		// 	aggregateFuncName = aggregateMethod
		// }
		if aggregateMethod != "" {
			aggregateFuncName = aggregateMethod
		}

		queryString := fmt.Sprintf(
			`SELECT last("createdAt") AS createdAt, %s("value") AS "eventValue" FROM "subscribe" WHERE time > now() - 1d AND "eventName"='%s' GROUP BY time(1d), "eventName", "customerId" ORDER BY time ASC`,
			aggregateFuncName,
			ctxEventName,
		)

		// log.Println(queryString)

		q := influxdb1.NewQuery(queryString, SERIES_DATABASE_NAME, "default")
		resp, err := c.Query(q)
		if err != nil {
			fmt.Println("Error querying: ", err.Error())
			return err
		}
		if resp.Error() != nil {
			fmt.Println("Error querying: ", resp.Error())
			return err
		}

		for _, result := range resp.Results {
			// log.Println(result)

			for _, row := range result.Series {
				eventValueColumnIndex := core.IndexInStringSlice("eventValue", row.Columns)
				createdAtColumnIndex := core.IndexInStringSlice("createdAt", row.Columns)

				eventNameTagValue := row.Tags["eventName"]
				customerIDTagValue := row.Tags["customerId"]

				for _, rowRecord := range row.Values {
					// log.Println(rowRecord)

					if rowRecord[eventValueColumnIndex] == nil {
						continue
					}

					if rowRecord[createdAtColumnIndex] == nil {
						continue
					}

					var (
						timeValue           time.Time
						nextEventValue      float64
						nextEventOccurrence int64   = 1
						newEventValue       float64 = 0.0000
						newEventOccurrence  int64   = 1
					)

					if ctxEventName != eventNameTagValue {
						continue
					}

					if et, err := (rowRecord[createdAtColumnIndex].(json.Number)).Float64(); err == nil {
						// to milliseconds
						ts := int64(math.Round(et / 1e3))
						timeValue = time.Unix(ts, 0)

						if err != nil {
							log.Println(err)
						}
					} else {
						log.Println(err)
						// continue
					}

					if ev, err := (rowRecord[eventValueColumnIndex].(json.Number)).Float64(); err == nil {
						nextEventValue = float64(ev)
					} else {
						log.Println(err)
					}

					var eventData CustomerEventData = CustomerEventData{
						EventName:       null.StringFrom(eventNameTagValue),
						CustomerID:      null.StringFrom(customerIDTagValue),
						AggregateMethod: null.StringFrom(aggregateMethod),
						EventOccurrence: null.IntFrom(0),
						EventValue:      null.FloatFrom(0.0000),
					}

					queryResult := db.
						// Debug().
						Model(CustomerEventData{}).
						// Clauses(clause.OnConflict{
						// 	Columns:   []clause.Column{{Name: "id"}},
						// 	DoUpdates: clause.Assignments(map[string]interface{}{
						// 		"updated_at": "CURRENT_TIMESTAMP()",
						// 	}),
						// }).
						Table("customer_event_series").
						Where("event_name = ?", eventNameTagValue).
						Where("customer_id = ?", customerIDTagValue).
						Where(gorm.Expr("DATEDIFF(IFNULL(snapshot_date, CURRENT_DATE), CURRENT_DATE)=0")).
						FirstOrInit(&eventData, CustomerEventData{
							// SnapshotDate: null.TimeFrom(time.Now()),
							// UpdatedAt: eventData.UpdatedAt.ValueOrZero()
							// PrevEventValue: eventData.EventValue.ValueOrZero()
							// PrevEventOccurrence: eventData.EventOccurrence.ValueOrZero()
						})

					if queryResult.Error != nil {
						fmt.Println("Error querying: ", queryResult.Error)
						// return err
					}

					var (
						updatedAt           time.Time
						prevEventValue      float64 = 0.0000
						prevEventOccurrence int64   = 0
						isWithinDay         bool    = updatedAt.Sub(timeValue).Minutes() <= (60 * time.Minute * 24).Minutes()
					)

					// log.Println(queryResult.RowsAffected)

					if queryResult.RowsAffected > 0 {
						updatedAt = eventData.UpdatedAt.ValueOrZero()
						prevEventValue = eventData.EventValue.ValueOrZero()
						prevEventOccurrence = eventData.EventOccurrence.ValueOrZero()
					}

					// log.Println(timeValue, updatedAt, timeValue.Before(updatedAt))
					if timeValue.Before(updatedAt) {
						continue
					}

					aggregationEntry := &AggregationEntry{
						AggregateMethod:     aggregateMethod,
						PrevEventOccurrence: prevEventOccurrence,
						PrevEventValue:      prevEventValue,
						NextEventOccurrence: nextEventOccurrence,
						NextEventValue:      nextEventValue,
					}

					switch true {
					case core.IncludeInStringSlice(aggregateMethod, []string{"count", "mean", "sum"}):
						{
							if isWithinDay {
								newEventOccurrence, newEventValue = aggregationEntry.Aggregate(2)
								break
							}
							newEventOccurrence, newEventValue = aggregationEntry.Aggregate(3)
							break
						}
					default:
						{
							newEventOccurrence, newEventValue = aggregationEntry.Aggregate(3)
							break
						}
					}
					// log.Println(isWithinDay, newEventOccurrence, newEventValue)

					eventData.SnapshotDate = null.TimeFrom(time.Now())
					eventData.EventValue = null.FloatFrom(newEventValue)
					eventData.EventOccurrence = null.IntFrom(newEventOccurrence)
					// eventData.UpdatedAt = null.TimeFrom(timeValue)

					_ = db.
						// Debug().
						Clauses(clause.OnConflict{
							Columns: []clause.Column{
								{Name: "event_name"},
								{Name: "customer_id"},
								{Name: "snapshot_date"},
							},
							DoUpdates: clause.AssignmentColumns([]string{
								"updated_at",
								"event_occurrence",
								"event_value",
							}),
						}).
						Create(&eventData)
				}
			}
		}

		return err
	})
}

// newEventValue

// prevEventValue = 3.0277
// prevEventOccurrence = 3

// nextEventValue = 3.3166
// nextEventOccurrence = 4

// // Update comment
// func (ctx *CustomerEventData) Update(db *gorm.DB, infoValues map[string]interface{}) (err error) {
// 	return ctx.Transact(db, func(tx *gorm.DB) error {
// 		var err error

// 		err = tx.
// 			Where(ctx).
// 			FirstOrInit(&ctx).
// 			Error

// 		if err != nil {
// 			return err
// 		}

// 		if tx.NewRecord(ctx) {
// 			return err
// 		}

// 		allowedAttributes := []string{
// 			"title",
// 			"description",
// 			"synchronizationInterval",
// 		}
// 		infoValues = core.FilterMap(infoValues, func(v interface{}, k string) bool {
// 			return core.IncludeInStringSlice(k, allowedAttributes)
// 		})

// 		err = db.
// 			Model(&ctx).
// 			Updates(infoValues).
// 			Error

// 		return err
// 	})
// }

// // Delete comment
// func (ctx *CustomerEventData) Delete(db *gorm.DB) (err error) {
// 	return ctx.Transact(db, func(tx *gorm.DB) error {
// 		return nil
// 	})
// }

// Aggregate ...
func (ctx *AggregationEntry) Aggregate(incrementalMode int) (int64, float64) {
	var (
		newEventOccurrence  int64   = 1
		newEventValue       float64 = 0.0000
		prevEventOccurrence int64   = ctx.PrevEventOccurrence
		prevEventValue      float64 = ctx.PrevEventValue
		nextEventOccurrence int64   = ctx.NextEventOccurrence
		nextEventValue      float64 = ctx.NextEventValue
	)

	switch incrementalMode {
	// force set new value
	case 1:
		{
			newEventValue = nextEventValue
			newEventOccurrence = nextEventOccurrence
			break
		}
	// set new value
	case 2:
		{
			newEventValue = nextEventValue
			newEventOccurrence = prevEventOccurrence + nextEventOccurrence
			break
		}
	// calculate and set new value
	case 3:
		{
			// log.Println(ctx.AggregateMethod, incrementalMode, nextEventValue, prevEventValue)
			switch true {
			case ctx.AggregateMethod == "count":
				{
					newEventValue = nextEventValue + prevEventValue
					break
				}
			case ctx.AggregateMethod == "first":
				{
					newEventValue = nextEventValue
					break
				}
			case ctx.AggregateMethod == "last":
				{
					newEventValue = nextEventValue
					break
				}
			case ctx.AggregateMethod == "min":
				{
					if prevEventValue != 0 {
						newEventValue = math.Min(nextEventValue, prevEventValue)
					} else {
						newEventValue = nextEventValue
					}
					break
				}
			case ctx.AggregateMethod == "max":
				{
					newEventValue = math.Max(nextEventValue, prevEventValue)
					break
				}
			case ctx.AggregateMethod == "mean":
				{
					if prevEventOccurrence != 0 {
						newEventValue = nextEventValue
						break
					}
					newEventValue = ((prevEventValue * float64(nextEventOccurrence)) + (nextEventValue * float64(prevEventOccurrence))) / float64(prevEventOccurrence+nextEventOccurrence)
					break
				}
			case ctx.AggregateMethod == "sum":
				{
					newEventValue = nextEventValue + prevEventValue
					break
				}
			// case ctx.AggregateMethod == "stddev":
			// 	{
			// 		newEventValue = math.Sqrt((math.Pow(prevEventValue, float64(2)) * float64(prevEventOccurrence)) + (math.Pow(nextEventValue, float64(2)) * float64(nextEventOccurrence)))
			// 		break
			// 	}
			default:
				{
					break
				}
			}
			newEventOccurrence = prevEventOccurrence + nextEventOccurrence
			break
		}
	default:
		return ctx.Aggregate(3)
	}

	return newEventOccurrence, newEventValue
}

/*
SET @range = 100;
SELECT
    FLOOR(event_value/@range)*@range AS `from`,
    (CEIL(event_value/@range)+IF(MOD(event_value,@range)=0,1,0))*@range AS `to`,
    CONCAT(FLOOR(event_value/@range)*@range,'-',(CEIL(event_value/@range)+IF(MOD(event_value,@range)=0,1,0))*@range) AS `label`,
    COUNT(customer_id) AS `value`
FROM customer_event_series
WHERE event_name='facebook_subscribed'
GROUP BY `label`
ORDER BY `from` ASC;
*/
