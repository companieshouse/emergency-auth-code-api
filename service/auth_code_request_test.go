package service

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/companieshouse/emergency-auth-code-api/config"
	"github.com/companieshouse/emergency-auth-code-api/mocks"
	"github.com/companieshouse/emergency-auth-code-api/models"
	"github.com/companieshouse/emergency-auth-code-api/oracle"
	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	. "github.com/smartystreets/goconvey/convey"
)

const authCodeRequestID = "123"
const companyNumber = "87654321"
const testRequestID = "xyz123"

func TestUnitUpdateAuthCodeRequestOfficer(t *testing.T) {

	Convey("Update Auth Code Request Officer", t, func() {

		Convey("error updating authcode request", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockDaoService.EXPECT().UpdateAuthCodeRequestOfficer(gomock.Any()).Return(fmt.Errorf("error"))
			svc := AuthCodeRequestService{DAO: mockDaoService}

			authCodeReq := models.AuthCodeRequestResourceDao{}
			officer := oracle.Officer{}

			responseType := svc.UpdateAuthCodeRequestOfficer(&authCodeReq, authCodeRequestID, &officer)
			So(responseType, ShouldEqual, Error)
		})

		Convey("officer update - success", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockDaoService.EXPECT().UpdateAuthCodeRequestOfficer(gomock.Any()).Return(nil)
			svc := AuthCodeRequestService{DAO: mockDaoService}

			authCodeReq := models.AuthCodeRequestResourceDao{}
			officer := oracle.Officer{}

			responseType := svc.UpdateAuthCodeRequestOfficer(&authCodeReq, authCodeRequestID, &officer)
			So(responseType, ShouldEqual, Success)
		})
	})
}

func TestUnitUpdateAuthCodeRequestStatus(t *testing.T) {
	Convey("Update Auth Code Request Status", t, func() {

		Convey("error updating authcode request", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockDaoService.EXPECT().UpdateAuthCodeRequestStatus(gomock.Any()).Return(fmt.Errorf("error"))
			svc := AuthCodeRequestService{DAO: mockDaoService}

			authCodeReq := models.AuthCodeRequestResourceDao{}

			responseType := svc.UpdateAuthCodeRequestStatusSubmitted(&authCodeReq, authCodeRequestID, false)
			So(responseType, ShouldEqual, Error)
		})

		Convey("status update - success", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockDaoService.EXPECT().UpdateAuthCodeRequestStatus(gomock.Any()).Return(nil)
			svc := AuthCodeRequestService{DAO: mockDaoService}

			authCodeReq := models.AuthCodeRequestResourceDao{}

			responseType := svc.UpdateAuthCodeRequestStatusSubmitted(&authCodeReq, authCodeRequestID, false)
			So(responseType, ShouldEqual, Success)
		})
	})
}

func TestUnitSendAuthCodeRequestQueueAPIAuthCodeFlowErrors(t *testing.T) {
	Convey("send auth code request", t, func() {
		// build test config
		cfg, _ := config.Get()
		cfg.NewAuthCodeAPIFlow = false
		cfg.QueueAPILocalURL = "http://local.test"
		cfg.QueueAPILocalPath = "/api/queue/authcode"

		Convey("error getting officer details", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			svc := AuthCodeRequestService{
				DAO:    mockDaoService,
				Config: cfg,
			}

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			responder := httpmock.NewStringResponder(http.StatusInternalServerError, "")
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/eligible-officers/123", responder)

			authCodeReq := models.AuthCodeRequestResourceDao{
				Data: models.AuthCodeRequestDataDao{
					OfficerID: "987",
				},
			}

			responseType := svc.SendAuthCodeRequest(&authCodeReq, companyNumber, "email@companieshouse.gov.uk", testRequestID, true)
			So(responseType, ShouldEqual, Error)

		})

		Convey("officer not found", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			svc := AuthCodeRequestService{
				DAO:    mockDaoService,
				Config: cfg,
			}

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			responder := httpmock.NewStringResponder(http.StatusNotFound, "")
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/eligible-officers/987", responder)

			authCodeReq := models.AuthCodeRequestResourceDao{
				Data: models.AuthCodeRequestDataDao{
					OfficerID: "987",
				},
			}

			responseType := svc.SendAuthCodeRequest(&authCodeReq, companyNumber, "email@companieshouse.gov.uk", testRequestID, true)
			So(responseType, ShouldEqual, NotFound)
		})

		Convey("error sending queue item", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			svc := AuthCodeRequestService{
				DAO:    mockDaoService,
				Config: cfg,
			}

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			responder := httpmock.NewStringResponder(http.StatusOK, `{"surname":"bloggs"}`)
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/eligible-officers/987", responder)

			authCodeReq := models.AuthCodeRequestResourceDao{
				Data: models.AuthCodeRequestDataDao{
					OfficerID: "987",
				},
			}

			responseType := svc.SendAuthCodeRequest(&authCodeReq, companyNumber, "email@companieshouse.gov.uk", testRequestID, true)
			So(responseType, ShouldEqual, Error)
		})
	})
}

