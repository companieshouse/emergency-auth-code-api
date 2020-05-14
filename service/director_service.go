package service

import (
	"database/sql"
	"errors"

	"github.com/companieshouse/emergency-auth-code-api/config"
	"github.com/companieshouse/emergency-auth-code-api/models"
)

// DirectorDatabase interface declares how to interact with the service layer that has access to the directors database
type DirectorDatabase interface {
	// GetCompanyDirectors returns all non-disqualified, natural, directors for a company
	GetCompanyDirectors(string) (*models.CompanyOfficers, error)
}

// NewOracleService will create a new instance of the DirectorDatabase interface, with a connection to an oracle database driver
func NewOracleService(cfg *config.Config) (*OracleService, error) {
	db, err := sql.Open("godror", generateOracleUrl(*cfg))
	if err != nil {
		return nil, errors.New("failed to connect to oracle database")
	}

	return &OracleService{
		DatabaseDriver: db,
	}, nil
}

// generateOracleUrl generates an Oracle database URL from the config
func generateOracleUrl(cfg config.Config) string {
	return cfg.DirectorDatabaseUsername + "/" + cfg.DirectorDatabasePassword + "@//" + cfg.DirectorDatabaseUrl
}
