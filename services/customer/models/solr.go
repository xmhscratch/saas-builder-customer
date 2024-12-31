package models

import (
	"encoding/json"
	"errors"
	"fmt"

	// "log"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
	"unsafe"

	"localdomain/customer/core"

	"github.com/vanng822/go-solr/solr"
)

type ListConfigSetsResponse struct {
	ConfigSets []string `json:"configSets"`
}

// type CoreStatusResponse struct {
// 	ResponseHeader struct {
// 		Status int `json:"status"`
// 		QTime  int `json:"QTime"`
// 	} `json:"responseHeader"`
// 	Status map[string]map[string]interface{} `json:"status"`
// }

type ListCollectionResponse struct {
	ResponseHeader struct {
		Status int `json:"status"`
		QTime  int `json:"QTime"`
	} `json:"responseHeader"`
	Collections []string `json:"collections"`
}

type RequestStatusResponse struct {
	ResponseHeader struct {
		Status int `json:"status"`
		QTime  int `json:"QTime"`
	} `json:"responseHeader"`
	Status struct {
		State string `json:"state"`
		Msg   string `json:"msg"`
	} `json:"status"`
}

func SolrGetRequest(cfg *core.Config, resourcePath string) (respBody []byte, err error) {
	resourceURLString := fmt.Sprintf(
		"http://%s:%s/solr/%s",
		cfg.Solr.Host,
		cfg.Solr.Port,
		strings.TrimPrefix(resourcePath, "/"),
	)

	d, _ := time.ParseDuration("30s")
	respBody, err = solr.HTTPGet(resourceURLString, nil, cfg.Solr.Username, cfg.Solr.Password, d)
	if !json.Valid(respBody) {
		err = errors.New("invalid solr GET response format")
	}

	return respBody, err
}

func SolrPostRequest(cfg *core.Config, resourcePath string, data *[]byte, headers [][]string) (respBody []byte, err error) {
	resourceURLString := fmt.Sprintf(
		"http://%s:%s/solr/%s",
		cfg.Solr.Host,
		cfg.Solr.Port,
		strings.TrimPrefix(resourcePath, "/"),
	)

	d, _ := time.ParseDuration("30s")
	respBody, err = solr.HTTPPost(resourceURLString, data, headers, cfg.Solr.Username, cfg.Solr.Password, d)
	if !json.Valid(respBody) {
		err = errors.New("invalid solr GET response format")
	}

	return respBody, err
}

func SolrInit(cfg *core.Config, configSetsName string, organizationID string) (err error) {
	var (
		hasConfigSets bool = false
		hasCollection bool = false
	)

	if hasConfigSets, err = SolrHasConfigSets(cfg, configSetsName); err != nil {
		return err
	}
	if !hasConfigSets {
		if err = SolrEnsureConfigSets(cfg, configSetsName); err != nil {
			return err
		}
	}

	if hasCollection, err = SolrHasCollection(cfg, configSetsName, organizationID); err != nil {
		return err
	}
	if !hasCollection {
		if err = SolrEnsureCollection(cfg, configSetsName, organizationID); err != nil {
			return err
		}
	}

	return err
}

func SolrHasConfigSets(cfg *core.Config, configSetsName string) (hasConfigSets bool, err error) {
	listConfigSetsResp := &ListConfigSetsResponse{}
	if listConfigSetsRespBody, err := SolrGetRequest(cfg, "/admin/configs?action=LIST&omitHeader=true"); err != nil {
		return false, err
	} else {
		if err := json.Unmarshal(listConfigSetsRespBody, &listConfigSetsResp); err != nil {
			return false, err
		}
	}

	if !core.IncludeInStringSlice(configSetsName, listConfigSetsResp.ConfigSets) {
		return false, err
	}
	return true, err
}

