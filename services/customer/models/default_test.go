package models_test

import (
	"fmt"
	"testing"

	"localdomain/customer/core"

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
	cfg.Solr.Host = "127.0.0.1"

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