func TestUnitSendAuthCodeRequestQueueAPIAuthCodeFlowSuccess(t *testing.T) {
	Convey("send auth code request", t, func() {
		Convey("send auth code request - success", func() {
			// build test config
			cfg, _ := config.Get()
			cfg.NewAuthCodeAPIFlow = false
			cfg.QueueAPILocalURL = "http://local.test"
			cfg.QueueAPILocalPath = "/api/queue/authcode"

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			svc := AuthCodeRequestService{
				DAO:    mockDaoService,
				Config: cfg,
			}

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			responder := httpmock.NewStringResponder(http.StatusOK, `{"forename":"joe","surname":"bloggs"}`)
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/eligible-officers/987", responder)
			queueAPIResponder := httpmock.NewStringResponder(http.StatusOK, `{}`)
			httpmock.RegisterResponder(http.MethodPost, cfg.QueueAPILocalPath, queueAPIResponder)

			authCodeReq := models.AuthCodeRequestResourceDao{
				Data: models.AuthCodeRequestDataDao{
					OfficerID: "987",
				},
			}

			responseType := svc.SendAuthCodeRequest(&authCodeReq, companyNumber, "email@companieshouse.gov.uk", testRequestID, true)
			So(responseType, ShouldEqual, Success)
		})
	})
}

func TestUnitSendAuthCodeRequestAuthCodeAPIAuthCodeFlow(t *testing.T) {
	Convey("send auth code request", t, func() {
		Convey("send auth code request - success", func() {
			// build test config
			cfg, _ := config.Get()
			cfg.NewAuthCodeAPIFlow = true
			cfg.AuthCodeAPILocalURL = "http://local.test"
			cfg.AuthCodeAPILocalPath = "/private/company/%s/authcode/request"

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			svc := AuthCodeRequestService{
				DAO:    mockDaoService,
				Config: cfg,
			}

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			responder := httpmock.NewStringResponder(http.StatusOK, `{"forename":"joe","surname":"bloggs"}`)
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/eligible-officers/987", responder)
			queueAPIResponder := httpmock.NewStringResponder(http.StatusOK, `{}`)
			httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf(cfg.AuthCodeAPILocalPath, "87654321"), queueAPIResponder)

			authCodeReq := models.AuthCodeRequestResourceDao{
				Data: models.AuthCodeRequestDataDao{
					OfficerID: "987",
				},
			}

			responseType := svc.SendAuthCodeRequest(&authCodeReq, companyNumber, "email@companieshouse.gov.uk", testRequestID, true)
			So(responseType, ShouldEqual, Success)
		})
	})
}

