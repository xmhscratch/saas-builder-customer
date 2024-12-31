package models

import (
	"encoding/json"
	// "log"
	"strings"

	"localdomain/customer/core"

	null "gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Attribute model
type Attribute struct {
	DefaultGormModel
	// ID                    null.String `gorm:"column:id;primary_key" sql:"type:varchar(27)" json:"id"`
	EntityType       null.String `gorm:"column:entity_type;not null;type:string;default:varchar;" sql:"type:enum('varchar','datetime','decimal','int','boolean','text','blob');not null;" json:"entityType"`
	CodeName         null.String `gorm:"column:code_name;not null" sql:"type:varchar(100);not null" json:"codeName"`
	Metadata         null.String `gorm:"column:metadata" sql:"type:mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin" json:"metadata"`
	Label            null.String `gorm:"column:label;not null" sql:"type:varchar(255)" json:"label"`
	Description      null.String `gorm:"column:description" sql:"type:varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci" json:"description"`
	DefaultValue     null.String `gorm:"column:default_value" sql:"type:text" json:"defaultValue"`
	SourceData       null.String `gorm:"column:source_data" sql:"type:varchar(255)" json:"sourceData"`
	InputRenderer    null.String `gorm:"column:input_renderer" sql:"type:varchar(255)" json:"inputRenderer"`
	ListRenderer     null.String `gorm:"column:list_renderer" sql:"type:varchar(255)" json:"listRenderer"`
	IsFilterable     null.Bool   `gorm:"column:is_filterable;type:bool;not null;default:0;" sql:"type:tinyint(1)" json:"isFilterable"`
	IsVisibleOnFront null.Bool   `gorm:"column:is_visible_on_front;type:bool;not null;default:1;" sql:"type:tinyint(1)" json:"isVisibleOnFront"`
	IsVisibleInList  null.Bool   `gorm:"column:is_visible_in_list;type:bool;not null;default:0;" sql:"type:tinyint(1)" json:"isVisibleInList"`
	IsConfigurable   null.Bool   `gorm:"column:is_configurable;type:bool;not null;default:0;" sql:"type:tinyint(1)" json:"isConfigurable"`
	IsUserDefined    null.Bool   `gorm:"column:is_user_defined;type:bool;not null;default:1;" sql:"type:tinyint(1)" json:"isUserDefined"`
	IsReadOnly       null.Bool   `gorm:"column:is_read_only;type:bool;not null;default:0;" sql:"type:tinyint(1)" json:"isReadOnly"`
	IsRequired       null.Bool   `gorm:"column:is_required;type:bool;not null;default:0;" sql:"type:tinyint(1)" json:"isRequired"`
	IsUnique         null.Bool   `gorm:"column:is_unique;type:bool;not null;default:0;" sql:"type:tinyint(1)" json:"isUnique"`
	ListColumnSize   null.Int    `gorm:"column:list_column_size;not null;default:3" sql:"type:smallint(2)" json:"listColumnSize"`
	DisplayFormat    null.String `gorm:"column:display_format" sql:"type:varchar(12)" json:"displayFormat"`
	SortOrder        null.Int    `gorm:"column:sort_order;not null;default:0;" sql:"type:int(11)" json:"sortOrder"`
	Note             null.String `gorm:"column:note" sql:"type:tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci" json:"note"`
	TimestampModel
	CustomerAttributeBlobs     []CustomerAttributeBlob     `gorm:"foreignKey:attribute_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CustomerAttributeBooleans  []CustomerAttributeBoolean  `gorm:"foreignKey:attribute_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CustomerAttributeDateTimes []CustomerAttributeDateTime `gorm:"foreignKey:attribute_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CustomerAttributeDecimals  []CustomerAttributeDecimal  `gorm:"foreignKey:attribute_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CustomerAttributeInts      []CustomerAttributeInt      `gorm:"foreignKey:attribute_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CustomerAttributeTexts     []CustomerAttributeText     `gorm:"foreignKey:attribute_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CustomerAttributeVarchars  []CustomerAttributeVarchar  `gorm:"foreignKey:attribute_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName specifies table name
func (ctx *Attribute) TableName() string {
	return "attributes"
}

// MarshalJSON implements json.Marshaler.
func (ctx *Attribute) MarshalJSON() ([]byte, error) {
	var (
		err error
		m   map[string]interface{}
	)
	type AttributeClone Attribute

	if ctx.Metadata.IsZero() {
		m = nil
	} else {
		b := []byte(ctx.Metadata.ValueOrZero())
		if strings.ReplaceAll(strings.TrimSpace(string(b[:])), " ", "") == "{}" {
			m = nil
		} else {
			err = json.Unmarshal(b, &m)
		}
	}

	if err != nil {
		return []byte("null"), err
	}

	return json.Marshal(&struct {
		Metadata map[string]interface{} `json:"metadata"`
		*AttributeClone
	}{
		Metadata:       m,
		AttributeClone: (*AttributeClone)(ctx),
	})
}

// Normalize comment
func (ctx *Attribute) Normalize(infoValues map[string]interface{}) (err error) {
	if infoValues["entityType"] != nil {
		val, err := ParseString(infoValues["entityType"])
		if err != nil {
			return err
		}
		ctx.EntityType = null.StringFromPtr(&val)
	}
	if infoValues["codeName"] != nil {
		val, err := ParseString(infoValues["codeName"])
		if err != nil {
			return err
		}
		ctx.CodeName = null.StringFromPtr(&val)
	}
	if infoValues["metadata"] != nil {
		val, err := ParseString(infoValues["metadata"])
		if err != nil {
			return err
		}
		ctx.Metadata = null.StringFromPtr(&val)
	}
	if infoValues["label"] != nil {
		val, err := ParseString(infoValues["label"])
		if err != nil {
			return err
		}
		ctx.Label = null.StringFromPtr(&val)
	}
	if infoValues["description"] != nil {
		val, err := ParseString(infoValues["description"])
		if err != nil {
			return err
		}
		ctx.Description = null.StringFromPtr(&val)
	}
	if infoValues["defaultValue"] != nil {
		val, err := ParseString(infoValues["defaultValue"])
		if err != nil {
			return err
		}
		ctx.DefaultValue = null.StringFromPtr(&val)
	}
	if infoValues["sourceData"] != nil {
		val, err := ParseString(infoValues["sourceData"])
		if err != nil {
			return err
		}
		ctx.SourceData = null.StringFromPtr(&val)
	}
	if infoValues["inputRenderer"] != nil {
		val, err := ParseString(infoValues["inputRenderer"])
		if err != nil {
			return err
		}
		ctx.InputRenderer = null.StringFromPtr(&val)
	}
	if infoValues["listRenderer"] != nil {
		val, err := ParseString(infoValues["listRenderer"])
		if err != nil {
			return err
		}
		ctx.ListRenderer = null.StringFromPtr(&val)
	}
	if infoValues["isFilterable"] != nil {
		val, err := ParseBool(infoValues["isFilterable"])
		if err != nil {
			return err
		}
		ctx.IsFilterable = null.BoolFromPtr(&val)
	}
	if infoValues["isVisibleOnFront"] != nil {
		val, err := ParseBool(infoValues["isVisibleOnFront"])
		if err != nil {
			return err
		}
		ctx.IsVisibleOnFront = null.BoolFromPtr(&val)
	}
	if infoValues["isVisibleInList"] != nil {
		val, err := ParseBool(infoValues["isVisibleInList"])
		if err != nil {
			return err
		}
		ctx.IsVisibleInList = null.BoolFromPtr(&val)
	}
	if infoValues["isConfigurable"] != nil {
		val, err := ParseBool(infoValues["isConfigurable"])
		if err != nil {
			return err
		}
		ctx.IsConfigurable = null.BoolFromPtr(&val)
	}
	if infoValues["isUserDefined"] != nil {
		val, err := ParseBool(infoValues["isUserDefined"])
		if err != nil {
			return err
		}
		ctx.IsUserDefined = null.BoolFromPtr(&val)
	}
	if infoValues["isReadOnly"] != nil {
		val, err := ParseBool(infoValues["isReadOnly"])
		if err != nil {
			return err
		}
		ctx.IsReadOnly = null.BoolFromPtr(&val)
	}
	if infoValues["isRequired"] != nil {
		val, err := ParseBool(infoValues["isRequired"])
		if err != nil {
			return err
		}
		ctx.IsRequired = null.BoolFromPtr(&val)
	}
	if infoValues["isUnique"] != nil {
		val, err := ParseBool(infoValues["isUnique"])
		if err != nil {
			return err
		}
		ctx.IsUnique = null.BoolFromPtr(&val)
	}
	if infoValues["listColumnSize"] != nil {
		val, err := ParseInt(infoValues["listColumnSize"])
		if err != nil {
			return err
		}
		ctx.ListColumnSize = null.IntFromPtr(&val)
	}
	if infoValues["displayFormat"] != nil {
		val, err := ParseString(infoValues["displayFormat"])
		if err != nil {
			return err
		}
		ctx.DisplayFormat = null.StringFromPtr(&val)
	}
	if infoValues["sortOrder"] != nil {
		val, err := ParseInt(infoValues["sortOrder"])
		if err != nil {
			return err
		}
		ctx.SortOrder = null.IntFromPtr(&val)
	}
	if infoValues["note"] != nil {
		val, err := ParseString(infoValues["note"])
		if err != nil {
			return err
		}
		ctx.Note = null.StringFromPtr(&val)
	}
	return err
}

// Register comment
func (ctx *Attribute) Register(db *gorm.DB) (err error) {
	return ctx.Transact(db, func(tx *gorm.DB) error {
		var (
			codeName  null.String = ctx.CodeName
			sortOrder int64
		)

		if err := tx.
			// Debug().
			Model(ctx).
			Where("code_name", codeName).
			FirstOrInit(&ctx).
			Error; err != nil {
			return err
		}

		if err := tx.
			// Debug().
			Model(ctx).
			Select("COUNT(id)").
			Count(&sortOrder).
			Error; err != nil {
			return err
		}
		ctx.SortOrder = null.IntFrom(sortOrder + 1)

		if ctx.DefaultGormModel.ID.IsZero() || ctx.DefaultGormModel.ID.ValueOrZero() == "" {
			ctx.DefaultGormModel.ID = null.StringFrom(core.GenerateID())
		}

		return tx.
			// Debug().
			Model(ctx).
			Clauses(
				clause.OnConflict{
					Columns: []clause.Column{
						{Name: "updated_at"},
					},
					DoUpdates: clause.Assignments(map[string]interface{}{
						"updated_at": gorm.Expr("CURRENT_TIMESTAMP()"),
					}),
				},
			).
			FirstOrCreate(&ctx, DefaultGormModel{
				ID: null.StringFrom(core.GenerateID()),
			}).
			Error
	})
}

// Update comment
func (ctx *Attribute) Update(db *gorm.DB, infoValues map[string]interface{}) (err error) {
	return ctx.Transact(db, func(tx *gorm.DB) error {
		allowedAttributes := []string{
			// "entityType",
			// "codeName",
			"metadata",
			"label",
			"description",
			"defaultValue",
			// "sourceData",
			// "inputRenderer",
			// "listRenderer",
			"isFilterable",
			// "isVisibleOnFront",
			"isVisibleInList",
			// "isConfigurable",
			// "isUserDefined",
			"isReadOnly",
			"isRequired",
			"isUnique",
			"listColumnSize",
			// "displayFormat",
			"sortOrder",
			"note",
		}
		infoValues = core.FilterMap(infoValues, func(val interface{}, key string) bool {
			return core.IncludeInStringSlice(key, allowedAttributes)
		})
		infoValues, err = ConvertMapToKeySnakeCase(infoValues)
		if err != nil {
			return err
		}
		infoValues["updated_at"] = gorm.Expr("CURRENT_TIMESTAMP()")

		return db.
			// Debug().
			Model(&ctx).
			Updates(infoValues).
			Error
	})
}

// Delete comment
func (ctx *Attribute) Delete(db *gorm.DB) (err error) {
	return ctx.Transact(db, func(tx *gorm.DB) error {
		return db.
			Where(ctx).
			Delete(&ctx).
			Error
	})
}

// FindAllAttributes comment
func FindAllAttributes(db *gorm.DB, attrFilter *AttributeFilter, opts *PaginationOptions) ([]*Attribute, error) {
	// limit := opts.GetLimit()
	// offset := opts.GetOffset()

	queryBuilder := db.
		// Debug().
		Model(&Attribute{}).
		Table("attributes AS e")

	queryBuilder = queryBuilder.
		// Debug().
		Where(attrFilter)
		// Where(attrFilter.FilterQueryString()).
		// Order(attrFilter.SortByQueryString())
		// Limit(limit).
		// Offset(offset).

	// if true == attrFilter.IsFilterable.ValueOrZero() {
	// 	queryBuilder = queryBuilder.Where("is_filterable = ?", attrFilter.IsFilterable.ValueOrZero())
	// }
	// if true == attrFilter.IsVisibleOnFront.ValueOrZero() {
	// 	queryBuilder = queryBuilder.Where("is_visible_on_front = ?", attrFilter.IsVisibleOnFront.ValueOrZero())
	// }
	// if true == attrFilter.IsVisibleInList.ValueOrZero() {
	// 	queryBuilder = queryBuilder.Where("is_visible_in_list = ?", attrFilter.IsVisibleInList.ValueOrZero())
	// }
	// if true == attrFilter.IsConfigurable.ValueOrZero() {
	// 	queryBuilder = queryBuilder.Where("is_configurable = ?", attrFilter.IsConfigurable.ValueOrZero())
	// }
	// if true == attrFilter.IsUserDefined.ValueOrZero() {
	// 	queryBuilder = queryBuilder.Where("is_user_defined = ?", attrFilter.IsUserDefined.ValueOrZero())
	// }
	// if true == attrFilter.IsReadOnly.ValueOrZero() {
	// 	queryBuilder = queryBuilder.Where("is_read_only = ?", attrFilter.IsReadOnly.ValueOrZero())
	// }
	// if true == attrFilter.IsRequired.ValueOrZero() {
	// 	queryBuilder = queryBuilder.Where("is_required = ?", attrFilter.IsRequired.ValueOrZero())
	// }
	// if true == attrFilter.IsUnique.ValueOrZero() {
	// 	queryBuilder = queryBuilder.Where("is_unique = ?", attrFilter.IsUnique.ValueOrZero())
	// }
	// log.Println(attrFilter.FilterQueryString())
	// log.Println(attrFilter.SortByQueryString())

	queryBuilder = queryBuilder.
		Where(attrFilter.FilterQueryString()).
		Order(attrFilter.SortByQueryString())

	rows, err := queryBuilder.
		// Debug().
		Rows()
	if err != nil {
		// log.Println(err)
		return nil, err
	}
	defer rows.Close()

	attrs := make([]*Attribute, 0)

	for rows.Next() {
		var attr Attribute
		db.ScanRows(rows, &attr)

		if attrFilter.WithAssociations {
			db.
				// Debug().
				Model(&attr).
				Association("CustomerAttributeBlobs").
				Find(&attr.CustomerAttributeBlobs)
			db.
				// Debug().
				Model(&attr).
				Association("CustomerAttributeBooleans").
				Find(&attr.CustomerAttributeBooleans)
			db.
				// Debug().
				Model(&attr).
				Association("CustomerAttributeDateTimes").
				Find(&attr.CustomerAttributeDateTimes)
			db.
				// Debug().
				Model(&attr).
				Association("CustomerAttributeDecimals").
				Find(&attr.CustomerAttributeDecimals)
			db.
				// Debug().
				Model(&attr).
				Association("CustomerAttributeInts").
				Find(&attr.CustomerAttributeInts)
			db.
				// Debug().
				Model(&attr).
				Association("CustomerAttributeTexts").
				Find(&attr.CustomerAttributeTexts)
			db.
				// Debug().
				Model(&attr).
				Association("CustomerAttributeVarchars").
				Find(&attr.CustomerAttributeVarchars)
		}

		attrs = append(attrs, &attr)
	}

	return attrs, nil
}

// CountAllAttributes comment
func CountAllAttributes(db *gorm.DB, attrFilter *AttributeFilter) int64 {
	var total int64

	queryBuilder := db.Model(Attribute{}).
		Table("attributes AS e")

	queryBuilder = queryBuilder.
		// Debug().
		Where(attrFilter).
		Where(attrFilter.FilterQueryString()).
		Order(attrFilter.SortByQueryString())
		// Limit(limit).
		// Offset(offset).

	// if true == attrFilter.IsFilterable.ValueOrZero() {
	// 	queryBuilder = queryBuilder.Where("is_filterable = ?", attrFilter.IsFilterable.ValueOrZero())
	// }
	// if true == attrFilter.IsVisibleOnFront.ValueOrZero() {
	// 	queryBuilder = queryBuilder.Where("is_visible_on_front = ?", attrFilter.IsVisibleOnFront.ValueOrZero())
	// }
	// if true == attrFilter.IsVisibleInList.ValueOrZero() {
	// 	queryBuilder = queryBuilder.Where("is_visible_in_list = ?", attrFilter.IsVisibleInList.ValueOrZero())
	// }
	// if true == attrFilter.IsUserDefined.ValueOrZero() {
	// 	queryBuilder = queryBuilder.Where("is_user_defined = ?", attrFilter.IsUserDefined.ValueOrZero())
	// }
	// if true == attrFilter.IsReadOnly.ValueOrZero() {
	// 	queryBuilder = queryBuilder.Where("is_read_only = ?", attrFilter.IsReadOnly.ValueOrZero())
	// }
	// if true == attrFilter.IsRequired.ValueOrZero() {
	// 	queryBuilder = queryBuilder.Where("is_required = ?", attrFilter.IsRequired.ValueOrZero())
	// }
	// if true == attrFilter.IsUnique.ValueOrZero() {
	// 	queryBuilder = queryBuilder.Where("is_unique = ?", attrFilter.IsUnique.ValueOrZero())
	// }

	err := queryBuilder.
		Count(&total).
		Error
	if err != nil {
		return int64(0)
	}
	return total
}

// FindOneAttribute comment
func FindOneAttribute(db *gorm.DB, codeName string) (*Attribute, error) {
	var (
		err  error
		attr Attribute
	)

	err = db.
		Model(Attribute{}).
		Where("code_name = ?", codeName).
		First(&attr).
		Error

	return &attr, err
}

// BuildFilterAttributeTextQuery comment
func BuildFilterAttributeTextQuery(db *gorm.DB, queryString string) (string, error) {
	var (
		err         error
		attrs       []*Attribute
		results     []string
		queryResult string
	)

	attrFilter := &AttributeFilter{
		EntityTypes: map[string]string{
			"0": "varchar",
			"1": "text",
		},
		IsFilterable: null.BoolFrom(true),
	}

	attrs, err = FindAllAttributes(db, attrFilter, nil)
	if err != nil {
		return "", err
	}

	results = append(results, core.BuildString(
		"emailAddress",
		":",
		core.BuildString("(*", queryString, "*~1000", ")"),
	))

	for _, attr := range attrs {
		results = append(results, core.BuildString(
			attr.CodeName.ValueOrZero(),
			"_",
			attr.EntityType.ValueOrZero(),
			":",
			core.BuildString("(*", queryString, "*~1000", ")"),
		))
	}
	if err != nil {
		return "", err
	}

	queryResult = strings.Join(results, " OR ")
	if len(queryResult) == 0 {
		return "*:*", err
	}
	return queryResult, err
}

// FetchDetailInfo comment
func (ctx *Attribute) FetchDetailInfo(db *gorm.DB) (map[string]interface{}, error) {
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
