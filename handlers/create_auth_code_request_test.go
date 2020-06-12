package handlers

import (
	"bytes"
	"context"
	"encoding/json"
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
	daoReqSvc dao.AuthcodeRequestDAOService) *httptest.ResponseRecorder {

	authCodeReqSvc := &service.AuthCodeRequestService{}

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
	Convey("CreateAuthCodeRequestHandler tests", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		Convey("authcode resource must be in context", func() {
			res := serveCreateAuthCodeRequestHandler(context.Background(), t, nil, nil)

			So(res.Code, ShouldEqual, http.StatusBadRequest)
			So(res.Body.String(), ShouldStartWith, `{"message":"failed to read request body"}`)
		})

		Convey("company number missing from request", func() {
			res := serveCreateAuthCodeRequestHandler(context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}), t, &models.AuthCodeRequest{}, nil)
			So(res.Code, ShouldEqual, http.StatusBadRequest)
			So(res.Body.String(), ShouldStartWith, `{"message":"company number missing from request"}`)
		})

		Convey("error calling oracle API for officer", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			defer httpmock.Reset()

			// stub the oracle query lookup
			responder := httpmock.NewStringResponder(http.StatusBadRequest, `{"total_results":3}`)
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/eligible-officers/12345678", responder)

			res := serveCreateAuthCodeRequestHandler(context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}), t, &models.AuthCodeRequest{CompanyNumber: "87654321", OfficerID: "12345678"}, nil)
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

			res := serveCreateAuthCodeRequestHandler(context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}), t, &models.AuthCodeRequest{CompanyNumber: "87654321", OfficerID: "12345678"}, nil)
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

			// stub the DB lookup
			mockReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)

			res := serveCreateAuthCodeRequestHandler(context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}), t, &models.AuthCodeRequest{CompanyNumber: "87654321", OfficerID: "12345678"}, mockReqService)
			So(res.Code, ShouldEqual, http.StatusInternalServerError)
			So(res.Body.String(), ShouldEqual, "")
		})

		Convey("successful Authcode Reminder", func() {
			defer httpmock.Reset()
			httpmock.RegisterResponder(http.MethodGet, "https://api.companieshouse.gov.uk/company/87654321", httpmock.NewStringResponder(http.StatusOK, companyDetailsResponse))

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// stub the oracle query lookup
			responder := httpmock.NewStringResponder(http.StatusOK, `{"total_results":3}`)
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/eligible-officers/12345678", responder)

			// stub the DB lookup
			mockReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockReqService.EXPECT().InsertAuthCodeRequest(gomock.Any()).Return(nil)

			res := serveCreateAuthCodeRequestHandler(context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}), t, &models.AuthCodeRequest{CompanyNumber: "87654321", OfficerID: "12345678"}, mockReqService)
			So(res.Code, ShouldEqual, http.StatusCreated)

			responseBody := decodeResponse(res, t)
			So(responseBody.CompanyName, ShouldEqual, "Test Company")
		})

		Convey("successful Authcode Request", func() {
			defer httpmock.Reset()
			httpmock.RegisterResponder(http.MethodGet, "https://api.companieshouse.gov.uk/company/87654321", httpmock.NewStringResponder(http.StatusOK, companyDetailsResponse))

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// stub the oracle query lookup
			responder := httpmock.NewStringResponder(http.StatusOK, `{"total_results":3}`)
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/eligible-officers/12345678", responder)

			// stub the DB lookup
			mockReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockReqService.EXPECT().InsertAuthCodeRequest(gomock.Any()).Return(nil)

			res := serveCreateAuthCodeRequestHandler(context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}), t, &models.AuthCodeRequest{CompanyNumber: "87654321", OfficerID: "12345678"}, mockReqService)
			So(res.Code, ShouldEqual, http.StatusCreated)

			responseBody := decodeResponse(res, t)
			So(responseBody.CompanyName, ShouldEqual, "Test Company")
		})
	})
}
