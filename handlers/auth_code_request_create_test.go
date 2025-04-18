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
	"github.com/companieshouse/emergency-auth-code-api/config"
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

var (
	companyDetailsResponse = `
		{
		"company_name": "Test Company",
		"registered_office_address" : {
			"postal_code" : "CF14 3UZ",
			"address_line_2" : "Cardiff",
			"address_line_1" : "1 Crown Way"
		}
		}
	`
	testBasePath = "http://test-path.gov"
	testResource = testBasePath + "/company/87654321"
)

func serveCreateAuthCodeRequestHandler(
	ctx context.Context,
	t *testing.T,
	reqBody *models.AuthCodeRequest,
	daoReqSvc dao.AuthcodeRequestDAOService) *httptest.ResponseRecorder {

	authCodeReqSvc := &service.AuthCodeRequestService{
		Config: &config.Config{
			APIBaseURL: testBasePath,
		},
	}

	if daoReqSvc != nil {
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

	h := CreateAuthCodeRequest(authCodeReqSvc)
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
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	Convey("CreateAuthCodeRequestHandler tests", t, func() {

		Convey("authcode resource must be in context", func() {
			res := serveCreateAuthCodeRequestHandler(context.Background(), t, nil, nil)

			So(res.Code, ShouldEqual, http.StatusBadRequest)
			So(res.Body.String(), ShouldStartWith, `{"message":"failed to read request body"}`)
		})

		Convey("company number missing from request", func() {
			res := serveCreateAuthCodeRequestHandler(context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}), t, &models.AuthCodeRequest{}, nil)
			So(res.Body.String(), ShouldStartWith, `{"message":"company number missing from request"}`)
		})

		Convey("error calling oracle API for company filing history", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			defer httpmock.Reset()

			// stub the oracle query lookup for the officer
			responderOfficer := httpmock.NewStringResponder(http.StatusOK, `{"total_results":3}`)
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/eligible-officers/12345678", responderOfficer)

			// stub the oracle query lookup for the filing history
			responderFilingHistory := httpmock.NewStringResponder(http.StatusBadRequest, `{"efiling_found_in_period":true}`)
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/efiling-status", responderFilingHistory)

			// stub the DB lookup
			mockReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockReqService.EXPECT().CheckMultipleCorporateBodySubmissions(gomock.Any()).Return(false, nil)
			mockReqService.EXPECT().CheckMultipleUserSubmissions(gomock.Any()).Return(false, nil)

			res := serveCreateAuthCodeRequestHandler(
				context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}),
				t,
				&models.AuthCodeRequest{CompanyNumber: "87654321", OfficerID: "12345678"},
				mockReqService,
			)

			b := res.Body.String()

			So(res.Code, ShouldEqual, http.StatusInternalServerError)
			So(b, ShouldStartWith, `{"message":"error checking corporate body"}`)
		})

		Convey("company has had a filing within recent filing period", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			defer httpmock.Reset()

			// stub the oracle query lookup for the officer
			responderOfficer := httpmock.NewStringResponder(http.StatusOK, `{"total_results":3}`)
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/eligible-officers/12345678", responderOfficer)

			// stub the oracle query lookup for the filing history
			responderFilingHistory := httpmock.NewStringResponder(http.StatusOK, `{"efiling_found_in_period":true}`)
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/efiling-status", responderFilingHistory)

			// stub the DB lookup
			mockReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockReqService.EXPECT().CheckMultipleCorporateBodySubmissions(gomock.Any()).Return(false, nil)
			mockReqService.EXPECT().CheckMultipleUserSubmissions(gomock.Any()).Return(false, nil)

			res := serveCreateAuthCodeRequestHandler(
				context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}),
				t,
				&models.AuthCodeRequest{CompanyNumber: "87654321", OfficerID: "12345678"},
				mockReqService,
			)
			So(res.Code, ShouldEqual, http.StatusForbidden)
			So(res.Body.String(), ShouldStartWith, `{"message":"request not permitted for corporate body"}`)
		})

		Convey("error calling oracle API for officer", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			defer httpmock.Reset()

			// stub the oracle query lookup
			responder := httpmock.NewStringResponder(http.StatusBadRequest, `{"total_results":3}`)
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/eligible-officers/12345678", responder)

			// stub the oracle query lookup for the filing history
			responderFilingHistory := httpmock.NewStringResponder(http.StatusOK, `{"efiling_found_in_period":false}`)
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/efiling-status", responderFilingHistory)

			// stub the DB lookup
			mockReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockReqService.EXPECT().CheckMultipleCorporateBodySubmissions(gomock.Any()).Return(false, nil)
			mockReqService.EXPECT().CheckMultipleUserSubmissions(gomock.Any()).Return(false, nil)

			res := serveCreateAuthCodeRequestHandler(
				context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}),
				t,
				&models.AuthCodeRequest{CompanyNumber: "87654321", OfficerID: "12345678"},
				mockReqService,
			)
			So(res.Code, ShouldEqual, http.StatusInternalServerError)
			So(res.Body.String(), ShouldStartWith, `{"message":"there was a problem communicating with the Oracle API"}`)
		})

		Convey("no officer with that ID found for company", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			defer httpmock.Reset()

			// stub the oracle query lookup
			responder := httpmock.NewStringResponder(http.StatusNotFound, `{"total_results":3}`)
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/eligible-officers/12345678", responder)

			// stub the oracle query lookup for the filing history
			responderFilingHistory := httpmock.NewStringResponder(http.StatusOK, `{"efiling_found_in_period":false}`)
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/efiling-status", responderFilingHistory)

			// stub the DB lookup
			mockReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockReqService.EXPECT().CheckMultipleCorporateBodySubmissions(gomock.Any()).Return(false, nil)
			mockReqService.EXPECT().CheckMultipleUserSubmissions(gomock.Any()).Return(false, nil)

			res := serveCreateAuthCodeRequestHandler(
				context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}),
				t,
				&models.AuthCodeRequest{CompanyNumber: "87654321", OfficerID: "12345678"},
				mockReqService,
			)
			So(res.Code, ShouldEqual, http.StatusNotFound)
			So(res.Body.String(), ShouldStartWith, `{"message":"No officer found"}`)
		})

		Convey("error getting company name", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			defer httpmock.Reset()

			// stub the oracle query lookup
			responder := httpmock.NewStringResponder(http.StatusOK, `{"total_results":3}`)
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/eligible-officers/12345678", responder)

			// stub the oracle query lookup for the filing history
			responderFilingHistory := httpmock.NewStringResponder(http.StatusOK, `{"efiling_found_in_period":false}`)
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/efiling-status", responderFilingHistory)

			// stub the DB lookup
			mockReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockReqService.EXPECT().CheckMultipleCorporateBodySubmissions(gomock.Any()).Return(false, nil)
			mockReqService.EXPECT().CheckMultipleUserSubmissions(gomock.Any()).Return(false, nil)

			res := serveCreateAuthCodeRequestHandler(
				context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}),
				t,
				&models.AuthCodeRequest{CompanyNumber: "87654321", OfficerID: "12345678"},
				mockReqService,
			)
			So(res.Code, ShouldEqual, http.StatusInternalServerError)
			So(res.Body.String(), ShouldEqual, "")
		})

		Convey("no eligible officers", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			defer httpmock.Reset()

			// stub the DB lookup
			mockReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockReqService.EXPECT().CheckMultipleCorporateBodySubmissions(gomock.Any()).Return(false, nil)
			mockReqService.EXPECT().CheckMultipleUserSubmissions(gomock.Any()).Return(false, nil)

			// stub the oracle query lookup
			responder := httpmock.NewStringResponder(http.StatusNotFound, "")
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/eligible-officers", responder)

			// stub the oracle query lookup for the filing history
			responderFilingHistory := httpmock.NewStringResponder(http.StatusOK, `{"efiling_found_in_period":false}`)
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/efiling-status", responderFilingHistory)

			res := serveCreateAuthCodeRequestHandler(
				context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}),
				t,
				&models.AuthCodeRequest{CompanyNumber: "87654321", OfficerID: ""},
				mockReqService)
			So(res.Code, ShouldEqual, http.StatusNotFound)

		})

		Convey("no eligible officers - error", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			defer httpmock.Reset()

			// stub the DB lookup
			mockReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockReqService.EXPECT().CheckMultipleCorporateBodySubmissions(gomock.Any()).Return(false, nil)
			mockReqService.EXPECT().CheckMultipleUserSubmissions(gomock.Any()).Return(false, nil)

			// stub the oracle query lookup
			responder := httpmock.NewStringResponder(http.StatusInternalServerError, "")
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/eligible-officers", responder)

			// stub the oracle query lookup for the filing history
			responderFilingHistory := httpmock.NewStringResponder(http.StatusOK, `{"efiling_found_in_period":false}`)
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/efiling-status", responderFilingHistory)

			res := serveCreateAuthCodeRequestHandler(
				context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}),
				t,
				&models.AuthCodeRequest{CompanyNumber: "87654321", OfficerID: ""},
				mockReqService,
			)
			So(res.Code, ShouldEqual, http.StatusInternalServerError)

		})

		Convey("successful Authcode Reminder", func() {
			defer httpmock.Reset()
			httpmock.RegisterResponder(http.MethodGet, testResource, httpmock.NewStringResponder(http.StatusOK, companyDetailsResponse))

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// stub the oracle query lookup
			responder := httpmock.NewStringResponder(http.StatusOK, `{"total_results":3}`)
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/eligible-officers/12345678", responder)

			// stub the DB lookup
			mockReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockReqService.EXPECT().InsertAuthCodeRequest(gomock.Any()).Return(nil)
			mockReqService.EXPECT().CheckMultipleCorporateBodySubmissions(gomock.Any()).Return(false, nil)
			mockReqService.EXPECT().CheckMultipleUserSubmissions(gomock.Any()).Return(false, nil)

			// stub the oracle query lookup for the filing history
			responderFilingHistory := httpmock.NewStringResponder(http.StatusOK, `{"efiling_found_in_period":false}`)
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/efiling-status", responderFilingHistory)

			res := serveCreateAuthCodeRequestHandler(
				context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}),
				t,
				&models.AuthCodeRequest{CompanyNumber: "87654321", OfficerID: "12345678"},
				mockReqService,
			)
			So(res.Code, ShouldEqual, http.StatusCreated)

			responseBody := decodeResponse(res, t)
			So(responseBody.CompanyName, ShouldEqual, "Test Company")
		})

		Convey("successful Authcode Request", func() {
			defer httpmock.Reset()
			httpmock.RegisterResponder(http.MethodGet, testResource, httpmock.NewStringResponder(http.StatusOK, companyDetailsResponse))

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// stub the oracle query lookup
			responder := httpmock.NewStringResponder(http.StatusOK, `{"total_results":3}`)
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/eligible-officers/12345678", responder)

			// stub the oracle query lookup for the filing history
			responderFilingHistory := httpmock.NewStringResponder(http.StatusOK, `{"efiling_found_in_period":false}`)
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/efiling-status", responderFilingHistory)

			// stub the DB lookup
			mockReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockReqService.EXPECT().InsertAuthCodeRequest(gomock.Any()).Return(nil)
			mockReqService.EXPECT().CheckMultipleCorporateBodySubmissions(gomock.Any()).Return(false, nil)
			mockReqService.EXPECT().CheckMultipleUserSubmissions(gomock.Any()).Return(false, nil)

			res := serveCreateAuthCodeRequestHandler(
				context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}),
				t,
				&models.AuthCodeRequest{CompanyNumber: "87654321", OfficerID: "12345678"},
				mockReqService,
			)
			So(res.Code, ShouldEqual, http.StatusCreated)

			responseBody := decodeResponse(res, t)
			So(responseBody.CompanyName, ShouldEqual, "Test Company")
		})
	})

	Convey("Multiple submissions for company", t, func() {

		Convey("multiple submissions", func() {

			defer httpmock.Reset()
			httpmock.RegisterResponder(http.MethodGet, testResource, httpmock.NewStringResponder(http.StatusOK, companyDetailsResponse))

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// stub the DB lookup
			mockReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockReqService.EXPECT().CheckMultipleCorporateBodySubmissions(gomock.Any()).Return(true, nil)

			res := serveCreateAuthCodeRequestHandler(
				context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}),
				t,
				&models.AuthCodeRequest{CompanyNumber: "87654321", OfficerID: "12345678"},
				mockReqService,
			)
			So(res.Code, ShouldEqual, http.StatusForbidden)

		})

		Convey("error checking submissions", func() {

			defer httpmock.Reset()
			httpmock.RegisterResponder(http.MethodGet, testResource, httpmock.NewStringResponder(http.StatusOK, companyDetailsResponse))

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// stub the DB lookup
			mockReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockReqService.EXPECT().CheckMultipleCorporateBodySubmissions(gomock.Any()).Return(false, fmt.Errorf("error"))

			res := serveCreateAuthCodeRequestHandler(
				context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}),
				t,
				&models.AuthCodeRequest{CompanyNumber: "87654321", OfficerID: "12345678"},
				mockReqService,
			)
			So(res.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})

	Convey("Multiple submissions for user", t, func() {

		Convey("multiple submissions", func() {

			defer httpmock.Reset()
			httpmock.RegisterResponder(http.MethodGet, testResource, httpmock.NewStringResponder(http.StatusOK, companyDetailsResponse))

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// stub the DB lookup
			mockReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockReqService.EXPECT().CheckMultipleCorporateBodySubmissions(gomock.Any()).Return(false, nil)
			mockReqService.EXPECT().CheckMultipleUserSubmissions(gomock.Any()).Return(true, nil)

			res := serveCreateAuthCodeRequestHandler(
				context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}),
				t,
				&models.AuthCodeRequest{CompanyNumber: "87654321", OfficerID: "12345678"},
				mockReqService,
			)
			So(res.Code, ShouldEqual, http.StatusForbidden)

		})

		Convey("error checking submissions", func() {

			defer httpmock.Reset()
			httpmock.RegisterResponder(http.MethodGet, testResource, httpmock.NewStringResponder(http.StatusOK, companyDetailsResponse))

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// stub the DB lookup
			mockReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockReqService.EXPECT().CheckMultipleCorporateBodySubmissions(gomock.Any()).Return(false, nil)
			mockReqService.EXPECT().CheckMultipleUserSubmissions(gomock.Any()).Return(false, fmt.Errorf("error"))

			res := serveCreateAuthCodeRequestHandler(
				context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}),
				t,
				&models.AuthCodeRequest{CompanyNumber: "87654321", OfficerID: "12345678"},
				mockReqService,
			)
			So(res.Code, ShouldEqual, http.StatusInternalServerError)

		})
	})
}
