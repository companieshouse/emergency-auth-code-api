package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"github.com/companieshouse/emergency-auth-code-api/models"

	"github.com/companieshouse/emergency-auth-code-api/service"

	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"

	. "github.com/smartystreets/goconvey/convey"
)

const companyNumber = "12345678"

func serveGetCompanyDirectorsHandler(directorService service.DirectorDatabase, companyNumber string) *httptest.ResponseRecorder {
	path := "/emergency-auth-code-service/company/" + companyNumber + "/officers"
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req = mux.SetURLVars(req, map[string]string{"company_number": companyNumber})
	res := httptest.NewRecorder()

	handler := GetCompanyDirectorsHandler(directorService)
	handler.ServeHTTP(res, req)

	return res
}

func directorServiceResponseValid() *models.CompanyOfficers {
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

func TestUnitGetCompanyDirectorsHandler(t *testing.T) {

	Convey("no company number provided in request", t, func() {
		httpmock.Activate()
		mockCtrl := gomock.NewController(t)
		defer httpmock.DeactivateAndReset()
		defer mockCtrl.Finish()

		mockService := service.NewMockDirectorDatabase(mockCtrl)

		res := serveGetCompanyDirectorsHandler(mockService, "")

		So(res.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("error receiving data from director database", t, func() {
		httpmock.Activate()
		mockCtrl := gomock.NewController(t)
		defer httpmock.DeactivateAndReset()
		defer mockCtrl.Finish()

		mockService := service.NewMockDirectorDatabase(mockCtrl)

		// expect the DirectorService to be called once for company directors and respond with valid data
		mockService.EXPECT().GetCompanyDirectors(companyNumber).Return(nil, errors.New("error receiving data from director database"))

		res := serveGetCompanyDirectorsHandler(mockService, companyNumber)

		So(res.Code, ShouldEqual, http.StatusInternalServerError)
	})

	Convey("no directors found for company", t, func() {
		httpmock.Activate()
		mockCtrl := gomock.NewController(t)
		defer httpmock.DeactivateAndReset()
		defer mockCtrl.Finish()

		mockService := service.NewMockDirectorDatabase(mockCtrl)

		// expect the DirectorService to be called once for company directors and respond with valid data
		mockService.EXPECT().GetCompanyDirectors(companyNumber).Return(&models.CompanyOfficers{Items: nil, TotalCount: 0}, nil)

		res := serveGetCompanyDirectorsHandler(mockService, companyNumber)

		So(res.Code, ShouldEqual, http.StatusNotFound)
	})

	Convey("successfully return company directors for company", t, func() {
		httpmock.Activate()
		mockCtrl := gomock.NewController(t)
		defer httpmock.DeactivateAndReset()
		defer mockCtrl.Finish()

		mockService := service.NewMockDirectorDatabase(mockCtrl)

		// expect the DirectorService to be called once for company directors and respond with valid data
		mockService.EXPECT().GetCompanyDirectors(companyNumber).Return(directorServiceResponseValid(), nil)

		res := serveGetCompanyDirectorsHandler(mockService, companyNumber)

		So(res.Code, ShouldEqual, http.StatusOK)
		So(res.Body, ShouldNotBeNil)
	})
}
