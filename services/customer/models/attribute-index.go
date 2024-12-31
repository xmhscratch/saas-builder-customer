package models

import (
	"gopkg.in/guregu/null.v4"
)

// CustomerAttributeIndex struct
type CustomerAttributeIndex struct {
	CustomerID  null.String `gorm:"column:customer_id;not null;primary_key" sql:"type:varchar(27)" json:"customerId"`
	AttributeID null.String `gorm:"column:attribute_id;not null;primary_key" sql:"type:varchar(27)" json:"attributeId"`
}
