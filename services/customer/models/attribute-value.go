package models

import (
	// "bytes"
	// "database/sql"
	// "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	// "log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"localdomain/customer/core"

	"github.com/araddon/dateparse"
	null "gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

// CustomerAttributeValue struct
type CustomerAttributeValue struct {
	CustomerAttributeIndex
	Value    null.String `gorm:"column:value;" json:"value"`
	Metadata null.String `gorm:"column:metadata;" json:"metadata"`
}

// CustomerAttribute model
type CustomerAttribute struct {
	CustomerAttributeValue
	CodeName      null.String `gorm:"column:code_name;type:string;" sql:"type:varchar(100)" json:"codeName"`
	EntityType    null.String `gorm:"column:entity_type;type:string;" sql:"type:enum('varchar','datetime','decimal','int','boolean','text','blob')" json:"entityType"`
	DisplayFormat null.String `gorm:"column:display_format;type:string;" sql:"type:varchar(12)" json:"displayFormat"`
}

// CustomerAttributeDictionary model
type CustomerAttributeDictionary struct {
	AttributeID   string      `json:"attributeId"`
	CustomerID    string      `json:"customerId"`
	EntityType    string      `json:"entityType"`
	CodeName      string      `json:"codeName"`
	DisplayFormat string      `json:"displayFormat"`
	Value         interface{} `json:"value"`
	Value2        interface{} `json:"value2"`
}

// MarshalJSON implements json.Marshaler.
func (ctx CustomerAttribute) MarshalJSON() ([]byte, error) {
	var (
		err error
		v   interface{}
		v2  interface{}
	)

	if !ctx.Value.IsZero() {
		switch ctx.EntityType.ValueOrZero() {
		case "blob":
			{
				// var vd []byte
				// if vd, err = GZipDecompress([]byte(ctx.Value.ValueOrZero())); err != nil {
				// 	vd = nil
				// }
				vStr := []uint8(ctx.Value.ValueOrZero())

				var byts []interface{}
				for _, val := range vStr {
					byts = append(byts, val)
				}

				vStrArr := strings.Join(strings.Fields(fmt.Sprintf("%d", byts)), ",")
				if err := json.Unmarshal([]byte(vStrArr), &v); err != nil {
					v = nil
				}
				if err := json.Unmarshal([]byte(ctx.Metadata.ValueOrZero()), &v2); err != nil {
					v2 = nil
				}

				break
			}
		case "datetime":
			{
				vStr, err := ParseString(ctx.Value.ValueOrZero())
				if err != nil || vStr == "" {
					v = nil
				}
				vSs := strings.Split(vStr, " - ")
				if len(vSs) == 1 {
					if v, err = dateparse.ParseAny(vSs[0]); err != nil {
						v = nil
					}
					v2 = nil
				} else if len(vSs) == 2 {
					if v, err = dateparse.ParseAny(vSs[0]); err != nil {
						v = nil
					}
					if v2, err = dateparse.ParseAny(vSs[1]); err != nil {
						v2 = nil
					}
				} else {
					v = nil
					v2 = nil
				}
				break
			}
		case "decimal":
			{
				v, err = ParseFloat(ctx.Value.ValueOrZero())
				if err != nil {
					v = float64(0.00)
				}
				break
			}
		case "int":
			{
				v, err = ParseInt(ctx.Value.ValueOrZero())
				if err != nil {
					v = int64(0)
				}
				break
			}
		case "boolean":
			{
				vStr, err := ParseString(ctx.Value.ValueOrZero())
				if err != nil {
					v = bool(false)
				} else {
					v, err = strconv.ParseBool(vStr)
					if err != nil {
						v = bool(false)
					}
				}
				break
			}
		case "text":
			{
				v, err = ParseString(ctx.Value.ValueOrZero())
				if err != nil {
					v = ""
				}
				break
			}
		case "varchar":
			{
				v, err = ParseString(ctx.Value.ValueOrZero())
				if err != nil {
					v = ""
				}
				break
			}
		default:
			{
				v = ctx.Value.ValueOrZero()
				break
			}
		}
	} else {
		v = nil
	}

	if err != nil {
		return []byte("null"), err
	}

	result := map[string]interface{}{
		"attributeId":   ctx.AttributeID.ValueOrZero(),
		"customerId":    ctx.CustomerID.ValueOrZero(),
		"entityType":    ctx.EntityType.ValueOrZero(),
		"codeName":      ctx.CodeName.ValueOrZero(),
		"displayFormat": ctx.DisplayFormat.ValueOrZero(),
		"value":         v,
		"value2":        v2,
	}

	return json.Marshal(result)
}

// UpsertCustomerAttributeValue comment
func UpsertCustomerAttributeValue(db *gorm.DB, customerID string, infoValues map[string]interface{}) error {
	var err error

	for codeName, infoValue := range infoValues {
		var (
			attr        Attribute
			attributeID string
			entityType  string
		)

		// if infoValue == nil {
		// 	continue
		// }

		if err = db.
			// Debug().
			Model(Attribute{}).
			Where("`code_name` = ?", codeName).
			Find(&attr).
			Error; err != nil {
			continue
		}

		attributeID = attr.ID.ValueOrZero()

		if attr.CodeName.IsZero() {
			err = errors.New(core.BuildString("attribute code name -", codeName, "- is empty"))
			continue
		}
		codeName = attr.CodeName.ValueOrZero()

		if attr.EntityType.IsZero() {
			err = errors.New(core.BuildString("attribute entity type -", codeName, "- is empty"))
			continue
		}
		entityType = attr.EntityType.ValueOrZero()
		if entityType == "" {
			entityType = "varchar"
		}

		isReadOnly := attr.IsReadOnly.ValueOrZero()

		rawValue := infoValue
		// rawValue := infoValue.(interface{})
		isMap := reflect.ValueOf(rawValue).Kind() == reflect.Map
		isJSON := attr.DisplayFormat.ValueOrZero() == "json"
		if isMap && isJSON {
			formattedValue, err := json.Marshal(rawValue)
			if err != nil {
				continue
			}
			rawValue = string(formattedValue[:])
		}

		customerAttributeIndex := CustomerAttributeIndex{
			CustomerID:  null.StringFrom(customerID),
			AttributeID: null.StringFrom(attributeID),
		}

		queryResult := db.
			// Debug().
			Where(customerAttributeIndex)

		// log.Println(codeName, entityType, rawValue)

		switch entityType {
		case "int":
			{
				var (
					cAttr    CustomerAttributeInt
					newValue int64
				)

				if newValue, err = ParseInt(rawValue); err != nil {
					continue
				}

				if !(newValue == newValue) && !(newValue > newValue) && !(newValue < newValue) {
					queryResult = queryResult.Delete(&cAttr)
					if queryResult.Error != nil {
						err = queryResult.Error
					}
					continue
				}

				newAtrrs := CustomerAttributeInt{
					CustomerAttributeIndex: customerAttributeIndex,
					Value:                  null.IntFrom(newValue),
				}

				if isReadOnly {
					queryResult = queryResult.Attrs(newAtrrs)
				} else {
					queryResult = queryResult.Assign(newAtrrs)
				}
				queryResult = queryResult.FirstOrCreate(&cAttr)
				if queryResult.Error != nil {
					err = queryResult.Error
					continue
				}

				break
			}
		case "boolean":
			{
				var (
					cAttr    CustomerAttributeBoolean
					newValue bool
				)

				if newValue, err = ParseBool(rawValue); err != nil {
					continue
				}

				newAtrrs := CustomerAttributeBoolean{
					CustomerAttributeIndex: customerAttributeIndex,
					Value:                  null.BoolFrom(newValue),
				}

				if isReadOnly {
					queryResult = queryResult.Attrs(newAtrrs)
				} else {
					queryResult = queryResult.Assign(newAtrrs)
				}
				queryResult = queryResult.FirstOrCreate(&cAttr)
				if queryResult.Error != nil {
					err = queryResult.Error
					continue
				}

				break
			}
		case "decimal":
			{
				var (
					cAttr    CustomerAttributeDecimal
					newValue float64
				)

				if newValue, err = ParseFloat(rawValue); err != nil {
					continue
				}

				if !(newValue == newValue) && !(newValue > newValue) && !(newValue < newValue) {
					queryResult = queryResult.Delete(&cAttr)
					if queryResult.Error != nil {
						err = queryResult.Error
					}
					continue
				}

				newAtrrs := CustomerAttributeDecimal{
					CustomerAttributeIndex: customerAttributeIndex,
					Value:                  null.FloatFrom(newValue),
				}

				if isReadOnly {
					queryResult = queryResult.Attrs(newAtrrs)
				} else {
					queryResult = queryResult.Assign(newAtrrs)
				}
				queryResult = queryResult.FirstOrCreate(&cAttr)
				if queryResult.Error != nil {
					err = queryResult.Error
					continue
				}

				break
			}
		case "datetime":
			{
				var (
					cAttr       CustomerAttributeDateTime
					newValueRaw string
					newValue1   *time.Time
					newValue2   *time.Time
				)

				if newValueRaw, err = ParseString(rawValue); err != nil {
					continue
				}

				if newValueRaw == "" {
					queryResult = queryResult.Delete(&cAttr)
					if queryResult.Error != nil {
						err = queryResult.Error
					}
					continue
				}

				vSs := strings.Split(newValueRaw, " - ")
				if len(vSs) == 1 {
					if nv, err := dateparse.ParseAny(vSs[0]); err != nil {
						newValue1 = nil
					} else {
						newValue1 = &nv
					}
					newValue2 = nil
				} else if len(vSs) == 2 {
					if nv, err := dateparse.ParseAny(vSs[0]); err != nil {
						newValue1 = nil
					} else {
						newValue1 = &nv
					}
					if nv, err := dateparse.ParseAny(vSs[1]); err != nil {
						newValue2 = nil
					} else {
						newValue2 = &nv
					}
				} else {
					newValue1 = nil
					newValue2 = nil
				}

				newAtrrs := CustomerAttributeDateTime{
					CustomerAttributeIndex: customerAttributeIndex,
					Value:                  null.TimeFromPtr(newValue1),
					Value2:                 null.TimeFromPtr(newValue2),
				}
				if newValue2 == nil {
					newAtrrs.Value2 = gorm.Expr("NULL")
				}

				if isReadOnly {
					queryResult = queryResult.Attrs(newAtrrs)
				} else {
					queryResult = queryResult.Assign(newAtrrs)
				}
				queryResult = queryResult.FirstOrCreate(&cAttr)
				if queryResult.Error != nil {
					err = queryResult.Error
					continue
				}

				break
			}
		case "text":
			{
				handleRawValue := func(rawValue []uint8) error {
					var (
						cAttr    CustomerAttributeText
						newValue string
					)

					if newValue, err = ParseString(rawValue); err != nil {
						return err
					}

					if newValue == "" {
						queryResult = queryResult.Delete(&cAttr)
						if queryResult.Error != nil {
							err = queryResult.Error
						}
						return err
					}

					newAtrrs := CustomerAttributeText{
						CustomerAttributeIndex: customerAttributeIndex,
						Value:                  null.StringFrom(newValue),
					}

					if isReadOnly {
						queryResult = queryResult.Attrs(newAtrrs)
					} else {
						queryResult = queryResult.Assign(newAtrrs)
					}
					queryResult = queryResult.FirstOrCreate(&cAttr)
					if queryResult.Error != nil {
						return queryResult.Error
					}

					return nil
				}

				switch reflect.TypeOf(rawValue) {
				case reflect.TypeOf([][]uint8{}):
					{
						var vals []string

						for _, val := range rawValue.([][]uint8) {
							if newVal, err := ParseString(val); err != nil {
								continue
							} else {
								vals = append(vals, newVal)
							}
						}
						var byts []uint8
						if byts, err = json.Marshal(vals); err != nil {
							continue
						} else {
							if err := handleRawValue(byts); err != nil {
								continue
							}
						}
						break
					}
				case reflect.TypeOf([]uint8{}):
					{
						if err := handleRawValue(rawValue.([]uint8)); err != nil {
							continue
						}
						break
					}
				case reflect.TypeOf(string("")):
					{
						if err := handleRawValue([]uint8(rawValue.(string))); err != nil {
							continue
						}
						break
					}
				default:
					{
						break
					}
				}
				break
			}
		case "varchar":
			{
				var (
					cAttr    CustomerAttributeVarchar
					newValue string
				)

				if newValue, err = ParseString(rawValue); err != nil {
					continue
				}

				if newValue == "" {
					queryResult = queryResult.Delete(&cAttr)
					if queryResult.Error != nil {
						err = queryResult.Error
					}
					continue
				}

				newAtrrs := CustomerAttributeVarchar{
					CustomerAttributeIndex: customerAttributeIndex,
					Value:                  null.StringFrom(newValue),
				}

				if isReadOnly {
					queryResult = queryResult.Attrs(newAtrrs)
				} else {
					queryResult = queryResult.Assign(newAtrrs)
				}
				queryResult = queryResult.FirstOrCreate(&cAttr)
				if queryResult.Error != nil {
					err = queryResult.Error
					continue
				}

				break
			}
		case "blob":
			{
				var (
					cAttr        CustomerAttributeBlob
					newValue     []byte
					newName      string
					newType      string
					newValueRaw  string
					newValueJSON map[string]interface{}
				)

				if newValueRaw, err = ParseString(rawValue); err != nil {
					continue
				}

				if newValueRaw == "" {
					queryResult = queryResult.Delete(&cAttr)
					if queryResult.Error != nil {
						err = queryResult.Error
					}
					continue
				}

				if err = json.Unmarshal([]byte(newValueRaw), &newValueJSON); err != nil {
					continue
				}

				newName = newValueJSON["name"].(string)
				newType = newValueJSON["type"].(string)

				var byts []uint8
				for _, val := range newValueJSON["data"].([]interface{}) {
					byts = append(byts, uint8(val.(float64)))
				}

				newValue = []byte(byts)
				// if newValue, err = GZipCompress([]byte(byts)); err != nil {
				// 	continue
				// }

				newAtrrs := CustomerAttributeBlob{
					CustomerAttributeIndex: customerAttributeIndex,
					Value:                  newValue,
					Name:                   null.StringFrom(newName),
					Type:                   null.StringFrom(newType),
				}

				if isReadOnly {
					queryResult = queryResult.Attrs(newAtrrs)
				} else {
					queryResult = queryResult.Assign(newAtrrs)
				}
				queryResult = queryResult.FirstOrCreate(&cAttr)
				if queryResult.Error != nil {
					err = queryResult.Error
					continue
				}

				break
			}
		default:
			{
				break
			}
		}
	}

	return err
}

// DeleteCustomerAttributeValue comment
func DeleteCustomerAttributeValue(db *gorm.DB, customerID string, codeName string) error {
	var (
		err         error
		attr        *Attribute
		entityType  string
		attributeID string
		isReadOnly  bool = true
	)

	if attr, err = FindOneAttribute(db, codeName); err != nil {
		return err
	}

	if attr.ID.IsZero() {
		return err
	}
	attributeID = attr.ID.ValueOrZero()

	if attr.EntityType.IsZero() {
		return err
	}
	entityType = attr.EntityType.ValueOrZero()

	if attr.IsReadOnly.IsZero() {
		isReadOnly = true
	}
	isReadOnly = attr.IsReadOnly.ValueOrZero()
	if isReadOnly {
		return err
	}

	if entityType == "" {
		entityType = "varchar"
	}

	attrIndex := CustomerAttributeIndex{
		CustomerID:  null.StringFrom(customerID),
		AttributeID: null.StringFrom(attributeID),
	}

	switch entityType {
	case "int":
		{
			err = db.
				// Debug().
				Model(CustomerAttributeInt{}).
				Where(attrIndex).
				Delete(&CustomerAttributeInt{}).
				Error
			break
		}
	case "boolean":
		{
			err = db.
				// Debug().
				Model(CustomerAttributeBoolean{}).
				Where(attrIndex).
				Delete(&CustomerAttributeBoolean{}).
				Error
			break
		}
	case "decimal":
		{
			err = db.
				// Debug().
				Model(CustomerAttributeDecimal{}).
				Where(attrIndex).
				Delete(&CustomerAttributeDecimal{}).
				Error
			break
		}
	case "datetime":
		{
			err = db.
				// Debug().
				Model(CustomerAttributeDateTime{}).
				Where(attrIndex).
				Delete(&CustomerAttributeDateTime{}).
				Error
			break
		}
	case "text":
		{
			err = db.
				// Debug().
				Model(CustomerAttributeText{}).
				Where(attrIndex).
				Delete(&CustomerAttributeText{}).
				Error
			break
		}
	case "varchar":
		{
			err = db.
				// Debug().
				Model(CustomerAttributeVarchar{}).
				Where(attrIndex).
				Delete(&CustomerAttributeVarchar{}).
				Error
			break
		}
	case "blob":
	default:
		{
			err = db.
				// Debug().
				Model(CustomerAttributeBlob{}).
				Where(attrIndex).
				Delete(&CustomerAttributeBlob{}).
				Error
			break
		}
	}

	return err
}
