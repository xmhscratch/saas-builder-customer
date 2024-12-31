package routers_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"localdomain/customer/core"
	"localdomain/customer/models"

	"github.com/vanng822/go-solr/solr"
	"gorm.io/gorm"
)

const organizationID = "1ef19f81-a3a1-45a2-9203-b792abcddc52"

var (
	cfg           *core.Config
	customerDb    string
	db            *gorm.DB
	solrInterface *solr.SolrInterface
)

func testInit(t *testing.T) {
	var err error

	cfg, err = core.NewConfig("/home/web/repos/customer/")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	customerDb = fmt.Sprintf("system_customer_%s", organizationID)
	db, err = core.NewDatabase(cfg, customerDb)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	solrInterface, err = solr.NewSolrInterface(
		fmt.Sprintf("http://%s:%s/solr", "solr_solr_svc", cfg.Solr.Port),
		cfg.Solr.Collection,
	)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestMigrateDb(t *testing.T) {
	testInit(t)

	var err error

	organizationID := "98374495-13f1-41ab-9bd9-aa3e3c1b599d"
	newDbName := fmt.Sprintf("system_customer_%s", organizationID)
	db, err = core.NewDatabase(cfg, newDbName)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	tx := db.Begin()

	err = tx.
		AutoMigrate(
			models.Customer{},                  // customers
			models.Attribute{},                 // attributes
			models.CustomerAttributeBlob{},     // customer_attribute_blob
			models.CustomerAttributeBoolean{},  // customer_attribute_boolean
			models.CustomerAttributeDateTime{}, // customer_attribute_datetime
			models.CustomerAttributeDecimal{},  // customer_attribute_decimal
			models.CustomerAttributeInt{},      // customer_attribute_int
			models.CustomerAttributeText{},     // customer_attribute_text
			models.CustomerAttributeVarchar{},  // customer_attribute_varchar
			models.CustomerEvent{},             // customer_events
			models.CustomerEventData{},         // customer_event_series
			models.Enumeration{},               // enumerations
			models.EnumerationValue{},          // enumeration_values
			models.GroupCategory{},             // group_categories
			models.Group{},                     // groups
			models.GroupCustomerLinker{},       // group_customer_linker
		)

	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	if err := tx.Commit().Error; err != nil {
		log.Println(err)
		t.Fail()
		return
	}

	log.Println(err)
	t.Fail()

	os.Exit(1)
}
