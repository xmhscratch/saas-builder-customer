package core

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/golang/groupcache"
)

// HTTPClientJar comment
var HTTPClientJar, _ = cookiejar.New(nil)

// HTTPClient comment
var HTTPClient = &http.Client{Jar: HTTPClientJar}

// RedisConnectionConfig ...
type RedisConnectionConfig struct {
	DBName   int    `json:"dbName"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

// SolrConnectionConfig ...
type SolrConnectionConfig struct {
	Host       string `json:"host"`
	Port       string `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Collection string `json:"collection"`
	MySQLHost  string `json:"mysqlHost"`
}

// // ElasticConnectionConfig ...
// type ElasticConnectionConfig struct {
// 	Addresses 				[]string		`json:"addresses"` 				// A list of Elasticsearch nodes to use.
// 	Username  				string 			`json:"username"`				// Username for HTTP Basic Authentication.
// 	Password  				string			`json:"password"`				// Password for HTTP Basic Authentication.

// 	CloudID 				string			`json:"cloudId"`				// Endpoint for the Elastic Service (https://elastic.co/cloud).
// 	APIKey  				string			`json:"apiKey"`					// Base64-encoded token for authorization; if set, overrides username and password.

// 	RetryOnStatus			[]int			`json:"retryOnStatus"`			// List of status codes for retry. Default: 502, 503, 504.
// 	DisableRetry			bool			`json:"disableRetry"`			// Default: false.
// 	EnableRetryOnTimeout 	bool			`json:"enableRetryOnTimeout"`	// Default: false.
// 	MaxRetries				int				`json:"maxRetries"`				// Default: 3.

// 	DiscoverNodesOnStart  	bool			`json:"discoverNodesOnStart"`	// Discover nodes when initializing the client. Default: false.
// 	DiscoverNodesInterval 	int				`json:"discoverNodesInterval"`	// Discover nodes periodically. Default: disabled.

// 	EnableMetrics     		bool			`json:"enableMetrics"`			// Enable the metrics collection.
// 	EnableDebugLogger 		bool			`json:"enableDebugLogger"`		// Enable the debug logging.
// }

// MySQLConnectionConfig ...
type MySQLConnectionConfig struct {
	User         string `json:"user"`
	Password     string `json:"password"`
	Host         string `json:"host"`
	Port         string `json:"port"`
	Type         string `json:"type"`
	MaxIdleConns int    `json:"maxIdleConns"`
	MaxOpenConns int    `json:"maxOpenConns"`
}

// ClusterHostNames comment
type ClusterHostNames struct {
	API         string `json:"api"`
	Integration string `json:"integration"`
	CustomerAPI string `json:"customer"`
}

// ConfigFile comment
type ConfigFile struct {
	Debug            bool                 `json:"debug"`
	Port             string               `json:"port"`
	HostName         string               `json:"hostName"`
	ClusterHostName  string               `json:"clusterHostName"`
	ClusterHostNames ClusterHostNames     `json:"clusterHostNames"`
	Solr             SolrConnectionConfig `json:"solr"`
	// Elastic ElasticConnectionConfig `json:"elastic"`
	MySQL                  MySQLConnectionConfig `json:"mysql"`
	Redis                  RedisConnectionConfig `json:"redis"`
	AMQPConnectionString   string                `json:"amqpConnectionString"`
	InfluxConnectionString string                `json:"influxConnectionString"`
	TmpPath                string                `json:"tmpPath"`
	DataPath               string                `json:"dataPath"`
}

// Config comment
type Config struct {
	ConfigFile
	AppDir     string
	GroupCache *groupcache.Group
}

// SystemAPI comment
type SystemAPI struct {
	cfg            *Config
	UserID         string
	OrganizationID string
	URL            *url.URL
	PostData       *url.Values
}
