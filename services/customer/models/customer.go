package models

import (
	"encoding/json"
	"errors"
	"fmt"

	// "log"
	// "reflect"
	"strings"
	"time"

	// sq "github.com/Masterminds/squirrel"

	"localdomain/customer/core"

	null "gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Customer model
type Customer struct {
	DefaultGormModel
	EmailAddress      null.String `gorm:"<-:create;column:email_address;uniqueIndex:uniq_email,unique;uniqueIndex:uniq_email_with_shopify,unique;" sql:"type:varchar(512);not null;unique_index:unq_email;unique_index:unq_email_with_shopify;" json:"emailAddress"`
	ShopifyCustomerID null.String `gorm:"column:shopify_customer_id;uniqueIndex:uniq_email_with_shopify,unique;" sql:"type:datetime;unique_index:unq_email_with_shopify;" json:"shopifyCustomerId"`
	InitialPoints     null.Float  `gorm:"column:initial_points;default:0;" sql:"type:decimal(12,4);not null;" json:"initialPoints"`
	CurrentPoints     null.Float  `gorm:"<-:false;column:current_points;default:0;" sql:"type:decimal(12,4);not null;" json:"currentPoints"`
	EarnedPoints      null.Float  `gorm:"<-:false;column:earned_points;default:0;" sql:"type:decimal(12,4);not null;" json:"earnedPoints"`
	SpentPoints       null.Float  `gorm:"<-:false;column:spent_points;default:0;" sql:"type:decimal(12,4);not null;" json:"spentPoints"`
	SyncedAt          null.Time   `gorm:"column:synced_at" sql:"type:datetime;" json:"syncedAt"`
	Groups            []Group     `gorm:"many2many:group_customer_linker" json:"groupInfos"`
	TimestampModel
	CustomerAttributeBlobs     []CustomerAttributeBlob     `gorm:"foreignKey:customer_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CustomerAttributeBooleans  []CustomerAttributeBoolean  `gorm:"foreignKey:customer_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CustomerAttributeDateTimes []CustomerAttributeDateTime `gorm:"foreignKey:customer_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CustomerAttributeDecimals  []CustomerAttributeDecimal  `gorm:"foreignKey:customer_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CustomerAttributeInts      []CustomerAttributeInt      `gorm:"foreignKey:customer_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CustomerAttributeTexts     []CustomerAttributeText     `gorm:"foreignKey:customer_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CustomerAttributeVarchars  []CustomerAttributeVarchar  `gorm:"foreignKey:customer_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName specifies table name
func (ctx *Customer) TableName() string {
	return "customers"
}

// CustomerAttributeFilter comment
type CustomerAttributeFilter struct {
	SortBy map[string]string
}

// FindAllCustomers comment
func FindAllCustomers(db *gorm.DB, filter *CustomerAttributeFilter, opts *PaginationOptions) ([]*Customer, error) {
	limit := opts.GetLimit()
	offset := opts.GetOffset()

	builder := db.
		// Debug().
		Model(Customer{}).
		Select([]string{
			"customers.*",
			"IFNULL(SUM(COALESCE(point_balance_histories.balance_income, 0.0000) - COALESCE(point_balance_histories.balance_outcome, 0.0000)), 0.0000) AS current_points",
			"IFNULL(SUM(COALESCE(point_balance_histories.balance_income, 0.0000)), 0.0000) AS earned_points",
			"IFNULL(SUM(COALESCE(point_balance_histories.balance_outcome, 0.0000)), 0.0000) AS spent_points",
		}).
		Joins("JOIN point_balance_histories ON point_balance_histories.customer_id = customers.id")

	rows, err := builder.
		Limit(limit).
		Offset(offset).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := make([]*Customer, 0)

	for rows.Next() {
		var customer Customer
		db.ScanRows(rows, &customer)

		customers = append(customers, &customer)
	}

	return customers, err
}

// CountAllCustomers comment
func CountAllCustomers(db *gorm.DB, filter *CustomerAttributeFilter) int64 {
	var total int64 = 0

	db.
		Model(Customer{}).
		Count(&total)

	return total
}

// FindOneCustomer comment
func FindOneCustomer(db *gorm.DB, customerID string) (*Customer, error) {
	var (
		err      error
		customer Customer
	)

	db.
		// Debug().
		Model(Customer{}).
		Table(
			"customers, (?) AS current_points, (?) AS earned_points, (?) AS spent_points",
			db.
				Model(&PointBalanceHistory{}).
				Select("IFNULL(SUM(COALESCE(point_balance_histories.balance_income, 0.0000) - COALESCE(point_balance_histories.balance_outcome, 0.0000)), 0.0000) AS current_points").
				Where("customer_id = ?", customerID),
			db.
				Model(&PointBalanceHistory{}).
				Select("IFNULL(SUM(COALESCE(point_balance_histories.balance_income, 0.0000)), 0.0000) AS earned_points").
				Where("customer_id = ?", customerID),
			db.
				Model(&PointBalanceHistory{}).
				Select("IFNULL(SUM(COALESCE(point_balance_histories.balance_outcome, 0.0000)), 0.0000) AS spent_points").
				Where("customer_id = ?", customerID),
		).
		Where("customers.id = ?", customerID).
		Limit(1).
		Find(&customer)

	return &customer, err
}

// FindOneCustomerByEmail comment
func FindOneCustomerByEmail(db *gorm.DB, emailAddress string) (*Customer, error) {
	var (
		err      error
		customer Customer
	)

	db.
		// Debug().
		Model(Customer{}).
		Select([]string{
			"customers.*",
			"IFNULL(SUM(COALESCE(point_balance_histories.balance_income, 0.0000) - COALESCE(point_balance_histories.balance_outcome, 0.0000)), 0.0000) AS current_points",
			"IFNULL(SUM(COALESCE(point_balance_histories.balance_income, 0.0000)), 0.0000) AS earned_points",
			"IFNULL(SUM(COALESCE(point_balance_histories.balance_outcome, 0.0000)), 0.0000) AS spent_points",
		}).
		Joins("JOIN point_balance_histories ON point_balance_histories.customer_id = customers.id").
		Where("customers.email_address = ?", emailAddress).
		First(&customer)

	return &customer, err
}

// FindCustomerID comment
func FindCustomerID(db *gorm.DB, emailAddress string) (string, error) {
	var (
		err        error
		customer   Customer
		customerID string
	)

	err = db.
		// Debug().
		Model(Customer{}).
		Select("id").
		Where("email_address = ?", emailAddress).
		First(&customer).
		Error

	if err != nil {
		return "", err
	}

	if customer.DefaultGormModel.ID.IsZero() {
		return "", err
	}

	customerID = customer.DefaultGormModel.ID.ValueOrZero()
	return customerID, err
}

// FindOneCustomerByEmailWithShopify comment
func FindOneCustomerByEmailWithShopify(db *gorm.DB, emailAddress string, shopifyCustomerID string) (*Customer, error) {
	var (
		err      error
		customer Customer
	)

	db.
		// Debug().
		Model(Customer{}).
		Select([]string{
			"customers.*",
			"IFNULL(SUM(COALESCE(point_balance_histories.balance_income, 0.0000) - COALESCE(point_balance_histories.balance_outcome, 0.0000)), 0.0000) AS current_points",
			"IFNULL(SUM(COALESCE(point_balance_histories.balance_income, 0.0000)), 0.0000) AS earned_points",
			"IFNULL(SUM(COALESCE(point_balance_histories.balance_outcome, 0.0000)), 0.0000) AS spent_points",
		}).
		Joins("JOIN point_balance_histories ON point_balance_histories.customer_id = customers.id").
		Where("customers.email_address = ?", emailAddress).
		Where("customers.shopify_customer_id = ?", shopifyCustomerID).
		First(&customer)

	return &customer, err
}

// Register comment
func (ctx *Customer) Register(db *gorm.DB, inputPassword string, forceCreate bool) error {
	var (
		emailAddress      null.String = ctx.EmailAddress
		shopifyCustomerID null.String = ctx.ShopifyCustomerID
	)

	err := db.
		// Debug().
		Model(ctx).
		Or("email_address", emailAddress).
		FirstOrInit(&ctx).
		Error
	if err != nil {
		// log.Println(err)
		return err
	}

	if !forceCreate && !emailAddress.IsZero() {
		isValid := ctx.ValidatePassword(db, inputPassword)
		if !isValid {
			return errors.New("incorrect email or password")
		}
	}

	if shopifyCustomerID.IsZero() || shopifyCustomerID.ValueOrZero() == "" {
		shopifyCustomerID = null.StringFromPtr(nil)
	}

	if !ctx.ShopifyCustomerID.IsZero() || ctx.ShopifyCustomerID.ValueOrZero() != "" {
		shopifyCustomerID = ctx.ShopifyCustomerID
	}

	if ctx.DefaultGormModel.ID.IsZero() || ctx.DefaultGormModel.ID.ValueOrZero() == "" {
		ctx.DefaultGormModel.ID = null.StringFrom(core.GenerateID())
	}

	return db.
		// Debug().
		Model(ctx).
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{
					{Name: "updated_at"},
					{Name: "shopify_customer_id"},
				},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"updated_at":          gorm.Expr("CURRENT_TIMESTAMP()"),
					"shopify_customer_id": shopifyCustomerID,
				}),
			},
		).
		Create(&ctx).
		Error
}

// Update comment
func (ctx *Customer) Update(db *gorm.DB, infoValues map[string]interface{}) (err error) {
	infoValues, err = ConvertMapToKeySnakeCase(infoValues)
	if err != nil {
		return err
	}
	infoValues["updated_at"] = gorm.Expr("CURRENT_TIMESTAMP()")

	return db.
		// Debug().
		Model(&ctx).
		Or("email_address", ctx.EmailAddress.ValueOrZero()).
		Updates(infoValues).
		Error
}

// Delete comment
func (ctx *Customer) Delete(db *gorm.DB) (err error) {
	return ctx.Transact(db, func(tx *gorm.DB) error {
		return db.
			Where(ctx).
			Delete(&ctx).
			Error
	})
}

// FetchDetailInfo comment
func (ctx *Customer) FetchDetailInfo(db *gorm.DB, attrFilter *AttributeFilter) (map[string]interface{}, error) {
	var attrs []*CustomerAttribute
	attrMap, err := ctx.GetAttributes(db, attrFilter)
	if err != nil {
		return nil, err
	}
	for _, v := range attrMap {
		attrs = append(attrs, v)
	}

	var dict []*CustomerAttributeDictionary
	b, err := json.Marshal(attrs)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &dict)
	if err != nil {
		return nil, err
	}

	details := map[string]interface{}{
		"id":                ctx.ID,
		"emailAddress":      ctx.EmailAddress,
		"shopifyCustomerId": ctx.ShopifyCustomerID,
		"initialPoints":     ctx.InitialPoints,
		"currentPoints":     ctx.CurrentPoints,
		"earnedPoints":      ctx.EarnedPoints,
		"spentPoints":       ctx.SpentPoints,
		"createdAt":         ctx.CreatedAt,
		"updatedAt":         ctx.UpdatedAt,
	}
	for _, a := range dict {
		switch a.DisplayFormat {
		case "json":
			{
				if a.Value == nil {
					details[a.CodeName] = make(map[string]interface{})
					break
				}

				var rawValue string
				if rawValue, err = ParseString(a.Value); err != nil {
					details[a.CodeName] = make(map[string]interface{})
					break
				}

				var vals interface{}
				_ = json.Unmarshal([]byte(rawValue), &vals)

				switch formattedValue := vals.(type) {
				case map[string]interface{}, []interface{}:
					{
						details[a.CodeName] = formattedValue
						break
					}
				default:
					{
						details[a.CodeName] = []string{rawValue}
						break
					}
				}

				break
			}
		case "array":
			{
				if a.Value == nil {
					details[a.CodeName] = make([]string, 0)
					break
				}

				var rawValue string
				if rawValue, err = ParseString(a.Value); err != nil {
					details[a.CodeName] = make([]string, 0)
					break
				}

				var vals interface{}
				_ = json.Unmarshal([]byte(rawValue), &vals)

				switch formattedValue := vals.(type) {
				case []interface{}, []string:
					{
						details[a.CodeName] = formattedValue
						break
					}
				default:
					{
						details[a.CodeName] = []string{rawValue}
						break
					}
				}

				break
			}
		default:
			{
				if a.Value2 != nil {
					details[a.CodeName] = []interface{}{a.Value, a.Value2}
				} else {
					details[a.CodeName] = a.Value
				}
				break
			}
		}

		switch true {
		case a.EntityType == "datetime" && a.Value2 != nil:
			{
				details[a.CodeName] = strings.Join([]string{a.Value.(string), a.Value2.(string)}, " - ")
				break
			}
		default:
			{
				break
			}
		}
	}

	createdAt := details["createdAt"].(null.Time).ValueOrZero()
	details["createdAt"] = createdAt.Format(time.RFC3339)

	return details, err
}

// GetAttributes comment
func (ctx *Customer) GetAttributes(db *gorm.DB, attrFilter *AttributeFilter) (map[string]*CustomerAttribute, error) {
	list := make(map[string]*CustomerAttribute)

	defaultAttrs, err := _acquireDefaultAttributes(db, attrFilter)
	if err != nil {
		return nil, err
	}

	for _, cAttr := range defaultAttrs {
		if cAttr.CodeName.IsZero() {
			continue
		}
		codeName := cAttr.CodeName.ValueOrZero()
		list[codeName] = cAttr
	}

	queryExprVarcharBuilder := db.
		Session(&gorm.Session{DryRun: true}).
		Select("e.code_name, e.sort_order, e.display_format, 'varchar' AS entity_type, v.value, NULL AS metadata").
		Table("attributes AS e").
		Joins("JOIN customer_attribute_varchar AS v ON v.attribute_id = e.id").
		Where(fmt.Sprintf("v.customer_id = '%s'", ctx.ID.ValueOrZero()))

	if attrFilter != nil {
		queryExprVarcharBuilder = queryExprVarcharBuilder.
			// Where(attrFilter).
			Where(attrFilter.FilterQueryString()).
			Order(attrFilter.SortByQueryString())
	}
	queryExprVarcharStmt := queryExprVarcharBuilder.Find(nil).Statement
	queryExprVarchar := db.Dialector.Explain(queryExprVarcharStmt.SQL.String(), queryExprVarcharStmt.Vars...)
	// log.Println(queryExprVarchar)

	queryExprDatetimeBuilder := db.
		Session(&gorm.Session{DryRun: true}).
		Select("e.code_name, e.sort_order, e.display_format, 'datetime' AS entity_type, CONCAT_WS(' - ', IF(TRIM(v.value)='',NULL,v.value), IF(TRIM(v.value2)='',NULL,v.value2)) AS value, NULL AS metadata").
		Table("attributes AS e").
		Joins("JOIN customer_attribute_datetime AS v ON v.attribute_id = e.id").
		Where(fmt.Sprintf("v.customer_id = '%s'", ctx.ID.ValueOrZero()))

	if attrFilter != nil {
		queryExprDatetimeBuilder = queryExprDatetimeBuilder.
			// Where(attrFilter).
			Where(attrFilter.FilterQueryString()).
			Order(attrFilter.SortByQueryString())
	}
	queryExprDatetimeStmt := queryExprDatetimeBuilder.Find(nil).Statement
	queryExprDatetime := db.Dialector.Explain(queryExprDatetimeStmt.SQL.String(), queryExprDatetimeStmt.Vars...)

	queryExprDecimalBuilder := db.
		Session(&gorm.Session{DryRun: true}).
		Select("e.code_name, e.sort_order, e.display_format, 'decimal' AS entity_type, v.value, NULL AS metadata").
		Table("attributes AS e").
		Joins("JOIN customer_attribute_decimal AS v ON v.attribute_id = e.id").
		Where(fmt.Sprintf("v.customer_id = '%s'", ctx.ID.ValueOrZero()))

	if attrFilter != nil {
		queryExprDecimalBuilder = queryExprDecimalBuilder.
			// Where(attrFilter).
			Where(attrFilter.FilterQueryString()).
			Order(attrFilter.SortByQueryString())
	}
	queryExprDecimalStmt := queryExprDecimalBuilder.Find(nil).Statement
	queryExprDecimal := db.Dialector.Explain(queryExprDecimalStmt.SQL.String(), queryExprDecimalStmt.Vars...)

	queryExprIntBuilder := db.
		Session(&gorm.Session{DryRun: true}).
		Select("e.code_name, e.sort_order, e.display_format, 'int' AS entity_type, v.value, NULL AS metadata").
		Table("attributes AS e").
		Joins("JOIN customer_attribute_int AS v ON v.attribute_id = e.id").
		Where(fmt.Sprintf("v.customer_id = '%s'", ctx.ID.ValueOrZero()))

	if attrFilter != nil {
		queryExprIntBuilder = queryExprIntBuilder.
			// Where(attrFilter).
			Where(attrFilter.FilterQueryString()).
			Order(attrFilter.SortByQueryString())
	}
	queryExprIntStmt := queryExprIntBuilder.Find(nil).Statement
	queryExprInt := db.Dialector.Explain(queryExprIntStmt.SQL.String(), queryExprIntStmt.Vars...)

	queryExprBooleanBuilder := db.
		Session(&gorm.Session{DryRun: true}).
		Select("e.code_name, e.sort_order, e.display_format, 'boolean' AS entity_type, v.value, NULL AS metadata").
		Table("attributes AS e").
		Joins("JOIN customer_attribute_boolean AS v ON v.attribute_id = e.id").
		Where(fmt.Sprintf("v.customer_id = '%s'", ctx.ID.ValueOrZero()))

	if attrFilter != nil {
		queryExprBooleanBuilder = queryExprBooleanBuilder.
			// Where(attrFilter).
			Where(attrFilter.FilterQueryString()).
			Order(attrFilter.SortByQueryString())
	}
	queryExprBooleanStmt := queryExprBooleanBuilder.Find(nil).Statement
	queryExprBoolean := db.Dialector.Explain(queryExprBooleanStmt.SQL.String(), queryExprBooleanStmt.Vars...)

	queryExprTextBuilder := db.
		Session(&gorm.Session{DryRun: true}).
		Select("e.code_name, e.sort_order, e.display_format, 'text' AS entity_type, v.value, NULL AS metadata").
		Table("attributes AS e").
		Joins("JOIN customer_attribute_text AS v ON v.attribute_id = e.id").
		Where(fmt.Sprintf("v.customer_id = '%s'", ctx.ID.ValueOrZero()))

	if attrFilter != nil {
		queryExprTextBuilder = queryExprTextBuilder.
			// Where(attrFilter).
			Where(attrFilter.FilterQueryString()).
			Order(attrFilter.SortByQueryString())
	}
	queryExprTextStmt := queryExprTextBuilder.Find(nil).Statement
	queryExprText := db.Dialector.Explain(queryExprTextStmt.SQL.String(), queryExprTextStmt.Vars...)

	queryExprBlobBuilder := db.
		Session(&gorm.Session{DryRun: true}).
		Select("e.code_name, e.sort_order, e.display_format, 'blob' AS entity_type, v.value, JSON_OBJECT('name', v.name, 'type', v.type) AS metadata").
		Table("attributes AS e").
		Joins("JOIN customer_attribute_blob AS v ON v.attribute_id = e.id").
		Where(fmt.Sprintf("v.customer_id = '%s'", ctx.ID.ValueOrZero()))

	if attrFilter != nil {
		queryExprBlobBuilder = queryExprBlobBuilder.
			// Or(attrFilter).
			Where(attrFilter.FilterQueryString()).
			Order(attrFilter.SortByQueryString())
	}
	queryExprBlobStmt := queryExprBlobBuilder.Find(nil).Statement
	queryExprBlob := db.Dialector.Explain(queryExprBlobStmt.SQL.String(), queryExprBlobStmt.Vars...)

	rows, err := db.
		// Debug().
		// Session(&gorm.Session{WithConditions: true, DryRun: true}).
		Raw(fmt.Sprintf(
			"(%s) UNION (%s) UNION (%s) UNION (%s) UNION (%s) UNION (%s) UNION (%s) ORDER BY sort_order ASC",
			queryExprVarchar,
			queryExprDatetime,
			queryExprDecimal,
			queryExprInt,
			queryExprBoolean,
			queryExprText,
			queryExprBlob,
		)).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cAttr CustomerAttribute
		db.ScanRows(rows, &cAttr)

		if cAttr.CodeName.IsZero() {
			continue
		}
		codeName := cAttr.CodeName.ValueOrZero()
		list[codeName] = &cAttr
	}

	return list, err
}

// SetAttributes comment
func (ctx *Customer) SetAttributes(db *gorm.DB, infoValues map[string]interface{}) (err error) {
	customerID := ctx.ID.ValueOrZero()

	// jsonString, _ := json.Marshal(infoValues)
	// log.Println(string(jsonString[:]))

	if err = ctx.Transact(db, func(tx *gorm.DB) error {
		return UpsertCustomerAttributeValue(db, customerID, infoValues)
	}); err != nil {
		return err
	}
	return err
}

// SetAttributes comment
func (ctx *Customer) ValidatePassword(db *gorm.DB, inputPassword string) bool {
	type QueryResult struct {
		CodeName   null.String
		CustomerID null.String
		Value      null.String
	}
	var queryResult QueryResult

	queryExprBuilder := db.
		Session(&gorm.Session{DryRun: true}).
		Select("e.code_name, v.customer_id, v.value").
		Table("customer_attribute_varchar AS v").
		Joins("JOIN attributes AS e ON e.id = v.attribute_id").
		Where("e.code_name = 'password'").
		Where(fmt.Sprintf("v.customer_id = '%s'", ctx.ID.ValueOrZero()))

	queryExpr := queryExprBuilder.Find(nil).Statement.SQL.String()

	db.
		// Session(&gorm.Session{WithConditions: true, DryRun: true}).
		Raw(queryExpr).
		Scan(&queryResult)

	if !queryResult.CustomerID.IsZero() {
		passwordValue := queryResult.Value.ValueOrZero()
		// err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
		return strings.Compare(inputPassword, passwordValue) == 0
	}
	return false
}

// FindAllCustomerGroupIDs comment
func FindAllCustomerGroupIDs(db *gorm.DB, customerID string) (list []*map[string]string, err error) {
	var (
		customer Customer
		groups   []*Group
	)

	db.
		Where("id = ?", customerID).
		Find(&customer)

	db.
		Model(&customer).
		Select("id").
		Association("Groups").
		Find(&groups)

	for _, group := range groups {
		list = append(list, &map[string]string{
			"customerId": customerID,
			"groupId":    group.ID.ValueOrZero(),
		})
	}

	return list, err
}

// _acquireDefaultAttributes comment
func _acquireDefaultAttributes(db *gorm.DB, attrFilter *AttributeFilter) (map[string]*CustomerAttribute, error) {
	list := make(map[string]*CustomerAttribute)

	queryBuilder := db.
		// Session(&gorm.Session{WithConditions: true, DryRun: true}).
		Select("e.code_name, e.sort_order, e.display_format, e.entity_type, e.default_value AS value").
		Table("attributes AS e")

	if attrFilter != nil {
		attrFilter.SortBy["e.sort_order"] = "ASC"
		queryBuilder = queryBuilder.
			Where(attrFilter).
			Where(attrFilter.FilterQueryString()).
			Order(attrFilter.SortByQueryString())
	}

	rows, err := queryBuilder.Rows()
	if err != nil {
		// log.Println(err)
		return list, err
	}
	defer rows.Close()

	for rows.Next() {
		var cAttr CustomerAttribute
		db.ScanRows(rows, &cAttr)

		codeName := cAttr.CodeName.ValueOrZero()
		if cAttr.CodeName.IsZero() {
			continue
		}
		list[codeName] = &cAttr
	}

	return list, err
}

// // AfterCreate comment
// func (ctx *Customer) AfterCreate(scope *gorm.Scope) (err error) {
// 	var dbName null.String

// 	scope.DB().CommonDB().
// 		QueryRow("SELECT DATABASE() AS dbName").
// 		Scan(&dbName)

// 	dbNameString := dbName.ValueOrZero()

// 	re := regexp.MustCompile(`^system_customer_([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89AB][0-9a-fA-F]{3}-[0-9a-fA-F]{12})$`)
// 	organizationId := re.FindAllString(dbNameString, -1)

// 	log.Println(organizationId)
// 	return err
// }
