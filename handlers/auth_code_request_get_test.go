package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/companieshouse/emergency-auth-code-api/dao"
	"github.com/companieshouse/emergency-auth-code-api/mocks"
	"github.com/companieshouse/emergency-auth-code-api/models"
	"github.com/companieshouse/emergency-auth-code-api/service"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"

	. "github.com/smartystreets/goconvey/convey"
)

var companyNumber = "12345678"
var companyName = "testCompany"

var daoResponse = models.AuthCodeRequestResourceDao{
	Data: models.AuthCodeRequestDataDao{
		CompanyNumber: companyNumber,
		CompanyName:   companyName,
	},
}

func serveGetAuthCodeRequest(daoReqSvc dao.AuthcodeRequestDAOService, hasAuthCodeRequestID bool) *httptest.ResponseRecorder {

	authCodeReqSvc := &service.AuthCodeRequestService{}

	if daoReqSvc != nil {
		authCodeReqSvc.DAO = daoReqSvc
	}

	h := GetAuthCodeRequest(authCodeReqSvc)
	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	if hasAuthCodeRequestID {
		req = mux.SetURLVars(req, map[string]string{"auth_code_request_id": companyNumber})
	}
	res := httptest.NewRecorder()

	h.ServeHTTP(res, req)

	return res
}

func TestUnitGetAuthCodeRequestHandler(t *testing.T) {
	Convey("GetAuthCodeRequest auth code request ID not provided", t, func() {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)

		res := serveGetAuthCodeRequest(mockDaoService, false)

		So(res.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("GetAuthCodeRequest returns error receiving auth code request", t, func() {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
		mockDaoService.EXPECT().GetAuthCodeRequest(companyNumber).Return(nil, fmt.Errorf("error"))

		res := serveGetAuthCodeRequest(mockDaoService, true)

		So(res.Code, ShouldEqual, http.StatusInternalServerError)
	})

	Convey("GetAuthCodeRequest returns no existing authcode request", t, func() {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
		mockDaoService.EXPECT().GetAuthCodeRequest(companyNumber).Return(nil, nil)

		res := serveGetAuthCodeRequest(mockDaoService, true)

		So(res.Code, ShouldEqual, http.StatusNotFound)
	})

	Convey("GetAuthCodeRequest successfully returns existing authcode request", t, func() {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
		mockDaoService.EXPECT().GetAuthCodeRequest(companyNumber).Return(&daoResponse, nil)

		res := serveGetAuthCodeRequest(mockDaoService, true)

		So(res.Code, ShouldEqual, http.StatusOK)
		responseBody := decodeResponse(res, t)
		So(responseBody.CompanyName, ShouldEqual, companyName)
		So(responseBody.CompanyNumber, ShouldEqual, companyNumber)
	})
}
