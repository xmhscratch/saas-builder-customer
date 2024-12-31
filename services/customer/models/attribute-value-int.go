package models

import (
	"gopkg.in/guregu/null.v4"
)

// CustomerAttributeInt model
type CustomerAttributeInt struct {
	CustomerAttributeIndex
	Value null.Int `gorm:"column:value;type:int;" sql:"type:int(11)" json:"value"`
}

// TableName specifies table name
func (ctx *CustomerAttributeInt) TableName() string {
	return "customer_attribute_int"
}
