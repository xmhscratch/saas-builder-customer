package models

import (
	null "gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

// EnumerationValue model
type EnumerationValue struct {
	DefaultGormModelWithoutID
	EnumID null.String `gorm:"column:enum_id;type:string;" sql:"type:varchar(27);not null" json:"enumId"`
	Value  null.Int    `gorm:"column:value;type:int;" sql:"type:int(11);not null" json:"value"`
	Label  null.String `gorm:"column:label;type:string;" sql:"type:varchar(255);not null" json:"label"`
	Data   null.String `gorm:"column:data;type:string;" sql:"type:tinytext" json:"data"`
}

// TableName specifies table name
func (ctx *EnumerationValue) TableName() string {
	return "enumeration_values"
}

// Register comment
func (ctx *EnumerationValue) Register(db *gorm.DB) (err error) {
	return ctx.Transact(db, func(tx *gorm.DB) error {
		return tx.Create(&ctx).Error
	})
}
