package dao

import (
	"github.com/companieshouse/emergency-auth-code-api/config"
	"github.com/companieshouse/emergency-auth-code-api/models"
)

// AuthcodeDAOService interface declares how to interact with the persistence layer regardless of underlying technology
type AuthcodeDAOService interface {
	// CompanyHasAuthCode will check if the supplied company number has an auth code
	CompanyHasAuthCode(companyNumber string) (bool, error)
}

// AuthcodeRequestDAOService interface declares how to interact with the persistence layer regardless of underlying technology
type AuthcodeRequestDAOService interface {
	// InsertAuthcodeRequest creates an auth-code-request
	InsertAuthCodeRequest(dao *models.AuthCodeRequestResourceDao) error
	// GetAuthCodeRequest returns an auth-code-request
	GetAuthCodeRequest(authCodeRequestID string) (*models.AuthCodeRequestResourceDao, error)
	// UpdateAuthCodeRequestOfficer updates the officer details in an auth-code-request
	UpdateAuthCodeRequestOfficer(dao *models.AuthCodeRequestResourceDao) error
	// UpdateAuthCodeRequestStatus updates the status in an auth-code-request
	UpdateAuthCodeRequestStatus(dao *models.AuthCodeRequestResourceDao) error
	// CheckMultipleCorporateBodySubmissions checks whether multiple requests have been made for a company
	CheckMultipleCorporateBodySubmissions(companyNumber string) (bool, error)
}

// NewAuthCodeDAOService will create a new instance of the AuthCode Service interface.
// All details about its implementation and the
// database driver will be hidden from outside of this package
func NewAuthCodeDAOService(cfg *config.Config) AuthcodeDAOService {
	database := getMongoDatabase(cfg.MongoDBURL, cfg.MongoAuthcodeDatabase)
	return &MongoService{
		db:             database,
		CollectionName: cfg.MongoAuthCodeCollection,
	}
}

// NewAuthCodeRequestDAOService will create a new instance of the AuthCode Request Service interface.
// All details about its implementation and the
// database driver will be hidden from outside of this package
func NewAuthCodeRequestDAOService(cfg *config.Config) AuthcodeRequestDAOService {
	database := getMongoDatabase(cfg.MongoDBURL, cfg.MongoAuthcodeRequestDatabase)
	return &MongoService{
		db:             database,
		CollectionName: cfg.MongoAuthCodeRequestCollection,
	}
}
