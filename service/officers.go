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
func GetOfficers(companyNumber string, startIndex string, itemsPerPage string) (*models.OfficerListResponse, ResponseType, error) {
	oracleAPIResponse, responseType, err := getOfficers(companyNumber, startIndex, itemsPerPage)
	if err != nil || responseType != Success {
		return nil, responseType, err
	}

	resp := transformers.OfficerListResponse(oracleAPIResponse)

	return resp, Success, nil
}

// CheckOfficers checks if a company has any eligible officers
func CheckOfficers(companyNumber string, startIndex string, itemsPerPage string) (bool, error) {
	_, responseType, err := getOfficers(companyNumber, startIndex, itemsPerPage)
	if err != nil {
		return false, err
	}

	if responseType == NotFound {
		return false, nil
	}
	return true, nil
}

func getOfficers(companyNumber string, startIndex string, itemsPerPage string) (*oracle.GetOfficersResponse, ResponseType, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, Error, nil
	}

	client := oracle.NewClient(cfg.OracleQueryAPIURL)
	oracleAPIResponse, err := client.GetOfficers(companyNumber, startIndex, itemsPerPage)

	if err != nil {
		log.Error(fmt.Errorf("error getting officer list: [%v]", err))
		return nil, Error, err
	}

	if oracleAPIResponse == nil {
		return nil, NotFound, nil
	}

	return oracleAPIResponse, Success, nil
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

// CheckCompanyFilingHistory returns a bool displaying whether the company has filed within the time period or not
func CheckCompanyFilingHistory(companyNumber string) (bool, error) {
	cfg, err := config.Get()
	if err != nil {
		return false, err
	}

	client := oracle.NewClient(cfg.OracleQueryAPIURL)
	filingHistoryCheck, err := client.CheckFilingHistory(companyNumber)

	if err != nil {
		log.Error(fmt.Errorf("error checking filing history: [%v]", err))
		return false, err
	}

	return filingHistoryCheck.EFilingFoundInPeriod, nil
}
