package service

import (
	"fmt"
	"testing"

	"github.com/companieshouse/emergency-auth-code-api/mocks"
	"github.com/companieshouse/emergency-auth-code-api/models"
	"github.com/golang/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"
)

const companyNumber = "12345678"

func TestUnitGetCompanyOfficers(t *testing.T) {

	officer1 := models.Items{
		ID:        "11111111",
		Forename1: "test1",
		Forename2: "test1",
		Surname:   "test1",
	}

	officer2 := models.Items{
		ID:        "22222222",
		Forename1: "test2",
		Forename2: "test2",
		Surname:   "test2",
	}

	officerServiceResponseValid := models.CompanyOfficers{
		Items:      []models.Items{officer1, officer2},
		TotalCount: 2,
	}

	Convey("Error returning officer details", t, func() {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockDaoService := mocks.NewMockOfficerDAOService(mockCtrl)
		mockDaoService.EXPECT().GetCompanyOfficers(companyNumber).Return(nil, fmt.Errorf("error returning officer details from database"))
		svc := OfficersService{DAO: mockDaoService}

		officers, err := svc.GetListOfCompanyOfficers(companyNumber)
		So(officers, ShouldBeNil)
		So(err, ShouldBeError, fmt.Errorf("error retrieving officer list from database: [error returning officer details from database]"))
	})

	Convey("Successfully return officer details", t, func() {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockDaoService := mocks.NewMockOfficerDAOService(mockCtrl)
		mockDaoService.EXPECT().GetCompanyOfficers(companyNumber).Return(&officerServiceResponseValid, nil)
		svc := OfficersService{DAO: mockDaoService}

		officers, err := svc.GetListOfCompanyOfficers(companyNumber)
		So(officers, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}
