package models

import (
	"gopkg.in/guregu/null.v4"
)

// CustomerAttributeBoolean model
type CustomerAttributeBoolean struct {
	CustomerAttributeIndex
	Value null.Bool `gorm:"column:value;type:bool;" sql:"type:tinyint(1)" json:"value"`
}

// TableName specifies table name
func (ctx *CustomerAttributeBoolean) TableName() string {
	return "customer_attribute_boolean"
}