func SolrEnsureConfigSets(cfg *core.Config, configSetsName string) error {
	var (
		err               error
		configSetsZipData []byte
	)

	configSetsSourcePath := filepath.Join(cfg.AppDir, "./config/solr/conf")
	if byts, err := ZipSourceDirectory(configSetsSourcePath); err != nil {
		return err
	} else {
		configSetsZipData = byts.Bytes()
	}

	resourceURLString := fmt.Sprintf(
		"/admin/configs?action=UPLOAD&name=%s",
		configSetsName, // name
	)
	_, err = SolrPostRequest(cfg, resourceURLString, (*[]byte)(unsafe.Pointer(&configSetsZipData)), [][]string{
		{"Content-Type", "application/octet-stream"},
	})

	return err
}

func SolrEraseConfigSets(cfg *core.Config, configSetsName string) error {
	var (
		err  error
		mlck sync.Mutex
		wg   *sync.WaitGroup = &sync.WaitGroup{}
	)

	hasConfigSets := false
	if hasConfigSets, err = SolrHasConfigSets(cfg, configSetsName); err != nil {
		return err
	}
	if !hasConfigSets {
		return err
	}

	wg.Add(3)

	mlck.Lock()
	go func() {
		defer wg.Done()

		listCollectionResp := &ListCollectionResponse{}
		listCollectionRespBody, err := SolrGetRequest(cfg, "/admin/collections?action=LIST")
		if err != nil {
			mlck.Unlock()
			return
		} else {
			err = json.Unmarshal(listCollectionRespBody, &listCollectionResp)
			if err != nil {
				mlck.Unlock()
				return
			}
		}

		for _, collectionName := range listCollectionResp.Collections {
			isMatched, _ := regexp.MatchString(core.BuildString("^", configSetsName, "_(.+)$"), collectionName)
			if !isMatched {
				continue
			}
			organizationID := strings.Replace(collectionName, core.BuildString(configSetsName, "_"), "", 1)
			err = SolrDeleteCollection(cfg, configSetsName, organizationID)
			if err != nil {
				mlck.Unlock()
				return
			}
		}

		mlck.Unlock()
	}()

	mlck.Lock()
	go func() {
		defer wg.Done()

		resourceURLString := fmt.Sprintf(
			"/%s/update?commit=true",
			configSetsName,
		)
		postBytes := []byte("<delete><query>*.*</query></delete>")
		_, err = SolrPostRequest(cfg, resourceURLString, (*[]byte)(unsafe.Pointer(&postBytes)), [][]string{
			{"Content-Type", "text/xml"},
		})

		mlck.Unlock()
	}()

	mlck.Lock()
	go func() {
		defer wg.Done()

		resourceURLString := fmt.Sprintf(
			"/admin/configs?action=DELETE&name=%s&omitHeader=true",
			configSetsName,
		)
		_, err = SolrGetRequest(cfg, resourceURLString)

		mlck.Unlock()
	}()

	wg.Wait()
	return err
}

func SolrHasCollection(cfg *core.Config, configSetsName string, organizationID string) (hasCollection bool, err error) {
	hasConfigSets, err := SolrHasConfigSets(cfg, configSetsName)
	if err != nil {
		return false, err
	}
	if !hasConfigSets {
		return false, err
	}

	listCollectionResp := &ListCollectionResponse{}
	listCollectionRespBody, err := SolrGetRequest(cfg, "/admin/collections?action=LIST")
	if err != nil {
		return false, err
	} else {
		if err := json.Unmarshal(listCollectionRespBody, &listCollectionResp); err != nil {
			return false, err
		}
	}

	collectionName := core.BuildString(configSetsName, "_", organizationID)
	if !core.IncludeInStringSlice(collectionName, listCollectionResp.Collections) {
		return false, err
	}

	return true, err
}

