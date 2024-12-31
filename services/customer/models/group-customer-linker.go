package models

import (
	"gopkg.in/guregu/null.v4"
)

// GroupCustomerLinker struct
type GroupCustomerLinker struct {
	GroupID    null.String `gorm:"column:group_id;primary_key;not null;type:string;" sql:"type:char(27) character set ascii collate ascii_general_ci;" json:"groupId"`
	CustomerID null.String `gorm:"column:customer_id;primary_key;not null;type:string;" sql:"type:char(27) character set ascii collate ascii_general_ci;" json:"customerId"`
	Group      Group       `gorm:"foreignKey:group_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Customer   Customer    `gorm:"foreignKey:customer_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName specifies table name
func (ctx *GroupCustomerLinker) TableName() string {
	return "group_customer_linker"
}
