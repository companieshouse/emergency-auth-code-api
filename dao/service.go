package dao

import "github.com/companieshouse/emergency-auth-code-api/config"

// Service interface declares how to interact with the persistence layer regardless of underlying technology
type Service interface {
	// CompanyHasAuthCode will check if the supplied company number has an auth code
	CompanyHasAuthCode(companyNumber string) (bool, error)
}

// NewAuthCodeDAOService will create a new instance of the AuthCode Service interface.
// All details about its implementation and the
// database driver will be hidden from outside of this package
func NewAuthCodeDAOService(cfg *config.Config) Service {
	database := getMongoDatabase(cfg.MongoDBURL, cfg.MongoAccountDatabase)
	return &MongoService{
		db:             database,
		CollectionName: cfg.MongoAuthCodeCollection,
	}
}
