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

// GetOfficer returns the list of officers for the supplied company number
func GetOfficer(companyNumber, officerID string) (*models.Officer, ResponseType, error) {
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

	resp := transformers.OfficerResponse(oracleAPIResponse)

	return resp, Success, nil

}
