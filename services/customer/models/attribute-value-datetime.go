package models

import (
	"gopkg.in/guregu/null.v4"
)

// CustomerAttributeDateTime model
type CustomerAttributeDateTime struct {
	CustomerAttributeIndex
	Value  null.Time   `gorm:"column:value;type:time;" sql:"type:datetime" json:"value"`
	Value2 interface{} `gorm:"column:value2;type:time;" sql:"type:datetime" json:"value2"`
}

// TableName specifies table name
func (ctx *CustomerAttributeDateTime) TableName() string {
	return "customer_attribute_datetime"
}
