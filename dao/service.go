package dao

import (
	"github.com/companieshouse/emergency-auth-code-api/config"
	"github.com/companieshouse/emergency-auth-code-api/models"
)

// AuthcodeDAOService interface declares how to interact with the Authcode database layer regardless of underlying technology.
type AuthcodeDAOService interface {
	// CompanyHasAuthCode will check if the supplied company number has an auth code
	CompanyHasAuthCode(companyNumber string) (bool, error)
}

// OfficerDAOService interface declares how to interact with the Officer database layer regardless of underlying technology.
type OfficerDAOService interface {
	// GetCompanyOfficers returns all non-disqualified, natural, officers for a company
	GetCompanyOfficers(string) (*models.CompanyOfficers, error)
}

// NewAuthCodeDAOService will create a new instance of the AuthCodeDaoService interface using mongo.
// All details about its implementation and the
// database driver will be hidden from outside of this package
func NewAuthCodeDAOService(cfg *config.Config) AuthcodeDAOService {
	database := getMongoDatabase(cfg.MongoDBURL, cfg.MongoAccountDatabase)
	return &MongoService{
		db:             database,
		CollectionName: cfg.MongoAuthCodeCollection,
	}
}

// NewOfficerDAOService will create a new instance of the OfficerDaoService interface using oracle.
// All details about its implementation and the
// database driver will be hidden from outside this package
func NewOfficerDAOService(cfg *config.Config) OfficerDAOService {
	database := getOracleDatabase(*cfg)
	return &OracleService{
		db: database,
	}
}
