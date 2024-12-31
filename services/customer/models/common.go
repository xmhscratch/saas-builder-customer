package models

import (
	"net/http"
	"net/http/cookiejar"

	null "gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

// SERIES_DATABASE_NAME comment
const SERIES_DATABASE_NAME = "test_db"

var (
	HttpClientJar, _ = cookiejar.New(nil)
	HttpClient       = &http.Client{Jar: HttpClientJar}
)

var (
	exchange     = "datadump" // Durable, non-auto-deleted AMQP exchange name
	exchangeType = "topic"    // Exchange type - direct|fanout|topic|x-custom
	queueName    = "datadump" // Ephemeral AMQP queue name
	bindingKey   = "datadump" // AMQP binding key
	routingKey   = "datadump" // AMQP binding key
	// consumerTag  = "datadump" // AMQP consumer tag (should not be blank)
)

// TimestampModel ...
type TimestampModel struct {
	CreatedAt null.Time `gorm:"<-:create;column:created_at;not null;default:current_timestamp();type:time;" sql:"type:datetime;" json:"createdAt"`
	UpdatedAt null.Time `gorm:"column:updated_at;type:time;" sql:"type:datetime" json:"updatedAt"`
	DeletedAt null.Time `gorm:"column:deleted_at;type:time;" sql:"type:datetime" json:"deletedAt"`
}

// DefaultGormModelWithoutID ...
type DefaultGormModelWithoutID struct{}

// DefaultGormModel ...
type DefaultGormModel struct {
	DefaultGormModelWithoutID `gorm:"-" json:"-"`
	ID                        null.String `gorm:"<-:create;column:id;primary_key;type:string;" sql:"type:char(27)" json:"id"`
}

// Transact comment
func (ctx *DefaultGormModelWithoutID) Transact(db *gorm.DB, txFunc func(*gorm.DB) error) (err error) {
	tx := db.Begin()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit().Error // err is nil; if Commit returns error update err
		}
	}()

	err = txFunc(tx)
	return err
}

// // FetchDetailInfo comment
// func (ctx *DefaultGormModelWithoutID) FetchDetailInfo(db *gorm.DB) (map[string]interface{}, error) {
// 	var details map[string]interface{}

// 	b, err := json.Marshal(*ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = json.Unmarshal(b, &details)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return details, err
// }

// Transact comment
func (ctx *DefaultGormModel) Transact(db *gorm.DB, txFunc func(*gorm.DB) error) (err error) {
	return ctx.DefaultGormModelWithoutID.Transact(db, txFunc)
}

// // FetchDetailInfo comment
// func (ctx *DefaultGormModel) FetchDetailInfo(db *gorm.DB) (map[string]interface{}, error) {
// 	return ctx.DefaultGormModelWithoutID.FetchDetailInfo(db)
// }

// CustomerDataRecord comment
type CustomerDataRecord struct {
	JobID          string                 `json:"jobId"`
	SessionID      string                 `json:"sessionId"`
	OrganizationID string                 `json:"organizationId"`
	MethodName     string                 `json:"methodName"`
	Index          int                    `json:"index"`
	Header         map[string]string      `json:"header"`
	Info           map[string]interface{} `json:"info"`
	Event          map[string]float64     `json:"event"`
	DeleteInfos    []string               `json:"deleteInfos"`
	DeleteEvents   []string               `json:"deleteEvents"`
}
