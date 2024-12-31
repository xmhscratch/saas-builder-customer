package models

import (
	// "database/sql"
	null "gopkg.in/guregu/null.v4"
)

// CustomerAttributeBlob model
type CustomerAttributeBlob struct {
	CustomerAttributeIndex
	Value []byte      `gorm:"column:value;type:bytes;" sql:"type:mediumblob" json:"value"`
	Name  null.String `gorm:"column:name;type:string;" sql:"type:varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci" json:"name"`
	Type  null.String `gorm:"column:type;type:string;" sql:"type:varchar(100)" json:"type"`
}

// TableName specifies table name
func (ctx *CustomerAttributeBlob) TableName() string {
	return "customer_attribute_blob"
}
