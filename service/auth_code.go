package service

import (
	"fmt"

	"github.com/companieshouse/emergency-auth-code-api/config"
	"github.com/companieshouse/emergency-auth-code-api/dao"
)

// AuthCodeService contains the DAO for db access
type AuthCodeService struct {
	DAO    dao.AuthcodeDAOService
	Config *config.Config
}

// CheckAuthCodeExists checks whether the specified company has an active auth code
func (s *AuthCodeService) CheckAuthCodeExists(companyNumber string) (bool, error) {
	companyHasAuthCode, err := s.DAO.CompanyHasAuthCode(companyNumber)
	if err != nil {
		err = fmt.Errorf("error checking DB for auth code: [%v]", err)
	}

	if !companyHasAuthCode {
		// backend processing expects an authcode item to exist, so need to create one here
		err := s.DAO.UpsertEmptyAuthCode(companyNumber)
		if err != nil {
			return false, err
		}
	}

	return companyHasAuthCode, err
}
