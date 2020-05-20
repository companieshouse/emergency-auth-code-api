package service

import (
	"net/http"

	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/go-sdk-manager/manager"
)

// GetCompanyName will attempt to get the company name from the CompanyProfileAPI.
func GetCompanyName(companyNumber string, req *http.Request) (string, error) {

	api, err := manager.GetSDK(req)
	if err != nil {
		log.ErrorR(req, err, log.Data{"company_number": companyNumber})
		return "", err
	}

	companyProfile, err := api.Profile.Get(companyNumber).Do()
	if err != nil {
		log.ErrorR(req, err, log.Data{"company_number": companyNumber})
		return "", err
	}

	return companyProfile.CompanyName, nil
}
