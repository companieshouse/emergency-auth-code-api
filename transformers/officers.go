package transformers

import (
	"fmt"

	"github.com/companieshouse/emergency-auth-code-api/models"
	"github.com/companieshouse/emergency-auth-code-api/oracle"
)

// OfficerListResponse converts an Officer List from the Oracle API into the required format to be returned
func OfficerListResponse(oracleAPIResp *oracle.GetOfficersResponse) *models.OfficerListResponse {
	resp := models.OfficerListResponse{
		ItemsPerPage: oracleAPIResp.ItemsPerPage,
		StartIndex:   oracleAPIResp.StartIndex,
		TotalResults: oracleAPIResp.TotalResults,
	}
	for i := range oracleAPIResp.Items {
		officer := &oracleAPIResp.Items[i]

		officerItem := models.Officer{
			ID:          officer.ID,
			Name:        getOfficerName(officer.Forename, officer.Surname),
			OfficerRole: officer.OfficerRole,
			DateOfBirth: models.DateOfBirth{
				Month: officer.DateOfBirth.Month,
				Year:  officer.DateOfBirth.Year,
			},
			AppointedOn:        officer.AppointedOn,
			Nationality:        officer.Nationality,
			CountryOfResidence: officer.CountryOfResidence,
			Occupation:         officer.Occupation,
		}

		resp.Items = append(resp.Items, officerItem)
	}
	return &resp
}

// OfficerResponse converts an Officer from the Oracle API into the required format to be returned
func OfficerResponse(oracleAPIResp *oracle.Officer) *models.Officer {
	return &models.Officer{
		ID:          oracleAPIResp.ID,
		Name:        getOfficerName(oracleAPIResp.Forename, oracleAPIResp.Surname),
		OfficerRole: oracleAPIResp.OfficerRole,
		DateOfBirth: models.DateOfBirth{
			Month: oracleAPIResp.DateOfBirth.Month,
			Year:  oracleAPIResp.DateOfBirth.Year,
		},
		AppointedOn:        oracleAPIResp.AppointedOn,
		Nationality:        oracleAPIResp.Nationality,
		CountryOfResidence: oracleAPIResp.CountryOfResidence,
		Occupation:         oracleAPIResp.Occupation,
	}
}

func getOfficerName(forename, surname string) string {
	if forename != "" {
		return fmt.Sprintf("%s %s", forename, surname)
	}
	return surname
}
