package models

import (
	"encoding/json"

	null "gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

// Enumeration model
type Enumeration struct {
	DefaultGormModel
	Label        null.String        `gorm:"column:label" sql:"type:varchar(255)" json:"label"`
	Description  null.String        `gorm:"column:description" sql:"type:varchar(512)" json:"description"`
	DefaultValue null.Int           `gorm:"column:default_value" sql:"type:int(11)" json:"defaultValue"`
	ListRenderer null.String        `gorm:"column:list_renderer" sql:"type:varchar(255)" json:"listRenderer"`
	Values       []EnumerationValue `gorm:"foreignKey:EnumID" json:"values"`
}

// TableName specifies table name
func (ctx *Enumeration) TableName() string {
	return "enumerations"
}

// Register comment
func (ctx *Enumeration) Register(db *gorm.DB) (err error) {
	return ctx.Transact(db, func(tx *gorm.DB) error {
		return tx.Create(&ctx).Error
	})
}

// FindAllEnumerations comment
func FindAllEnumerations(db *gorm.DB, enumFilter *EnumerationFilter) ([]*Enumeration, error) {
	rows, err := db.
		Model(Enumeration{}).
		// Table("enumerations AS e").
		Where(enumFilter).
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	enums := make([]*Enumeration, 0)

	for rows.Next() {
		var enum Enumeration
		db.ScanRows(rows, &enum)

		enums = append(enums, &enum)
	}

	return enums, err
}

// CountAllEnumerations comment
func CountAllEnumerations(db *gorm.DB, filter *EnumerationFilter) int64 {
	var total int64 = 0

	db.Model(Enumeration{}).
		// Table("enumerations AS e").
		Where(filter).
		Count(&total)

	return total
}

// FindOneEnumeration comment
func FindOneEnumeration(db *gorm.DB, id string) (*Enumeration, error) {
	var (
		err  error
		enum Enumeration
	)

	db.
		Model(Enumeration{}).
		Where("id = ?", id).
		First(&enum)

	db.
		Model(&enum).
		Association("Values").
		Find(&enum.Values)

	return &enum, err
}

// FetchDetailInfo comment
func (ctx *Enumeration) FetchDetailInfo(db *gorm.DB) (map[string]interface{}, error) {
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
