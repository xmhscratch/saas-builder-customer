package models

// GroupCategoryFilter comment
type GroupCategoryFilter struct {
	CodeNames map[string]string `gorm:"-"`
}

// TableName specifies table name
func (ctx *GroupCategoryFilter) TableName() string {
	return "e"
}
