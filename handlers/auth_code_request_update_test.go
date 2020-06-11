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
	"github.com/gorilla/mux"
	"github.com/jarcoal/httpmock"
	. "github.com/smartystreets/goconvey/convey"
)

func serveUpdateAuthCodeRequestHandler(
	ctx context.Context,
	t *testing.T,
	reqBody *models.AuthCodeRequest,
	authCodeReqID string,
	daoSvc dao.AuthcodeDAOService,
	daoReqSvc dao.AuthcodeRequestDAOService) *httptest.ResponseRecorder {

	authCodeSvc := &service.AuthCodeService{}
	authCodeReqSvc := &service.AuthCodeRequestService{}

	if daoSvc != nil {
		authCodeSvc.DAO = daoSvc
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

	h := UpdateAuthCodeRequest(authCodeSvc, authCodeReqSvc)
	req := httptest.NewRequest(http.MethodPost, "/", body).WithContext(ctx)

	if authCodeReqID != "" {
		req = mux.SetURLVars(req, map[string]string{"auth_code_request_id": authCodeReqID})
	}

	res := httptest.NewRecorder()

	h.ServeHTTP(res, req)

	return res
}

func TestUnitUpdateAuthCodeRequestHandler(t *testing.T) {
	Convey("UpdateAuthCodeRequestHandler tests", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		Convey("authcode resource must be in context", func() {
			res := serveUpdateAuthCodeRequestHandler(context.Background(), t, nil, "", nil, nil)

			So(res.Code, ShouldEqual, http.StatusBadRequest)
			So(res.Body.String(), ShouldStartWith, `{"message":"failed to read request body"}`)
		})

		Convey("authcode request ID missing from request", func() {
			res := serveUpdateAuthCodeRequestHandler(context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}), t, &models.AuthCodeRequest{}, "", nil, nil)
			So(res.Code, ShouldEqual, http.StatusBadRequest)
			So(res.Body.String(), ShouldStartWith, `{"message":"auth code request ID missing from request"}`)
		})

		Convey("company number missing from request", func() {
			res := serveUpdateAuthCodeRequestHandler(context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}), t, &models.AuthCodeRequest{}, "123", nil, nil)
			So(res.Code, ShouldEqual, http.StatusBadRequest)
			So(res.Body.String(), ShouldStartWith, `{"message":"company number missing from request"}`)
		})

		Convey("no valid changes", func() {
			res := serveUpdateAuthCodeRequestHandler(context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}), t, &models.AuthCodeRequest{CompanyNumber: "87654321"}, "123", nil, nil)
			So(res.Code, ShouldEqual, http.StatusBadRequest)
			So(res.Body.String(), ShouldStartWith, `{"message":"no valid changes supplied"}`)
		})

		Convey("error reading authcode request", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			authCodeDaoResponse := models.AuthCodeRequestResourceDao{
				Data: models.AuthCodeRequestDataDao{
					CompanyNumber: "87654321",
					Status:        "submitted",
				},
			}

			mockDaoReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockDaoReqService.EXPECT().GetAuthCodeRequest("123").Return(&authCodeDaoResponse, fmt.Errorf("error"))

			res := serveUpdateAuthCodeRequestHandler(context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}), t, &models.AuthCodeRequest{CompanyNumber: "87654321", OfficerID: "98765432"}, "123", nil, mockDaoReqService)
			So(res.Code, ShouldEqual, http.StatusInternalServerError)
			So(res.Body.String(), ShouldStartWith, `{"message":"error reading auth code request"}`)
		})

		Convey("request already submitted", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			authCodeDaoResponse := models.AuthCodeRequestResourceDao{
				Data: models.AuthCodeRequestDataDao{
					CompanyNumber: "87654321",
					Status:        "submitted",
				},
			}

			mockDaoReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockDaoReqService.EXPECT().GetAuthCodeRequest("123").Return(&authCodeDaoResponse, nil)

			res := serveUpdateAuthCodeRequestHandler(context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}), t, &models.AuthCodeRequest{CompanyNumber: "87654321", OfficerID: "98765432"}, "123", nil, mockDaoReqService)
			So(res.Code, ShouldEqual, http.StatusBadRequest)
			So(res.Body.String(), ShouldStartWith, `{"message":"request already submitted"}`)
		})

		Convey("officer update", func() {

			Convey("error calling Oracle API", func() {
				mockCtrl := gomock.NewController(t)
				defer mockCtrl.Finish()

				authCodeDaoResponse := models.AuthCodeRequestResourceDao{
					Data: models.AuthCodeRequestDataDao{
						CompanyNumber: "87654321",
					},
				}

				mockDaoReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
				mockDaoReqService.EXPECT().GetAuthCodeRequest("123").Return(&authCodeDaoResponse, nil)

				res := serveUpdateAuthCodeRequestHandler(context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}), t, &models.AuthCodeRequest{CompanyNumber: "87654321", OfficerID: "98765432"}, "123", nil, mockDaoReqService)
				So(res.Code, ShouldEqual, http.StatusInternalServerError)
				So(res.Body.String(), ShouldStartWith, `{"message":"there was a problem communicating with the Oracle API"}`)
			})

			Convey("officer not found", func() {
				mockCtrl := gomock.NewController(t)
				defer mockCtrl.Finish()

				authCodeDaoResponse := models.AuthCodeRequestResourceDao{
					Data: models.AuthCodeRequestDataDao{
						CompanyNumber: "87654321",
					},
				}

				mockDaoReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
				mockDaoReqService.EXPECT().GetAuthCodeRequest("123").Return(&authCodeDaoResponse, nil)

				httpmock.Activate()
				defer httpmock.DeactivateAndReset()
				responder := httpmock.NewStringResponder(http.StatusNotFound, "")
				httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/eligible-officers/98765432", responder)

				res := serveUpdateAuthCodeRequestHandler(context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}), t, &models.AuthCodeRequest{CompanyNumber: "87654321", OfficerID: "98765432"}, "123", nil, mockDaoReqService)
				So(res.Code, ShouldEqual, http.StatusNotFound)
				So(res.Body.String(), ShouldStartWith, `{"message":"No officer found"}`)
			})

			Convey("error updating officer details - error", func() {
				mockCtrl := gomock.NewController(t)
				defer mockCtrl.Finish()

				authCodeDaoResponse := models.AuthCodeRequestResourceDao{
					Data: models.AuthCodeRequestDataDao{
						CompanyNumber: "87654321",
					},
				}

				mockDaoReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
				mockDaoReqService.EXPECT().GetAuthCodeRequest("123").Return(&authCodeDaoResponse, nil)
				mockDaoReqService.EXPECT().UpdateAuthCodeRequestOfficer(gomock.Any()).Return(fmt.Errorf("error"))

				httpmock.Activate()
				defer httpmock.DeactivateAndReset()
				responder := httpmock.NewStringResponder(http.StatusOK, `{"total_results":1}`)
				httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/eligible-officers/98765432", responder)

				res := serveUpdateAuthCodeRequestHandler(context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}), t, &models.AuthCodeRequest{CompanyNumber: "87654321", OfficerID: "98765432"}, "123", nil, mockDaoReqService)
				So(res.Code, ShouldEqual, http.StatusInternalServerError)
				So(res.Body.String(), ShouldStartWith, `{"message":"error updating officer details in authcode request"}`)
			})

			Convey("successful officer update", func() {
				mockCtrl := gomock.NewController(t)
				defer mockCtrl.Finish()

				authCodeDaoResponse := models.AuthCodeRequestResourceDao{
					Data: models.AuthCodeRequestDataDao{
						CompanyNumber: "87654321",
					},
				}

				mockDaoReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
				mockDaoReqService.EXPECT().GetAuthCodeRequest("123").Return(&authCodeDaoResponse, nil)
				mockDaoReqService.EXPECT().UpdateAuthCodeRequestOfficer(gomock.Any()).Return(nil)

				httpmock.Activate()
				defer httpmock.DeactivateAndReset()
				responder := httpmock.NewStringResponder(http.StatusOK, `{"total_results":1}`)
				httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/eligible-officers/98765432", responder)

				res := serveUpdateAuthCodeRequestHandler(context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}), t, &models.AuthCodeRequest{CompanyNumber: "87654321", OfficerID: "98765432"}, "123", nil, mockDaoReqService)
				So(res.Code, ShouldEqual, http.StatusCreated)
				So(res.Body.String(), ShouldContainSubstring, `"company_number":"87654321"`)
			})

		})

		Convey("status update", func() {

			Convey("no officer details", func() {
				mockCtrl := gomock.NewController(t)
				defer mockCtrl.Finish()

				authCodeDaoResponse := models.AuthCodeRequestResourceDao{
					Data: models.AuthCodeRequestDataDao{
						CompanyNumber: "87654321",
					},
				}

				mockDaoReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
				mockDaoReqService.EXPECT().GetAuthCodeRequest("123").Return(&authCodeDaoResponse, nil)

				res := serveUpdateAuthCodeRequestHandler(context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}), t, &models.AuthCodeRequest{CompanyNumber: "87654321", Status: "submitted"}, "123", nil, mockDaoReqService)
				So(res.Code, ShouldEqual, http.StatusBadRequest)
				So(res.Body.String(), ShouldStartWith, `{"message":"officer details not supplied"}`)
			})

			Convey("error retrieving authcode", func() {
				mockCtrl := gomock.NewController(t)
				defer mockCtrl.Finish()

				authCodeDaoResponse := models.AuthCodeRequestResourceDao{
					Data: models.AuthCodeRequestDataDao{
						CompanyNumber: "87654321",
						OfficerID:     "321",
					},
				}

				mockDaoReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
				mockDaoReqService.EXPECT().GetAuthCodeRequest("123").Return(&authCodeDaoResponse, nil)

				mockDaoAuthcodeService := mocks.NewMockAuthcodeDAOService(mockCtrl)
				mockDaoAuthcodeService.EXPECT().CompanyHasAuthCode(gomock.Any()).Return(false, fmt.Errorf("error"))

				res := serveUpdateAuthCodeRequestHandler(context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}), t, &models.AuthCodeRequest{CompanyNumber: "87654321", Status: "submitted"}, "123", mockDaoAuthcodeService, mockDaoReqService)
				So(res.Code, ShouldEqual, http.StatusInternalServerError)
				So(res.Body.String(), ShouldStartWith, `{"message":"error retrieving Auth Code from DB"}`)
			})

			Convey("error sending status queue item", func() {
				mockCtrl := gomock.NewController(t)
				defer mockCtrl.Finish()

				authCodeDaoResponse := models.AuthCodeRequestResourceDao{
					Data: models.AuthCodeRequestDataDao{
						CompanyNumber: "87654321",
						OfficerID:     "321",
					},
				}

				mockDaoReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
				mockDaoReqService.EXPECT().GetAuthCodeRequest("123").Return(&authCodeDaoResponse, nil)

				mockDaoAuthcodeService := mocks.NewMockAuthcodeDAOService(mockCtrl)
				mockDaoAuthcodeService.EXPECT().CompanyHasAuthCode(gomock.Any()).Return(false, nil)

				res := serveUpdateAuthCodeRequestHandler(context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}), t, &models.AuthCodeRequest{CompanyNumber: "87654321", Status: "submitted"}, "123", mockDaoAuthcodeService, mockDaoReqService)
				So(res.Code, ShouldEqual, http.StatusInternalServerError)
				So(res.Body.String(), ShouldStartWith, `{"message":"error sending queue item"}`)
			})

			Convey("officer not found", func() {
				mockCtrl := gomock.NewController(t)
				defer mockCtrl.Finish()

				authCodeDaoResponse := models.AuthCodeRequestResourceDao{
					Data: models.AuthCodeRequestDataDao{
						CompanyNumber: "87654321",
						OfficerID:     "321",
					},
				}

				mockDaoReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
				mockDaoReqService.EXPECT().GetAuthCodeRequest("123").Return(&authCodeDaoResponse, nil)

				mockDaoAuthcodeService := mocks.NewMockAuthcodeDAOService(mockCtrl)
				mockDaoAuthcodeService.EXPECT().CompanyHasAuthCode(gomock.Any()).Return(false, nil)

				httpmock.Activate()
				defer httpmock.DeactivateAndReset()
				responder := httpmock.NewStringResponder(http.StatusNotFound, "")
				httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/eligible-officers/321", responder)

				res := serveUpdateAuthCodeRequestHandler(context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}), t, &models.AuthCodeRequest{CompanyNumber: "87654321", Status: "submitted"}, "123", mockDaoAuthcodeService, mockDaoReqService)
				So(res.Code, ShouldEqual, http.StatusNotFound)
				So(res.Body.String(), ShouldStartWith, `{"message":"officer not found"}`)

			})

			Convey("error updating status", func() {
				mockCtrl := gomock.NewController(t)
				defer mockCtrl.Finish()

				authCodeDaoResponse := models.AuthCodeRequestResourceDao{
					Data: models.AuthCodeRequestDataDao{
						CompanyNumber: "87654321",
						OfficerID:     "321",
					},
				}

				mockDaoReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
				mockDaoReqService.EXPECT().GetAuthCodeRequest("123").Return(&authCodeDaoResponse, nil)
				mockDaoReqService.EXPECT().UpdateAuthCodeRequestStatus(gomock.Any()).Return(fmt.Errorf("error"))

				mockDaoAuthcodeService := mocks.NewMockAuthcodeDAOService(mockCtrl)
				mockDaoAuthcodeService.EXPECT().CompanyHasAuthCode(gomock.Any()).Return(false, nil)

				httpmock.Activate()
				defer httpmock.DeactivateAndReset()
				responder := httpmock.NewStringResponder(http.StatusOK, `{"total_results":1}`)
				httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/eligible-officers/321", responder)
				queueAPIResponder := httpmock.NewStringResponder(http.StatusOK, `{}`)
				httpmock.RegisterResponder(http.MethodPost, "/api/queue/authcode", queueAPIResponder)

				res := serveUpdateAuthCodeRequestHandler(context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}), t, &models.AuthCodeRequest{CompanyNumber: "87654321", Status: "submitted"}, "123", mockDaoAuthcodeService, mockDaoReqService)
				So(res.Code, ShouldEqual, http.StatusInternalServerError)
				So(res.Body.String(), ShouldStartWith, `{"message":"error updating status"}`)
			})

			Convey("successful status update", func() {
				mockCtrl := gomock.NewController(t)
				defer mockCtrl.Finish()

				authCodeDaoResponse := models.AuthCodeRequestResourceDao{
					Data: models.AuthCodeRequestDataDao{
						CompanyNumber: "87654321",
						OfficerID:     "321",
					},
				}

				mockDaoReqService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
				mockDaoReqService.EXPECT().GetAuthCodeRequest("123").Return(&authCodeDaoResponse, nil)
				mockDaoReqService.EXPECT().UpdateAuthCodeRequestStatus(gomock.Any()).Return(nil)

				mockDaoAuthcodeService := mocks.NewMockAuthcodeDAOService(mockCtrl)
				mockDaoAuthcodeService.EXPECT().CompanyHasAuthCode(gomock.Any()).Return(false, nil)

				httpmock.Activate()
				defer httpmock.DeactivateAndReset()
				responder := httpmock.NewStringResponder(http.StatusOK, `{"total_results":1}`)
				httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/eligible-officers/321", responder)
				queueAPIResponder := httpmock.NewStringResponder(http.StatusOK, `{}`)
				httpmock.RegisterResponder(http.MethodPost, "/api/queue/authcode", queueAPIResponder)

				res := serveUpdateAuthCodeRequestHandler(context.WithValue(context.Background(), authentication.ContextKeyUserDetails, authentication.AuthUserDetails{}), t, &models.AuthCodeRequest{CompanyNumber: "87654321", Status: "submitted"}, "123", mockDaoAuthcodeService, mockDaoReqService)
				So(res.Code, ShouldEqual, http.StatusCreated)
				So(res.Body.String(), ShouldContainSubstring, `"company_number":"87654321"`)
			})
		})
	})
}
