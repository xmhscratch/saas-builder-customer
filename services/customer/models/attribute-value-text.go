package models

import (
	"gopkg.in/guregu/null.v4"
)

// CustomerAttributeText model
type CustomerAttributeText struct {
	CustomerAttributeIndex
	Value null.String `gorm:"column:value;type:string;" sql:"type:text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin" json:"value"`
}

// TableName specifies table name
func (ctx *CustomerAttributeText) TableName() string {
	return "customer_attribute_text"
}
