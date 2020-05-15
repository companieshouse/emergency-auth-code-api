package service

import (
	"fmt"

	"github.com/companieshouse/emergency-auth-code-api/config"
	"github.com/companieshouse/emergency-auth-code-api/dao"
	"github.com/companieshouse/emergency-auth-code-api/models"
)

// OfficersService contains the dao for db access
type OfficersService struct {
	Config *config.Config
	DAO    dao.OfficerDAOService
}

// GetListOfCompanyOfficers returns valid company officers from the officer database
func (s *OfficersService) GetListOfCompanyOfficers(companyNumber string) (*models.CompanyOfficers, error) {
	companyOfficers, err := s.DAO.GetCompanyOfficers(companyNumber)
	if err != nil {
		err = fmt.Errorf("error retrieving officer list from database: [%v]", err)
		return nil, err
	}

	return companyOfficers, err
}
