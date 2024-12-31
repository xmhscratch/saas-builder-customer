package models

import (
	"log"
	"net/url"
)

// ?
// s[filter_groups][0][filters][0][field]=group_gear&
// s[filter_groups][0][filters][0][value]=86&
// s[filter_groups][0][filters][0][type]=finset

// SearchCriteriaConditionTypes comment
type SearchCriteriaConditionTypes string

const (
	// Equals Equals.
	Equals SearchCriteriaConditionTypes = "eq"

	// FindInSet A value within a set of values
	FindInSet SearchCriteriaConditionTypes = "finset"

	// GreaterThan Greater than
	GreaterThan SearchCriteriaConditionTypes = "gt"

	// GreaterThanEquals Greater than or equal
	GreaterThanEquals SearchCriteriaConditionTypes = "gteq"

	// In In. The value can contain a comma-separated list of values.
	In SearchCriteriaConditionTypes = "in"

	// Like Like. The value can contain the SQL wildcard characters when like is specified.
	Like SearchCriteriaConditionTypes = "like"

	// LessThan Less than
	LessThan SearchCriteriaConditionTypes = "lt"

	// LessThanEquals Less than or equal
	LessThanEquals SearchCriteriaConditionTypes = "lteq"

	// // MoreThanEquals More or equal
	// MoreThanEquals SearchCriteriaConditionTypes = "moreq"

	// NotEquals Not equal
	NotEquals SearchCriteriaConditionTypes = "neq"

	// NotFindInSet A value that is not within a set of values
	NotFindInSet SearchCriteriaConditionTypes = "nfinset"

	// NotIn Not in. The value can contain a comma-separated list of values.
	NotIn SearchCriteriaConditionTypes = "nin"

	// NotNull Not null
	NotNull SearchCriteriaConditionTypes = "notnull"

	// Null Null
	Null SearchCriteriaConditionTypes = "null"

	// The beginning of a range. Must be used with to
	// SearchCriteriaConditionTypes = "from"

	// The end of a range. Must be used with from
	// SearchCriteriaConditionTypes = "to"
)

// SearchCriterialCondition comment
type SearchCriterialCondition struct {
	Field string                       `json:"field"`
	Value interface{}                  `json:"value"`
	Type  SearchCriteriaConditionTypes `json:"type"`
}

// SearchCriterial comment
type SearchCriterial struct {
	Query      string
	Conditions []*SearchCriterialCondition
}

// AddFilters comment
func (ctx *SearchCriterial) AddFilters(conditions []SearchCriterialCondition) *SearchCriterial {
	for _, condition := range conditions {
		ctx.Conditions = append(ctx.Conditions, &condition)
	}
	return ctx
}

// PrepareConditions comment
func (ctx *SearchCriterial) PrepareConditions(query string) (*SearchCriterial, error) {
	decodedValue, err := url.QueryUnescape(query)
	if err != nil {
		return ctx, err
	}
	log.Println(decodedValue)
	// err := decoder.Decode(&employeeStruct, query)
	return ctx, err
}