func TestUnitGetAuthCodeReqDao(t *testing.T) {
	Convey("Get Auth Code Request DAO", t, func() {
		Convey("error getting request", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockDaoService.EXPECT().GetAuthCodeRequest(gomock.Any()).Return(nil, fmt.Errorf("error"))
			svc := AuthCodeRequestService{DAO: mockDaoService}

			request, responseType := svc.GetAuthCodeReqDao(authCodeRequestID, companyNumber)
			So(request, ShouldBeNil)
			So(responseType, ShouldEqual, Error)
		})

		Convey("request not found", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockDaoService.EXPECT().GetAuthCodeRequest(gomock.Any()).Return(nil, nil)
			svc := AuthCodeRequestService{DAO: mockDaoService}

			request, responseType := svc.GetAuthCodeReqDao(authCodeRequestID, companyNumber)
			So(request, ShouldBeNil)
			So(responseType, ShouldEqual, NotFound)
		})

		Convey("company number mismatch", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			response := models.AuthCodeRequestResourceDao{
				Data: models.AuthCodeRequestDataDao{
					CompanyNumber: "mismatch",
				},
			}
			mockDaoService.EXPECT().GetAuthCodeRequest(gomock.Any()).Return(&response, nil)
			svc := AuthCodeRequestService{DAO: mockDaoService}

			request, responseType := svc.GetAuthCodeReqDao(authCodeRequestID, companyNumber)
			So(request, ShouldBeNil)
			So(responseType, ShouldEqual, InvalidData)
		})

		Convey("get auth code request - success", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			response := models.AuthCodeRequestResourceDao{
				Data: models.AuthCodeRequestDataDao{
					CompanyNumber: companyNumber,
				},
			}
			mockDaoService.EXPECT().GetAuthCodeRequest(gomock.Any()).Return(&response, nil)
			svc := AuthCodeRequestService{DAO: mockDaoService}

			request, responseType := svc.GetAuthCodeReqDao(authCodeRequestID, companyNumber)
			So(request.Data.CompanyNumber, ShouldEqual, companyNumber)
			So(responseType, ShouldEqual, Success)
		})
	})
}

func TestUnitGetLetterType(t *testing.T) {
	Convey("get letter type", t, func() {
		So(getLetterType(true), ShouldEqual, "reminder")
		So(getLetterType(false), ShouldEqual, "apply")
	})
}

func TestUnitCheckMultipleCorporateBodySubmissions(t *testing.T) {

	Convey("Check Multiple Corporate Body Submissions", t, func() {
		Convey("Error checking submissions", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			errorMessage := "error test"
			mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockDaoService.EXPECT().CheckMultipleCorporateBodySubmissions(gomock.Any()).Return(false, fmt.Errorf(errorMessage))
			svc := AuthCodeRequestService{DAO: mockDaoService}

			response, err := svc.CheckMultipleCorporateBodySubmissions(companyNumber)
			So(response, ShouldBeFalse)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, errorMessage)
		})

		Convey("Company has multiple submissions", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockDaoService.EXPECT().CheckMultipleCorporateBodySubmissions(gomock.Any()).Return(true, nil)
			svc := AuthCodeRequestService{DAO: mockDaoService}

			response, err := svc.CheckMultipleCorporateBodySubmissions(companyNumber)
			So(response, ShouldBeTrue)
			So(err, ShouldBeNil)
		})
	})
}

func TestUnitCheckMultipleUserSubmissions(t *testing.T) {

	Convey("Check Multiple User Submissions", t, func() {
		Convey("Error checking submissions", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			errorMessage := "error test"
			mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockDaoService.EXPECT().CheckMultipleUserSubmissions(gomock.Any()).Return(false, fmt.Errorf(errorMessage))
			svc := AuthCodeRequestService{DAO: mockDaoService}

			response, err := svc.CheckMultipleUserSubmissions(companyNumber)
			So(response, ShouldBeFalse)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, errorMessage)
		})

		Convey("Company has multiple submissions", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockDaoService.EXPECT().CheckMultipleUserSubmissions(gomock.Any()).Return(true, nil)
			svc := AuthCodeRequestService{DAO: mockDaoService}

			response, err := svc.CheckMultipleUserSubmissions(companyNumber)
			So(response, ShouldBeTrue)
			So(err, ShouldBeNil)
		})
	})
}
