package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/companieshouse/chs.go/authentication"
	"github.com/companieshouse/emergency-auth-code-api/dao"
	"github.com/companieshouse/emergency-auth-code-api/mocks"
	"github.com/companieshouse/emergency-auth-code-api/models"
	"github.com/companieshouse/emergency-auth-code-api/service"
	"github.com/companieshouse/go-session-handler/httpsession"
	"github.com/companieshouse/go-session-handler/session"
	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	. "github.com/smartystreets/goconvey/convey"
)

var companyDetailsResponse = `
{
  "company_name": "Test Company",
  "registered_office_address" : {
    "postal_code" : "CF14 3UZ",
    "address_line_2" : "Cardiff",
    "address_line_1" : "1 Crown Way"
  }
}
`

func serveCreateAuthCodeRequestHandler(
	ctx context.Context,
	t *testing.T,
	reqBody *models.AuthCodeRequest,
	daoSvc dao.AuthcodeDAOService,
	daoReqSvc dao.AuthcodeRequestDAOService) *httptest.ResponseRecorder {

	authCodeSvc := &service.AuthCodeService{}
	authCodeReqSvc := &service.AuthCodeRequestService{}

	if daoSvc != nil {
		authCodeSvc.DAO = daoSvc
		authCodeReqSvc.DAO = daoReqSvc
	}

	var body io.Reader
	if reqBody != nil {
		b, err := json.Marshal(reqBody)
		if err != nil {
			t.Fatal("failed to marshal request body")
		}
		body = bytes.NewReader(b)
	}

	ctx = context.WithValue(ctx, httpsession.ContextKeySession, &session.Session{})

	h := CreateAuthCodeRequest(authCodeSvc, authCodeReqSvc)
	req := httptest.NewRequest(http.MethodPost, "/", body).WithContext(ctx)
	res := httptest.NewRecorder()

	h.ServeHTTP(res, req.WithContext(ctx))

	return res
}

func decodeResponse(res *httptest.ResponseRecorder, t *testing.T) *models.AuthCodeRequestResourceResponse {
	if res.Body.Len() > 0 {
		var responseBody models.AuthCodeRequestResourceResponse
		err := json.NewDecoder(res.Body).Decode(&responseBody)
		if err != nil {
			t.Errorf("failed to read response body")
		}

		return &responseBody
	}
	return nil
}

func TestUnitCreateAuthCodeRequestHandler(t *testing.T) {
	Convey("CreateAuthCodeRequestHandler tests", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		Convey("authcode resource must be in context", func() {
			res := serveCreateAuthCodeRequestHandler(context.Background(), t, nil, nil, nil)

			So(res.Code, ShouldEqual, http.StatusBadRequest)
			So(res.Body.String(), ShouldStartWith, `{"message":"failed to read request body"}`)
		})

		Convey("company number missing from request", func() {
			res := serveCreateAuthCodeRequestHandler(context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}), t, &models.AuthCodeRequest{}, nil, nil)
			So(res.Code, ShouldEqual, http.StatusBadRequest)
			So(res.Body.String(), ShouldStartWith, `{"message":"company number missing from request"}`)
		})

		Convey("error checking DB for authcode", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// stub the DB lookup
			mockService := mocks.NewMockAuthcodeDAOService(mockCtrl)
			mockService.EXPECT().CompanyHasAuthCode(gomock.Any()).Return(false, fmt.Errorf("error"))
			mockReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)

			res := serveCreateAuthCodeRequestHandler(context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}), t, &models.AuthCodeRequest{CompanyNumber: "87654321", OfficerID: "12345678"}, mockService, mockReqService)
			So(res.Code, ShouldEqual, http.StatusInternalServerError)
			So(res.Body.String(), ShouldEqual, "")
		})

		Convey("successful Authcode Reminder", func() {
			defer httpmock.Reset()
			httpmock.RegisterResponder(http.MethodGet, "https://api.companieshouse.gov.uk/company/87654321", httpmock.NewStringResponder(http.StatusOK, companyDetailsResponse))

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// stub the DB lookup
			mockService := mocks.NewMockAuthcodeDAOService(mockCtrl)
			mockService.EXPECT().CompanyHasAuthCode(gomock.Any()).Return(true, nil)
			mockReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockReqService.EXPECT().InsertAuthCodeRequest(gomock.Any()).Return(nil)

			res := serveCreateAuthCodeRequestHandler(context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}), t, &models.AuthCodeRequest{CompanyNumber: "87654321", OfficerID: "12345678"}, mockService, mockReqService)
			So(res.Code, ShouldEqual, http.StatusCreated)

			responseBody := decodeResponse(res, t)
			So(responseBody.CompanyName, ShouldEqual, "Test Company")
		})

		Convey("successful Authcode Request", func() {
			defer httpmock.Reset()
			httpmock.RegisterResponder(http.MethodGet, "https://api.companieshouse.gov.uk/company/87654321", httpmock.NewStringResponder(http.StatusOK, companyDetailsResponse))

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// stub the DB lookup
			mockAuthcodeService := mocks.NewMockAuthcodeDAOService(mockCtrl)
			mockAuthcodeService.EXPECT().CompanyHasAuthCode(gomock.Any()).Return(false, nil)
			mockReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockReqService.EXPECT().InsertAuthCodeRequest(gomock.Any()).Return(nil)

			res := serveCreateAuthCodeRequestHandler(context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}), t, &models.AuthCodeRequest{CompanyNumber: "87654321", OfficerID: "12345678"}, mockAuthcodeService, mockReqService)
			So(res.Code, ShouldEqual, http.StatusCreated)

			responseBody := decodeResponse(res, t)
			So(responseBody.CompanyName, ShouldEqual, "Test Company")
		})
	})
}
