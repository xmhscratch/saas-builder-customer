package models

import (
	"encoding/json"
	"fmt"

	// "log"
	"strings"

	// "github.com/elastic/go-elasticsearch/v8"
	// "github.com/elastic/go-elasticsearch/v8/esapi"

	"localdomain/customer/core"

	"github.com/vanng822/go-solr/solr"
)

var solrInterface *solr.SolrInterface

type CustomerSearch struct {
	OrganizationID string
	Solr           *solr.SolrInterface
	Config         *core.Config
}

func NewCustomerSearch(cfg *core.Config, organizationID string, isForce bool) (*CustomerSearch, error) {
	var (
		searcher *CustomerSearch
		err      error
	)

	if solrInterface == nil {
		solrInterface, err = solr.NewSolrInterface(
			fmt.Sprintf("http://%s:%s/solr", cfg.Solr.Host, cfg.Solr.Port),
			"",
		)
		if err != nil {
			return nil, err
		}
		solrInterface.SetBasicAuth(cfg.Solr.Username, cfg.Solr.Password)
		if err := SolrInit(cfg, "system_customer", organizationID); err != nil {
			return nil, err
		}
	}

	// elasticConfig := elasticsearch.Config{
	// 	Addresses: cfg.Elastic.Addresses,
	// 	Username: cfg.Elastic.Username,
	// 	Password: cfg.Elastic.Password,
	// }

	// es, err := elasticsearch.NewClient(elasticConfig)
	// if err != nil {
	// 	return nil, err
	// }
	// res, err := es.Info()
	// if err != nil {
	// 	return nil, err
	// }
	// // Check response status
	// if res.IsError() {
	// 	log.Fatalf("Error: %s", res.String())
	// 	return nil, err
	// }

	searcher = &CustomerSearch{
		OrganizationID: organizationID,
		Solr:           solrInterface,
		Config:         cfg,
	}

	if err := searcher._syncSearchImport(isForce); err != nil {
		return searcher, nil
	}
	return searcher, nil
}

func (ctx *CustomerSearch) SearchIDs(q string, aq string, fq string, pageOffset int, pageSize int) ([]string, int, error) {
	var (
		customerIDs []string = []string{}
		total       int      = 0
		err         error    = nil
	)

	query := solr.NewQuery()
	queryStacks := []string{}
	if len(q) > 0 {
		queryStacks = append(queryStacks, core.BuildString("(", q, ")"))
	}
	if len(aq) > 0 {
		queryStacks = append(queryStacks, core.BuildString("(", aq, ")"))
	}
	queryString := strings.Join(queryStacks, " AND ")
	query.Q(queryString)

	query.FilterQuery(core.BuildString("organizationId:", ctx.OrganizationID))
	query.FilterQuery(core.BuildString("-deletedAt:", "[* TO *]"))
	if fq != "" {
		query.FilterQuery(fq)
	}

	query.FieldList("id")

	query.Start(pageOffset)
	query.Rows(pageSize)

	s := ctx.Solr.Search(query)
	r, err := s.Result(nil)
	if err != nil {
		return customerIDs, total, err
	}
	total = int(r.Results.NumFound)

	for _, searchResult := range r.Results.Docs {
		customerID, err := ParseString(searchResult.Get("id"))
		if err != nil {
			return customerIDs, total, err
		}
		if customerID == "" {
			total--
			continue
		}
		customerIDs = append(customerIDs, customerID)
	}
	return customerIDs, total, err
}

// _syncSearchImport
func (ctx *CustomerSearch) _syncSearchImport(isForce bool) (err error) {
	err = SolrEnsureCollection(ctx.Config, "system_customer", ctx.OrganizationID)
	if err != nil {
		return err
	}
	collectionName := core.BuildString("system_customer", "_", ctx.OrganizationID)
	ctx.Solr.SetCore(collectionName)

	resourcePath := fmt.Sprintf(
		"/%s/select?q=*:*&facet=true&rows=0&hostname=%s&user=%s&password=%s&organizationId=%s",
		collectionName,
		ctx.Config.Solr.MySQLHost,
		ctx.Config.MySQL.User,
		ctx.Config.MySQL.Password,
		ctx.OrganizationID,
	)
	// log.Println(resourcePath)

	rs, err := SolrGetRequest(ctx.Config, resourcePath)
	if err != nil {
		return err
	}
	// log.Println(string(rs[:]), err)

	var results map[string]interface{}
	if err := json.Unmarshal(rs, &results); err != nil {
		return err
	}
	resp := results["response"].(map[string]interface{})
	numFound := float64(resp["numFound"].(float64))

	// log.Println(resp, numFound)

	if !isForce {
		if numFound == 0 {
			// log.Println("_syncSearchFullImport")
			err = ctx._syncSearchFullImport()
		} else {
			// log.Println("_syncSearchDeltaImport")
			err = ctx._syncSearchDeltaImport()
		}
	} else {
		// log.Println("_syncSearchFullImport")
		err = ctx._syncSearchFullImport()
	}

	return err
}

func (ctx *CustomerSearch) _syncSearchDeltaImport() (err error) {
	resourcePath := fmt.Sprintf(
		"/%s/dataimport?wt=json&command=delta-import&clean=true&commit=true&hostname=%s&user=%s&password=%s&organizationId=%s&synchronous=%t",
		core.BuildString("system_customer", "_", ctx.OrganizationID),
		ctx.Config.Solr.MySQLHost,
		ctx.Config.MySQL.User,
		ctx.Config.MySQL.Password,
		ctx.OrganizationID,
		true,
	)
	// log.Println(resourcePath)
	_, err = SolrGetRequest(ctx.Config, resourcePath)
	return err
}

// _syncSearchFullImport
func (ctx *CustomerSearch) _syncSearchFullImport() (err error) {
	resourcePath := fmt.Sprintf(
		"/%s/dataimport?wt=json&command=full-import&clean=true&commit=true&hostname=%s&user=%s&password=%s&organizationId=%s&synchronous=%t",
		core.BuildString("system_customer", "_", ctx.OrganizationID),
		ctx.Config.Solr.MySQLHost,
		ctx.Config.MySQL.User,
		ctx.Config.MySQL.Password,
		ctx.OrganizationID,
		true,
	)
	// log.Println(resourcePath)
	_, err = SolrGetRequest(ctx.Config, resourcePath)
	// log.Println(string(rs[:]), err)
	return err
}
