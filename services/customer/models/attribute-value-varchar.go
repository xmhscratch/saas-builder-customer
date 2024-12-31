package models

import (
	"gopkg.in/guregu/null.v4"
)

// CustomerAttributeVarchar model
type CustomerAttributeVarchar struct {
	CustomerAttributeIndex
	Value null.String `gorm:"column:value;type:string;" sql:"type:varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin" json:"value"`
}

// TableName specifies table name
func (ctx *CustomerAttributeVarchar) TableName() string {
	return "customer_attribute_varchar"
}
