package models

import (
	null "gopkg.in/guregu/null.v4"
)

// EnumerationFilter comment
type EnumerationFilter struct {
	SeachKeyword null.String `gorm:"-"`
}

// TableName specifies table name
func (ctx *EnumerationFilter) TableName() string {
	return "e"
}
