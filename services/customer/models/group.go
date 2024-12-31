package models

import (
	"encoding/json"

	"localdomain/customer/core"

	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

// Group struct
type Group struct {
	DefaultGormModel
	CodeName     null.String `gorm:"column:code_name;not null;index:idx_code_name;" sql:"type:varchar(100) character set ascii collate ascii_general_ci;" json:"codeName"`
	CategoryName null.String `gorm:"column:category_name;not null;index:idx_category_name;" sql:"type:varchar(100) character set ascii collate ascii_general_ci;" json:"categoryName"`
	Title        null.String `gorm:"column:title;not null;" sql:"type:varchar(255) character set ascii collate ascii_general_ci;" json:"title"`
	Description  null.String `gorm:"column:description;" sql:"type:varchar(512) character set utf8mb4 collate utf8mb4_unicode_ci;" json:"description"`
	ParentID     null.String `gorm:"column:parent_id;index:idx_parent_id;" sql:"type:char(27) character set ascii collate ascii_general_ci;" json:"parentId"`
	SortOrder    null.Int    `gorm:"column:sort_order;not null;default:0;" sql:"type:smallint(6);" json:"sortOrder"`
	NestedSetTreeModel
	TimestampModel
	Customers     []Customer    `gorm:"many2many:group_customer_linker;"`
	GroupCategory GroupCategory `gorm:"foreignKey:category_name;references:category_name;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName specifies table name
func (ctx *Group) TableName() string {
	return "groups"
}

// FindAllGroups comment
func FindAllGroups(db *gorm.DB, catFilter *GroupFilter, opts *PaginationOptions) ([]*Group, error) {
	limit := opts.GetLimit()
	offset := opts.GetOffset()

	builder := db.
		Model(Group{}).
		Select("*")

	if catFilter.RootOnly {
		builder = builder.Where("parent_id IS NULL")
	}

	rows, err := builder.
		Order("node_left ASC, sort_order ASC").
		Limit(limit).
		Offset(offset).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := make([]*Group, 0)

	for rows.Next() {
		var group Group
		db.ScanRows(rows, &group)

		groups = append(groups, &group)
	}

	return groups, err
}

// CountAllGroups comment
func CountAllGroups(db *gorm.DB, catFilter *GroupFilter) int64 {
	var total int64 = 0

	builder := db.Model(Group{})
	if catFilter.RootOnly {
		builder = builder.Where("parent_id IS NULL")
	}
	builder.Count(&total)

	return total
}

// FindOneGroup comment
func FindOneGroup(db *gorm.DB, groupID string) (*Group, error) {
	var (
		err   error
		group Group
	)

	db.
		Model(Group{}).
		Where("id = ?", groupID).
		First(&group)

	return &group, err
}

// FindAllGroupDescendants comment
func FindAllGroupDescendants(db *gorm.DB, nodeID string) (list []*Group, err error) {
	var group Group

	db.
		Where("id = ?", nodeID).
		Find(&group)

	db.
		Where("node_right > ?", group.NodeLeft).
		Where("node_right < ?", group.NodeRight).
		Order("node_left ASC, sort_order ASC").
		Find(&list)

	return list, err
}

// FindAllGroupChildren comment
func FindAllGroupChildren(db *gorm.DB, nodeID string) (list []*Group, err error) {
	var group Group

	db.
		Where("id = ?", nodeID).
		Find(&group)

	db.
		Where("parent_id = ?", group.ID).
		Order("node_left ASC, sort_order ASC").
		Find(&list)

	return list, err
}

// FetchDetailInfo comment
func (ctx *Group) FetchDetailInfo(db *gorm.DB) (map[string]interface{}, error) {
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

// Register comment
func (ctx *Group) Register(db *gorm.DB) (err error) {
	return ctx.Transact(db, func(tx *gorm.DB) error {
		queryResult := tx.
			Where(ctx).
			FirstOrCreate(&ctx, DefaultGormModel{
				ID: null.StringFrom(core.GenerateID()),
			})

		if queryResult.RowsAffected == 0 {
			return tx.
				Where(ctx).
				Update("updated_at", gorm.Expr("CURRENT_TIMESTAMP()")).
				Error
		}
		return queryResult.Error
	})
}

// Update comment
func (ctx *Group) Update(db *gorm.DB, infoValues map[string]interface{}) (err error) {
	return ctx.Transact(db, func(tx *gorm.DB) error {
		infoValues, err = ConvertMapToKeySnakeCase(infoValues)
		if err != nil {
			return err
		}

		return db.
			Model(&ctx).
			Where(ctx).
			Updates(infoValues).
			Error
	})
}

// Delete comment
func (ctx *Group) Delete(db *gorm.DB) (err error) {
	return ctx.Transact(db, func(tx *gorm.DB) error {
		return db.
			Where(ctx).
			Delete(&ctx).
			Error
	})
}

// FindAllGroupCustomers comment
func FindAllGroupCustomers(db *gorm.DB, nodeID string, customerFilter *CustomerAttributeFilter, opts *PaginationOptions) (list []*Customer, err error) {
	var group Group

	limit := opts.GetLimit()
	offset := opts.GetOffset()

	db.
		Where("id = ?", nodeID).
		Find(&group)

	queryResult := db.
		Model(&group).
		// Order("title asc").
		Limit(limit).
		Offset(offset).
		Association("Customers")

	err = queryResult.Error
	if queryResult.Error != nil {
		return nil, queryResult.Error
	}

	queryResult.Find(&list)
	return list, err
}

// CountAllCategorieCustomers comment
func CountAllCategorieCustomers(db *gorm.DB, nodeID string, customerFilter *CustomerAttributeFilter) int64 {
	var (
		total int64 = 0
		group Group
	)

	db.
		Where("id = ?", nodeID).
		Find(&group)

	total = db.
		Model(&group).
		Association("Customers").
		Count()

	return total
}

// FindAllGroupPaths comment
func FindAllGroupPaths(db *gorm.DB, nodeID string) ([]*Group, error) {
	list := make([]*Group, 0)

	rows, err := db.
		Raw(`
			SELECT
				catParent.*
			FROM
				groups AS catNode,
				groups AS catParent
			WHERE (
				(catNode.node_left BETWEEN catParent.node_left AND catParent.node_right)
				AND catNode.id = ?
			)
			ORDER BY catParent.node_left;
		`, nodeID).
		Rows()
	if err != nil {
		return list, err
	}
	defer rows.Close()

	for rows.Next() {
		var group Group
		db.ScanRows(rows, &group)

		list = append(list, &group)
	}

	return list, err
}
