// Package config defines the environment variable and command-line flags
package config

import (
	"sync"

	"github.com/companieshouse/gofigure"
)

var cfg *Config
var mtx sync.Mutex

// Config defines the configuration options for this service.
type Config struct {
	BindAddr                       string   `env:"BIND_ADDR"                         flag:"bind-addr"                           flagDesc:"Bind address"`
	MongoDBURL                     string   `env:"MONGODB_URL"                       flag:"mongodb-url"                         flagDesc:"MongoDB server URL"`
	MongoAuthcodeDatabase          string   `env:"MONGO_AUTHCODE_DATABASE"           flag:"mongodb-authcode-database"           flagDesc:"MongoDB database for auth code data"`
	MongoAuthCodeCollection        string   `env:"MONGO_AUTHCODE_COLLECTION"         flag:"mongodb-authcode-collection"         flagDesc:"The name of the mongodb auth code collection"`
	MongoAuthcodeRequestDatabase   string   `env:"MONGO_AUTHCODE_REQUEST_DATABASE"   flag:"mongodb-authcode-request-database"   flagDesc:"MongoDB database for auth code request data"`
	MongoAuthCodeRequestCollection string   `env:"MONGO_AUTHCODE_REQUEST_COLLECTION" flag:"mongodb-authcode-request-collection" flagDesc:"The name of the mongodb auth code request collection"`
	OracleQueryAPIURL              string   `env:"ORACLE_QUERY_API_URL"              flag:"oracle-query-api-url"                flagDesc:"Oracle Query API URL"`
	QueueAPILocalURL               string   `env:"QUEUE_API_LOCAL_URL"               flag:"queue-api-local-url"                 flagDesc:"Queue API Local URL"`
	BrokerAddr                     []string `env:"KAFKA_BROKER_ADDR"                 flag:"broker-addr"                         flagDesc:"Kafka broker address"`
	SchemaRegistryURL              string   `env:"SCHEMA_REGISTRY_URL"               flag:"schema-registry-url"                 flagDesc:"Schema registry url"`
	CHSURL                         string   `env:"CHS_URL"                           flag:"chs-url"                             flagDesc:"CHS URL"`
}

// Get returns a pointer to a Config instance populated with values from environment or command-line flags
func Get() (*Config, error) {
	mtx.Lock()
	defer mtx.Unlock()

	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{}

	err := gofigure.Gofigure(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
