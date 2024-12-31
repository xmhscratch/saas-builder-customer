package routers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"localdomain/customer/core"
	"localdomain/customer/models"

	"github.com/vanng822/go-solr/solr"
	null "gopkg.in/guregu/null.v4"
)

const customerID = "1PEHGKQlfEe7lDlCZUMKSgPlisH"

func TestFindAllCustomer(t *testing.T) {
	testInit(t)

	customer, err := models.FindOneCustomer(db, customerID)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	info, err := customer.FetchDetailInfo(db, nil)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	var prettyJSON bytes.Buffer
	jsonString, _ := json.Marshal(&info)
	_ = json.Indent(&prettyJSON, jsonString, "", "\t")
	fmt.Println(string(prettyJSON.Bytes()))

	t.Fail()
}

func TestSearchCustomer(t *testing.T) {
	var (
		results []*string
	)

	testInit(t)

	q, err := models.BuildFilterAttributeTextQuery(db, "9")

	query := solr.NewQuery()
	query.Q(q)
	query.FilterQuery(core.BuildString("organizationId:", organizationID))

	// fq := ginCtx.Query("fq")
	// if fq != "" {
	// 	query.FilterQuery(fq)
	// }

	query.FieldList("id")

	query.Start(0)
	query.Rows(200)

	s := solrInterface.Search(query)

	r, err := s.Result(nil)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	for _, searchResult := range r.Results.Docs {
		customerID := string(searchResult.Get("id").(string))
		results = append(results, &customerID)
	}

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	var prettyJSON bytes.Buffer
	jsonString, _ := json.Marshal(&results)
	_ = json.Indent(&prettyJSON, jsonString, "", "\t")
	fmt.Println(len(results))
	fmt.Println(string(prettyJSON.Bytes()))

	os.Exit(1)
}

func TestGetAttributes(t *testing.T) {
	var (
		err     error
		attrs   []*models.Attribute
		results []*map[string]interface{}
	)

	testInit(t)

	attrFilter := &models.AttributeFilter{
		// CodeNames: map[string]string{
		// 	"0": "firstName",
		// },
		// WithAssociations: true,
		EntityTypes: map[string]string{
			"0": "varchar",
		},
		IsFilterable: null.BoolFrom(true),
	}

	attrs, err = models.FindAllAttributes(db, attrFilter, nil)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	for _, attr := range attrs {
		info, err := attr.FetchDetailInfo(db)
		if err != nil {
			break
		}
		results = append(results, &info)
	}
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	var prettyJSON bytes.Buffer
	jsonString, _ := json.Marshal(&results)
	_ = json.Indent(&prettyJSON, jsonString, "", "\t")

	fmt.Println(len(results))
	fmt.Println(string(prettyJSON.Bytes()))

	os.Exit(1)
}

func TestAddCustomer(t *testing.T) {
	testInit(t)

	var err error

	customer := &models.Customer{
		EmailAddress:  null.StringFrom("admin123@localhost.com"),
		InitialPoints: null.IntFrom(1),
		SyncedAt:      null.TimeFrom(time.Now()),
		TimestampModel: models.TimestampModel{
			CreatedAt: null.TimeFrom(time.Now()),
			UpdatedAt: null.TimeFrom(time.Now()),
			DeletedAt: null.TimeFrom(time.Now()),
		},
		// DefaultGormModel: models.DefaultGormModel{
		// 	ID: null.StringFrom("1bNcGadyDNsmrA7YbtgN8DUuMIj"),
		// },
	}

	err = customer.Register(db, "", true)
	if err != nil {
		t.Errorf("failed: %v", err)
		return
	}

	fmt.Println(customer)
	os.Exit(1)
}
