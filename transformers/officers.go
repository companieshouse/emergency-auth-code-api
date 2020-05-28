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

		var officerName string
		if officer.Forename != "" {
			officerName = fmt.Sprintf("%s %s", officer.Forename, officer.Surname)
		} else {
			officerName = officer.Surname
		}

		officerItem := models.OfficerListItem{
			ID:          officer.ID,
			Name:        officerName,
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
