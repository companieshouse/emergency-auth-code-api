package service

import (
	"fmt"

	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/emergency-auth-code-api/config"
	"github.com/companieshouse/emergency-auth-code-api/models"
	"github.com/companieshouse/emergency-auth-code-api/oracle"
	"github.com/companieshouse/emergency-auth-code-api/transformers"
)

// GetOfficers returns the list of officers for the supplied company number
func GetOfficers(companyNumber string) (*models.OfficerListResponse, ResponseType, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, Error, nil
	}

	client := oracle.NewClient(cfg.OracleQueryAPIURL)
	oracleAPIResponse, err := client.GetOfficers(companyNumber)

	if err != nil {
		log.Error(fmt.Errorf("error getting officer list: [%v]", err))
		return nil, Error, err
	}

	if oracleAPIResponse == nil {
		return nil, NotFound, nil
	}

	resp := transformers.OfficerListResponse(oracleAPIResponse)

	return resp, Success, nil
}

// GetOfficer returns a single officer to be returned by the API for the supplied company number and officer id
func GetOfficer(companyNumber, officerID string) (*models.Officer, ResponseType, error) {
	oracleAPIResponse, responseType, err := GetOfficerDetails(companyNumber, officerID)
	if err != nil {
		return nil, Error, err
	}
	if responseType != Success {
		return nil, responseType, nil
	}

	resp := transformers.OfficerResponse(oracleAPIResponse)

	return resp, Success, nil

}

// GetOfficerDetails returns a single officer with values such as URA to be used internally only
func GetOfficerDetails(companyNumber, officerID string) (*oracle.Officer, ResponseType, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, Error, nil
	}

	client := oracle.NewClient(cfg.OracleQueryAPIURL)
	oracleAPIResponse, err := client.GetOfficer(companyNumber, officerID)

	if err != nil {
		log.Error(fmt.Errorf("error getting officer: [%v]", err))
		return nil, Error, err
	}

	if oracleAPIResponse == nil {
		return nil, NotFound, nil
	}

	return oracleAPIResponse, Success, nil
}
