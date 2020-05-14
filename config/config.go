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
	BindAddr                 string `env:"BIND_ADDR" flag:"bind-addr" flagDesc:"Bind address"`
	DirectorDatabaseUsername string `env:"DIRECTOR_DATABASE_USERNAME" flag:"director-database-username" flagDesc:"Username to access directors database"`
	DirectorDatabasePassword string `env:"DIRECTOR_DATABASE_PASSWORD" flag:"director-database-password" flagDesc:"Password to access directors database"`
	DirectorDatabaseUrl      string `env:"DIRECTOR_DATABASE_URL"      flag:"director-database-url"      flagDesc:"URL to access directors database"`
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
