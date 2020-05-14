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
	BindAddr                string `env:"BIND_ADDR"              flag:"bind-addr"          flagDesc:"Bind address"`
	MongoDBURL              string `env:"MONGODB_URL"            flag:"mongodb-url"        flagDesc:"MongoDB server URL"`
	MongoAccountDatabase    string `env:"MONGO_AUTHCODE_DATABASE"   flag:"mongodb-authcode-database"   flagDesc:"MongoDB database for auth code data"`
	MongoAuthCodeCollection string `env:"MONGO_AUTHCODE_COLLECTION" flag:"mongodb-authcode-collection" flagDesc:"The name of the mongodb auth code collection"`
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
