package service

import (
	"fmt"
	"net/http"

	"github.com/companieshouse/emergency-auth-code-api/config"
	"github.com/companieshouse/emergency-auth-code-api/dao"
	"github.com/companieshouse/emergency-auth-code-api/models"
	"github.com/companieshouse/emergency-auth-code-api/transformers"
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

// GetAuthCodeRequest returns an auth code request from the database
func (s *AuthCodeRequestService) GetAuthCodeRequest(authCodeRequestId string) (*models.AuthCodeRequestResourceResponse, int) {
	authCodeRequest, err := s.DAO.GetAuthCodeRequest(authCodeRequestId)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	if authCodeRequest == nil {
		return nil, http.StatusNotFound
	}

	return transformers.AuthCodeRequestResourceDaoToResponse(authCodeRequest), http.StatusOK
}
