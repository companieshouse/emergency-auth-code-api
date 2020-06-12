package service

import (
	"fmt"
	"net/http"

	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/emergency-auth-code-api/config"
	"github.com/companieshouse/emergency-auth-code-api/dao"
	"github.com/companieshouse/emergency-auth-code-api/models"
	"github.com/companieshouse/emergency-auth-code-api/oracle"
	"github.com/companieshouse/emergency-auth-code-api/queueapi"
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

// UpdateAuthCodeRequestOfficer updates the officer details in an authcode request
func (s *AuthCodeRequestService) UpdateAuthCodeRequestOfficer(
	authCodeReqDao *models.AuthCodeRequestResourceDao, authCodeRequestID string, officer *oracle.Officer) (
	*models.AuthCodeRequestResourceResponse, ResponseType) {

	requestDao := models.AuthCodeRequestResourceDao{
		ID: authCodeRequestID,
		Data: models.AuthCodeRequestDataDao{
			OfficerID:       officer.ID,
			OfficerUraID:    officer.UsualResidentialAddress.ID,
			OfficerForename: officer.Forename,
			OfficerSurname:  officer.Surname,
		},
	}

	err := s.DAO.UpdateAuthCodeRequestOfficer(&requestDao)
	if err != nil {
		return nil, Error
	}

	authCodeReqDao.Data.OfficerID = officer.ID
	authCodeReqDao.Data.OfficerUraID = officer.UsualResidentialAddress.ID
	authCodeReqDao.Data.OfficerForename = officer.Forename
	authCodeReqDao.Data.OfficerSurname = officer.Surname

	return transformers.AuthCodeRequestResourceDaoToResponse(authCodeReqDao), Success
}

// UpdateAuthCodeRequestStatus updates the status in an authcode request
func (s *AuthCodeRequestService) UpdateAuthCodeRequestStatus(authCodeReqDao *models.AuthCodeRequestResourceDao, authCodeRequestID string, status string, companyHasAuthCode bool) (*models.AuthCodeRequestResourceResponse, ResponseType) {

	requestDao := models.AuthCodeRequestResourceDao{
		ID: authCodeRequestID,
		Data: models.AuthCodeRequestDataDao{
			Status: status,
			Type:   getLetterType(companyHasAuthCode),
		},
	}

	err := s.DAO.UpdateAuthCodeRequestStatus(&requestDao)
	if err != nil {
		return nil, Error
	}

	authCodeReqDao.Data.Status = status

	return transformers.AuthCodeRequestResourceDaoToResponse(authCodeReqDao), Success
}

// SendAuthCodeRequest sends a letter item to the Queue API
func (s *AuthCodeRequestService) SendAuthCodeRequest(authCodeReqDao *models.AuthCodeRequestResourceDao, companyNumber string, userEmail string, companyHasAuthCode bool) ResponseType {

	// get Officer residential address
	companyOfficer, responseType, err := GetOfficerDetails(companyNumber, authCodeReqDao.Data.OfficerID)

	if err != nil || responseType == Error {
		log.Error(fmt.Errorf("error calling Oracle API to get officer: %v", err))
		return Error
	}

	if responseType == NotFound {
		log.Error(fmt.Errorf("officer not found"))
		return NotFound
	}

	var officerName string
	if companyOfficer.Forename != "" {
		officerName = fmt.Sprintf("%s %s", companyOfficer.Forename, companyOfficer.Surname)
	} else {
		officerName = companyOfficer.Surname
	}

	queueItem := models.QueueItem{
		Type:          "authcode_put",
		Email:         userEmail,
		CompanyNumber: companyNumber,
		CompanyName:   officerName,
		Address: models.Address{
			POBox:        companyOfficer.UsualResidentialAddress.PoBox,
			Premises:     companyOfficer.UsualResidentialAddress.Premises,
			AddressLine1: companyOfficer.UsualResidentialAddress.AddressLine1,
			AddressLine2: companyOfficer.UsualResidentialAddress.AddressLine2,
			Locality:     companyOfficer.UsualResidentialAddress.Locality,
			Region:       companyOfficer.UsualResidentialAddress.Region,
			PostalCode:   companyOfficer.UsualResidentialAddress.Postcode,
			Country:      companyOfficer.UsualResidentialAddress.Country,
		},
		Status: getLetterType(companyHasAuthCode),
	}

	err = sendQueueAPI(&queueItem)
	if err != nil {
		log.Error(err)
		return Error
	}

	return Success
}

func sendQueueAPI(item *models.QueueItem) error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}

	client := queueapi.NewClient(cfg.QueueAPILocalURL)
	err = client.SendQueueItem(item)
	return err
}

// GetAuthCodeReqDao returns an authcode request db object
func (s *AuthCodeRequestService) GetAuthCodeReqDao(authCodeRequestID, companyNumber string) (*models.AuthCodeRequestResourceDao, ResponseType) {
	authCodeRequest, err := s.DAO.GetAuthCodeRequest(authCodeRequestID)
	if err != nil {
		return nil, Error
	}
	if authCodeRequest == nil {
		return nil, NotFound
	}

	if authCodeRequest.Data.CompanyNumber != companyNumber {
		return nil, InvalidData
	}

	return authCodeRequest, Success
}

func getLetterType(companyHasAuthCode bool) string {
	if companyHasAuthCode {
		return "reminder"
	}
	return "apply"
}
