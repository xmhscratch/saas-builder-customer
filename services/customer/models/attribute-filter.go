package models

import (
	"strings"

	"localdomain/customer/core"

	"gopkg.in/guregu/null.v4"
)

// AttributeFilter comment
type AttributeFilter struct {
	EntityTypes      map[string]string `gorm:"-"`
	CodeNames        map[string]string `gorm:"-"`
	SortBy           map[string]string `gorm:"-"` // ?r=created_at DESC&r=title ASC
	WithAssociations bool              `gorm:"-"`
	// SearchQuery      map[string]string `gorm:"-"` // ?wa=price BETWEEN 30 AND 140&wo=
	IsFilterable     null.Bool `gorm:"column:is_filterable;type:bool;" sql:"type:tinyint(1)" json:"isFilterable"`
	IsVisibleOnFront null.Bool `gorm:"column:is_visible_on_front;type:bool;" sql:"type:tinyint(1)" json:"isVisibleOnFront"`
	IsVisibleInList  null.Bool `gorm:"column:is_visible_in_list;type:bool;" sql:"type:tinyint(1)" json:"isVisibleInList"`
	IsConfigurable   null.Bool `gorm:"column:is_configurable;type:bool;" sql:"type:tinyint(1)" json:"isConfigurable"`
	IsUserDefined    null.Bool `gorm:"column:is_user_defined;type:bool;" sql:"type:tinyint(1)" json:"isUserDefined"`
	IsReadOnly       null.Bool `gorm:"column:is_read_only;type:bool;not null" sql:"type:tinyint(1)" json:"isReadOnly"`
	IsRequired       null.Bool `gorm:"column:is_required;type:bool;" sql:"type:tinyint(1)" json:"isRequired"`
	IsUnique         null.Bool `gorm:"column:is_unique;type:bool;" sql:"type:tinyint(1)" json:"isUnique"`
}

// TableName specifies table name
func (ctx *AttributeFilter) TableName() string {
	return "e"
}

// FilterQueryString comment
func (ctx *AttributeFilter) FilterQueryString() string {
	var (
		results     []string
		codeNames   []string
		entityTypes []string
	)

	for _, val := range ctx.CodeNames {
		codeNames = append(codeNames, core.BuildString("\"", val, "\""))
	}

	if len(codeNames) > 0 {
		results = append(results, core.BuildString(
			"`e`.`code_name` IN (",
			strings.Join(codeNames, ","),
			")",
		))
	}

	for _, val := range ctx.EntityTypes {
		entityTypes = append(entityTypes, core.BuildString("\"", val, "\""))
	}

	if len(entityTypes) > 0 {
		results = append(results, core.BuildString(
			"`e`.`entity_type` IN (",
			strings.Join(entityTypes, ","),
			")",
		))
	}

	if !ctx.IsFilterable.IsZero() {
		isFilterable := "0"
		if ctx.IsFilterable.ValueOrZero() {
			isFilterable = "1"
		}
		results = append(results, core.BuildString(
			"`e`.`is_filterable` = (", isFilterable, ")",
		))
	}

	if !ctx.IsVisibleOnFront.IsZero() {
		isVisibleOnFront := "0"
		if ctx.IsVisibleOnFront.ValueOrZero() {
			isVisibleOnFront = "1"
		}
		results = append(results, core.BuildString(
			"`e`.`is_visible_on_front` = (", isVisibleOnFront, ")",
		))
	}

	if !ctx.IsVisibleInList.IsZero() {
		isVisibleInList := "0"
		if ctx.IsVisibleInList.ValueOrZero() {
			isVisibleInList = "1"
		}
		results = append(results, core.BuildString(
			"`e`.`is_visible_in_list` = (", isVisibleInList, ")",
		))
	}

	if !ctx.IsConfigurable.IsZero() {
		isConfigurable := "0"
		if ctx.IsConfigurable.ValueOrZero() {
			isConfigurable = "1"
		}
		results = append(results, core.BuildString(
			"`e`.`is_configurable` = (", isConfigurable, ")",
		))
	}

	if !ctx.IsUserDefined.IsZero() {
		isUserDefined := "0"
		if ctx.IsUserDefined.ValueOrZero() {
			isUserDefined = "1"
		}
		results = append(results, core.BuildString(
			"`e`.`is_user_defined` = (", isUserDefined, ")",
		))
	}

	if !ctx.IsReadOnly.IsZero() {
		isReadOnly := "0"
		if ctx.IsReadOnly.ValueOrZero() {
			isReadOnly = "1"
		}
		results = append(results, core.BuildString(
			"`e`.`is_read_only` = (", isReadOnly, ")",
		))
	}

	if !ctx.IsRequired.IsZero() {
		isRequired := "0"
		if ctx.IsRequired.ValueOrZero() {
			isRequired = "1"
		}
		results = append(results, core.BuildString(
			"`e`.`is_required` = (", isRequired, ")",
		))
	}

	if !ctx.IsUnique.IsZero() {
		isUnique := "0"
		if ctx.IsUnique.ValueOrZero() {
			isUnique = "1"
		}
		results = append(results, core.BuildString(
			"`e`.`is_unique` = (", isUnique, ")",
		))
	}

	return strings.Join(results, " AND ")
}

// SortByQueryString comment
func (ctx *AttributeFilter) SortByQueryString() string {
	var sortByList []string

	if ctx.SortBy == nil {
		return "created_at DESC"
	}

	for byField, dirField := range ctx.SortBy {
		sortByList = append(sortByList, core.BuildString(byField, " ", dirField))
	}
	return strings.Join(sortByList, ", ")
}
