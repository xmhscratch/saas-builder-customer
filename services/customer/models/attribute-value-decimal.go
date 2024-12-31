package models

import (
	"gopkg.in/guregu/null.v4"
)

// CustomerAttributeDecimal model
type CustomerAttributeDecimal struct {
	CustomerAttributeIndex
	Value null.Float `gorm:"column:value;type:float;" sql:"type:decimal(12,4)" json:"value"`
}

// TableName specifies table name
func (ctx *CustomerAttributeDecimal) TableName() string {
	return "customer_attribute_decimal"
}
