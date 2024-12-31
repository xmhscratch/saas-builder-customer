package models

import (
	"strings"

	"localdomain/customer/core"
)

// GroupFilter comment
type GroupFilter struct {
	CodeNames map[string]string `gorm:"-"`
	RootOnly  bool              `gorm:"-"`
}

// TableName specifies table name
func (ctx *GroupFilter) TableName() string {
	return "e"
}

// GetCodeNameFilterQuery comment
func (ctx *GroupFilter) GetCodeNameFilterQuery() string {
	var codeNames []string
	for _, val := range ctx.CodeNames {
		codeNames = append(codeNames, core.BuildString("\"", val, "\""))
	}

	if len(codeNames) == 0 {
		return ""
	}

	return core.BuildString(
		"`e`.`code_name` IN (",
		strings.Join(codeNames, ","),
		")",
	)
}
