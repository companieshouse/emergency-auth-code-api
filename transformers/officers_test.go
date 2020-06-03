package transformers

import (
	"testing"

	"github.com/companieshouse/emergency-auth-code-api/models"
	"github.com/companieshouse/emergency-auth-code-api/oracle"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitOfficerTransformation(t *testing.T) {
	Convey("Officer List", t, func() {
		Convey("Response correctly converted", func() {
			input := oracle.GetOfficersResponse{
				ItemsPerPage: 1,
				StartIndex:   2,
				TotalResults: 3,
				Items: []oracle.Officer{
					{
						ID:          "123",
						Forename:    "Joe",
						Surname:     "Bloggs",
						OfficerRole: "director",
						DateOfBirth: oracle.DateOfBirth{
							Month: "3",
							Year:  "2012",
						},
						AppointedOn:        "01-01-2001",
						Nationality:        "British",
						CountryOfResidence: "Wales",
						Occupation:         "Director",
						UsualResidentialAddress: oracle.Address{
							ID:           "1",
							Premises:     "2",
							AddressLine1: "address-line-1",
							AddressLine2: "address-line-2",
							Locality:     "locality",
							Region:       "region",
							Postcode:     "CF14 3UZ",
						},
					},
				},
			}
			response := OfficerListResponse(&input)

			expected := models.OfficerListResponse{
				ItemsPerPage: 1,
				StartIndex:   2,
				TotalResults: 3,
				Items: []models.Officer{
					{
						ID: "12" +
							"3",
						Name:        "Joe Bloggs",
						OfficerRole: "director",
						DateOfBirth: models.DateOfBirth{
							Month: "3",
							Year:  "2012",
						},
						AppointedOn:        "01-01-2001",
						Nationality:        "British",
						CountryOfResidence: "Wales",
						Occupation:         "Director",
					},
				},
			}
			So(response.TotalResults, ShouldEqual, 3)
			So(response, ShouldResemble, &expected)
		})

		Convey("Officer with no forename converted", func() {
			input := oracle.GetOfficersResponse{Items: []oracle.Officer{{Surname: "Bloggs"}}}
			response := OfficerListResponse(&input)

			So(response.Items[0].Name, ShouldEqual, "Bloggs")
		})
	})

	Convey("Single Officer", t, func() {
		Convey("Response correctly converted", func() {
			input := oracle.Officer{
				ID:          "123",
				Forename:    "Joe",
				Surname:     "Bloggs",
				OfficerRole: "director",
				DateOfBirth: oracle.DateOfBirth{
					Month: "3",
					Year:  "2012",
				},
				AppointedOn:        "01-01-2001",
				Nationality:        "British",
				CountryOfResidence: "Wales",
				Occupation:         "Director",
				UsualResidentialAddress: oracle.Address{
					ID:           "1",
					Premises:     "2",
					AddressLine1: "address-line-1",
					AddressLine2: "address-line-2",
					Locality:     "locality",
					Region:       "region",
					Postcode:     "CF14 3UZ",
				},
			}
			response := OfficerResponse(&input)

			expected := models.Officer{
				ID: "12" +
					"3",
				Name:        "Joe Bloggs",
				OfficerRole: "director",
				DateOfBirth: models.DateOfBirth{
					Month: "3",
					Year:  "2012",
				},
				AppointedOn:        "01-01-2001",
				Nationality:        "British",
				CountryOfResidence: "Wales",
				Occupation:         "Director",
			}
			So(response, ShouldResemble, &expected)
		})

	})

}
