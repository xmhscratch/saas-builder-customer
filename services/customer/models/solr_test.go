package models_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"localdomain/customer/models"
)

// var prettyJSON bytes.Buffer
// jsonString, _ := json.Marshal(&info)
// _ = json.Indent(&prettyJSON, jsonString, "", "\t")
// fmt.Println(string(prettyJSON.Bytes()))

func TestSolrInitConfigSets(t *testing.T) {
	testInit(t)

	err := models.SolrInitConfigSets(cfg, "system_customer")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	t.Fail()
	os.Exit(1)
}

func TestSolrEnsureCollection(t *testing.T) {
	testInit(t)

	var err error
	time.Sleep(100 * time.Millisecond)

	go func() {
		err = models.SolrEnsureCollection(cfg, "system_customer", organizationID)
		if err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		err = models.SolrEnsureCollection(cfg, "system_customer", organizationID)
		if err != nil {
			fmt.Println(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	t.Fail()
	os.Exit(1)
}

func TestSolrDeleteCollection(t *testing.T) {
	testInit(t)

	err := models.SolrDeleteCollection(cfg, "system_customer", organizationID)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	t.Fail()
	os.Exit(1)
}

func TestSolrEraseConfigSets(t *testing.T) {
	testInit(t)

	var err error

	err = models.SolrEraseConfigSets(cfg, "system_customer_1ef19f81-a3a1-45a2-9203-b792abcddc52")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	// err = models.SolrEraseConfigSets(cfg, "system_customer")
	// if err != nil {
	// 	fmt.Println(err)
	// 	t.Fail()
	// }

	t.Fail()
	os.Exit(1)
}