func SolrEnsureCollection(cfg *core.Config, configSetsName string, organizationID string) error {
	var (
		err  error
		mlck sync.Mutex
		wg   *sync.WaitGroup = &sync.WaitGroup{}
	)

	if hasCollection, err := SolrHasCollection(cfg, configSetsName, organizationID); err != nil {
		return err
	} else {
		// log.Println(configSetsName, organizationID, hasCollection)
		if hasCollection {
			return err
		}
	}

	timestamp := (time.Now()).Unix()
	collectionName := core.BuildString(configSetsName, "_", organizationID)

	resourceURLString := fmt.Sprintf(
		"/admin/collections?_=%d&action=CREATE&autoAddReplicas=%t&collection.configName=%s&name=%s&maxShardsPerNode=%d&numShards=%d&replicationFactor=%d&router.name=compositeId&waitForFinalState=true&async=%s&wt=json",
		timestamp,
		false,          // autoAddReplicas
		configSetsName, // collection.configName
		collectionName, // name
		1,              // maxShardsPerNode
		1,              // numShards
		1,              // replicationFactor
		core.BuildString("create:", collectionName), // async
	)

	// log.Println(resourceURLString)

	_, err = SolrPostRequest(cfg, resourceURLString, nil, [][]string{
		{"Content-Type", "application/json;charset=utf-8"},
	})
	// log.Println(resourceURLString)
	if err != nil {
		return err
	}
	// log.Println(string(respBody[:]))

	wg.Add(3)

	mlck.Lock()
	go func() {
		defer wg.Done()
		_, err = SolrGetRequest(cfg, "/admin/collections?action=DELETESTATUS&flush=true&wt=json")
		mlck.Unlock()
	}()

	mlck.Lock()
	go func() {
		defer wg.Done()
		requestId := core.BuildString("create", ":", configSetsName, "_", organizationID)
		if ok := waitStateReady(cfg, "collections", requestId, 0); !ok {
			err = fmt.Errorf("cannot create collection with configsets: %s", configSetsName)
		}
		mlck.Unlock()
	}()

	mlck.Lock()
	go func() {
		defer wg.Done()
		_, err = SolrGetRequest(cfg, "/admin/collections?action=DELETESTATUS&flush=true&wt=json")
		mlck.Unlock()
	}()

	wg.Wait()
	return err
}

func SolrDeleteCollection(cfg *core.Config, configSetsName string, organizationID string) error {
	var err error

	hasConfigSets := false
	if hasConfigSets, err = SolrHasConfigSets(cfg, configSetsName); err != nil {
		return err
	}
	if !hasConfigSets {
		return err
	}

	collectionName := core.BuildString(configSetsName, "_", organizationID)

	resourceURLString := fmt.Sprintf(
		"/admin/collections?action=DELETE&name=%s&wt=json",
		collectionName, // name
	)
	_, err = SolrGetRequest(cfg, resourceURLString)
	// if err != nil {
	// 	// log.Fatalln(err)
	// } else {
	// 	log.Println(string(respBody[:]))
	// }

	return err
}

func SolrRequestState(cfg *core.Config, resourceName string, requestId string) (requestState string, err error) {
	requestStatusResp := &RequestStatusResponse{}
	requestStatusURLString := fmt.Sprintf(
		"/admin/%s?action=REQUESTSTATUS&requestid=%s&wt=json",
		resourceName,
		requestId,
	)
	if requestStatusRespBody, err := SolrGetRequest(cfg, requestStatusURLString); err != nil {
		return "failed", err
	} else {
		if err := json.Unmarshal(requestStatusRespBody, &requestStatusResp); err != nil {
			return "failed", err
		}
	}
	// log.Println(requestStatusResp)
	return requestStatusResp.Status.State, err
}

func waitStateReady(cfg *core.Config, resourceName string, requestId string, retried int) (ok bool) {
	requestState, err := SolrRequestState(cfg, resourceName, requestId)
	if err != nil {
		return ok
	}
	// log.Println(requestState, retried)
	if retried > 15*1000/250 {
		return ok
	}
	if requestState == "notfound" {
		time.Sleep(250 * time.Millisecond)
		go waitStateReady(cfg, resourceName, requestId, retried+1)
		return ok
	}
	if requestState == "running" || requestState == "submitted" {
		time.Sleep(50 * time.Millisecond)
		go waitStateReady(cfg, resourceName, requestId, 0)
		return ok
	}
	// completed
	return true
}
