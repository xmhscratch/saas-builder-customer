package models

import (
	"encoding/json"

	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

// GroupCategory struct
type GroupCategory struct {
	CategoryName null.String `gorm:"column:category_name;primary_key;type:string;" sql:"type:varchar(100) character set ascii collate ascii_general_ci;" json:"categoryName"`
	Title        null.String `gorm:"column:title;not null;type:string;" sql:"type:varchar(255) character set utf8mb4 collate utf8mb4_unicode_ci;" json:"title"`
	Description  null.String `gorm:"column:description;type:string;" sql:"type:text character set utf8mb4 collate utf8mb4_unicode_ci;" json:"description"`
	// Groups      []Group     `gorm:"foreignKey:category_name;references:category_name;" json:"groupInfos"`
}

// TableName specifies table name
func (ctx *GroupCategory) TableName() string {
	return "group_categories"
}

// FindAllGroupCategories comment
func FindAllGroupCategories(db *gorm.DB, filter *GroupCategoryFilter) ([]*GroupCategory, error) {
	builder := db.
		Model(GroupCategory{}).
		Select("*")

	rows, err := builder.
		Where(filter).
		Order("title ASC").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groupCategories := make([]*GroupCategory, 0)

	for rows.Next() {
		var groupCategory GroupCategory
		db.ScanRows(rows, &groupCategory)

		groupCategories = append(groupCategories, &groupCategory)
	}

	return groupCategories, err
}

// CountAllGroupCategories comment
func CountAllGroupCategories(db *gorm.DB, filter *GroupCategoryFilter) int64 {
	var total int64 = 0

	db.
		Model(GroupCategory{}).
		Count(&total)

	return total
}

// FindOneGroupCategory comment
func FindOneGroupCategory(db *gorm.DB, groupCategoryID string) (*GroupCategory, error) {
	var (
		err           error
		groupCategory GroupCategory
	)

	db.
		Model(GroupCategory{}).
		Where("id = ?", groupCategoryID).
		First(&groupCategory)

	return &groupCategory, err
}

// FetchDetailInfo comment
func (ctx *GroupCategory) FetchDetailInfo(db *gorm.DB) (map[string]interface{}, error) {
	var details map[string]interface{}

	b, err := json.Marshal(*ctx)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &details)
	if err != nil {
		return nil, err
	}
	return details, err
}
