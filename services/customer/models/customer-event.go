package models

import (
	// sq "github.com/Masterminds/squirrel"
	"encoding/json"

	// "errors"
	// "log"

	// "gorm.io/gorm"
	null "gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"localdomain/customer/core"
)

// CustomerEvent model
type CustomerEvent struct {
	DefaultGormModelWithoutID
	EventName               null.String `gorm:"column:event_name;primary_key;type:string;" sql:"type:varchar(255);not null;primary_key;" json:"eventName"`
	Title                   null.String `gorm:"column:title" sql:"type:varchar(255);not null;" json:"title"`
	Description             null.String `gorm:"column:description" sql:"type:varchar(512);" json:"description"`
	AggregateMethod         null.String `gorm:"column:aggregate_method;not null;default:count;" sql:"type:enum('mean','median','count','sum','first','last','min','max','stddev');not null;" json:"aggregateMethod"`
	SynchronizationInterval null.Int    `gorm:"column:synchronization_interval" sql:"type:int(11);" json:"synchronizationInterval"`
	// Series                  []CustomerEventDataInfo `gorm:"foreignKey:event_name;references:event_name" json:"series"`
	TimestampModel
	CustomerEventData []CustomerEventData `gorm:"foreignKey:event_name;references:event_name;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName specifies table name
func (ctx *CustomerEvent) TableName() string {
	return "customer_events"
}

// CustomerEventFilter comment
type CustomerEventFilter struct{}

// FindAllCustomerEvents comment
func FindAllCustomerEvents(db *gorm.DB, filter *CustomerEventFilter, opts *PaginationOptions) ([]*CustomerEvent, error) {
	limit := opts.GetLimit()
	offset := opts.GetOffset()

	var builder *gorm.DB = db.
		Model(CustomerEvent{}).
		Select("*")

	rows, err := builder.
		// Debug().
		Order("event_name").
		Limit(limit).
		Offset(offset).
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customerEvents := make([]*CustomerEvent, 0)

	for rows.Next() {
		var customerEvent CustomerEvent
		db.ScanRows(rows, &customerEvent)

		customerEvents = append(customerEvents, &customerEvent)
	}

	return customerEvents, err
}

// CountAllCustomerEvents comment
func CountAllCustomerEvents(db *gorm.DB, filter *CustomerEventFilter) int64 {
	var total int64 = 0

	db.
		Model(CustomerEvent{}).
		Count(&total)

	return total
}

// FindOneCustomerEvent comment
func FindOneCustomerEvent(db *gorm.DB, eventName string) (*CustomerEvent, error) {
	var (
		err           error
		customerEvent CustomerEvent
	)

	db.
		// Debug().
		Model(CustomerEvent{}).
		Where("event_name = ?", eventName).
		First(&customerEvent)

	return &customerEvent, err
}

// Register comment
func (ctx *CustomerEvent) Register(db *gorm.DB) (err error) {
	return ctx.Transact(db, func(tx *gorm.DB) error {
		return tx.
			Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "event_name"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"updated_at": gorm.Expr("CURRENT_TIMESTAMP()"),
				}),
			}).
			FirstOrCreate(&ctx).
			Error
	})
}

// Update comment
func (ctx *CustomerEvent) Update(db *gorm.DB, infoValues map[string]interface{}) (err error) {
	return ctx.Transact(db, func(tx *gorm.DB) error {
		allowedAttributes := []string{
			"title",
			"description",
			"synchronizationInterval",
		}
		infoValues = core.FilterMap(infoValues, func(v interface{}, k string) bool {
			return core.IncludeInStringSlice(k, allowedAttributes)
		})
		infoValues, err = ConvertMapToKeySnakeCase(infoValues)
		if err != nil {
			return err
		}

		return db.
			Model(&ctx).
			Where(ctx).
			Updates(infoValues).
			Error
	})
}

// Delete comment
func (ctx *CustomerEvent) Delete(db *gorm.DB) (err error) {
	return ctx.Transact(db, func(tx *gorm.DB) error {
		var err error

		findResult := tx.
			Where(ctx).
			FirstOrInit(&ctx)

		if findResult.Error != nil {
			return err
		}

		if findResult.RowsAffected == 0 {
			return err
		}

		err = db.
			Delete(&ctx).
			Error

		return err
	})
}

// FetchDetailInfo comment
func (ctx *CustomerEvent) FetchDetailInfo(db *gorm.DB) (map[string]interface{}, error) {
	var details map[string]interface{}

	b, err := json.Marshal(*ctx)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &details)
	if err != nil {
		return nil, err
	}
	return details, err
}
