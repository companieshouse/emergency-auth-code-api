package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/companieshouse/emergency-auth-code-api/dao"
	"github.com/companieshouse/emergency-auth-code-api/mocks"

	"github.com/gorilla/mux"

	"github.com/companieshouse/emergency-auth-code-api/models"

	"github.com/companieshouse/emergency-auth-code-api/service"

	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"

	. "github.com/smartystreets/goconvey/convey"
)

const companyNumber = "12345678"

func serveGetCompanyOfficersHandler(officerDaoService dao.OfficerDAOService, companyNumber string) *httptest.ResponseRecorder {
	path := "/emergency-auth-code-service/company/" + companyNumber + "/officers"
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req = mux.SetURLVars(req, map[string]string{"company_number": companyNumber})
	res := httptest.NewRecorder()
	svc := service.OfficersService{
		Config: nil,
		DAO:    officerDaoService,
	}

	handler := GetCompanyOfficersHandler(&svc)
	handler.ServeHTTP(res, req)

	return res
}

func officerServiceResponseValid() *models.CompanyOfficers {
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

	return &models.CompanyOfficers{
		Items:      []models.Items{officer1, officer2},
		TotalCount: 2,
	}
}

func TestUnitGetCompanyOfficersHandler(t *testing.T) {

	Convey("no company number provided in request", t, func() {
		httpmock.Activate()
		mockCtrl := gomock.NewController(t)
		defer httpmock.DeactivateAndReset()
		defer mockCtrl.Finish()

		mockService := mocks.NewMockOfficerDAOService(mockCtrl)

		res := serveGetCompanyOfficersHandler(mockService, "")

		So(res.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("error receiving data from officer database", t, func() {
		httpmock.Activate()
		mockCtrl := gomock.NewController(t)
		defer httpmock.DeactivateAndReset()
		defer mockCtrl.Finish()

		mockDaoService := mocks.NewMockOfficerDAOService(mockCtrl)

		// expect the OfficerService to be called for company officers and respond with valid data
		mockDaoService.EXPECT().GetCompanyOfficers(companyNumber).Return(nil, errors.New("error receiving data from officer database"))

		res := serveGetCompanyOfficersHandler(mockDaoService, companyNumber)

		So(res.Code, ShouldEqual, http.StatusInternalServerError)
	})

	Convey("no officers found for company", t, func() {
		httpmock.Activate()
		mockCtrl := gomock.NewController(t)
		defer httpmock.DeactivateAndReset()
		defer mockCtrl.Finish()

		mockDaoService := mocks.NewMockOfficerDAOService(mockCtrl)

		// expect the OfficerService to be called once for company officers and respond with valid data
		mockDaoService.EXPECT().GetCompanyOfficers(companyNumber).Return(&models.CompanyOfficers{Items: nil, TotalCount: 0}, nil)

		res := serveGetCompanyOfficersHandler(mockDaoService, companyNumber)

		So(res.Code, ShouldEqual, http.StatusNotFound)
	})

	Convey("successfully return company officers for company", t, func() {
		httpmock.Activate()
		mockCtrl := gomock.NewController(t)
		defer httpmock.DeactivateAndReset()
		defer mockCtrl.Finish()

		mockService := mocks.NewMockOfficerDAOService(mockCtrl)

		// expect the OfficerService to be called once for company officers and respond with valid data
		mockService.EXPECT().GetCompanyOfficers(companyNumber).Return(officerServiceResponseValid(), nil)

		res := serveGetCompanyOfficersHandler(mockService, companyNumber)

		So(res.Code, ShouldEqual, http.StatusOK)
		So(res.Body, ShouldNotBeNil)
	})
}
