package service

import (
	"fmt"

	"github.com/companieshouse/emergency-auth-code-api/config"
	"github.com/companieshouse/emergency-auth-code-api/dao"
	"github.com/companieshouse/emergency-auth-code-api/models"
)

// AuthCodeRequestService contains the DAO for db access
type AuthCodeRequestService struct {
	DAO    dao.AuthcodeRequestDAOService
	Config *config.Config
}

// CreateAuthCodeRequest insert an auth code request into the database
func (s *AuthCodeRequestService) CreateAuthCodeRequest(requestDao *models.AuthCodeRequestResourceDao) error {

	err := s.DAO.InsertAuthCodeRequest(requestDao)
	if err != nil {
		err = fmt.Errorf("error creating AuthCode request: [%v]", err)
	}

	return err
}
