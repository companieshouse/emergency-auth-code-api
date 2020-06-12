package service

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/companieshouse/emergency-auth-code-api/mocks"
	"github.com/companieshouse/emergency-auth-code-api/models"
	"github.com/companieshouse/emergency-auth-code-api/oracle"
	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	. "github.com/smartystreets/goconvey/convey"
)

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

			_, responseType := svc.UpdateAuthCodeRequestOfficer(&authCodeReq, "123", &officer)
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

			_, responseType := svc.UpdateAuthCodeRequestOfficer(&authCodeReq, "123", &officer)
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

			_, responseType := svc.UpdateAuthCodeRequestStatus(&authCodeReq, "123", "submitted", false)
			So(responseType, ShouldEqual, Error)
		})

		Convey("status update - success", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockDaoService.EXPECT().UpdateAuthCodeRequestStatus(gomock.Any()).Return(nil)
			svc := AuthCodeRequestService{DAO: mockDaoService}

			authCodeReq := models.AuthCodeRequestResourceDao{}

			_, responseType := svc.UpdateAuthCodeRequestStatus(&authCodeReq, "123", "submitted", false)
			So(responseType, ShouldEqual, Success)
		})
	})
}

func TestUnitSendAuthCodeRequest(t *testing.T) {

	Convey("send auth code request", t, func() {

		Convey("error getting officer details", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			svc := AuthCodeRequestService{DAO: mockDaoService}

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			responder := httpmock.NewStringResponder(http.StatusInternalServerError, "")
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/eligible-officers/123", responder)

			authCodeReq := models.AuthCodeRequestResourceDao{
				Data: models.AuthCodeRequestDataDao{
					OfficerID: "987",
				},
			}

			responseType := svc.SendAuthCodeRequest(&authCodeReq, "87654321", "email@companieshouse.gov.uk", true)
			So(responseType, ShouldEqual, Error)

		})

		Convey("officer not found", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			svc := AuthCodeRequestService{DAO: mockDaoService}

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			responder := httpmock.NewStringResponder(http.StatusNotFound, "")
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/eligible-officers/987", responder)

			authCodeReq := models.AuthCodeRequestResourceDao{
				Data: models.AuthCodeRequestDataDao{
					OfficerID: "987",
				},
			}

			responseType := svc.SendAuthCodeRequest(&authCodeReq, "87654321", "email@companieshouse.gov.uk", true)
			So(responseType, ShouldEqual, NotFound)
		})

		Convey("error sending queue item", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			svc := AuthCodeRequestService{DAO: mockDaoService}

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			responder := httpmock.NewStringResponder(http.StatusOK, `{"surname":"bloggs"}`)
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/eligible-officers/987", responder)

			authCodeReq := models.AuthCodeRequestResourceDao{
				Data: models.AuthCodeRequestDataDao{
					OfficerID: "987",
				},
			}

			responseType := svc.SendAuthCodeRequest(&authCodeReq, "87654321", "email@companieshouse.gov.uk", true)
			So(responseType, ShouldEqual, Error)
		})

		Convey("send auth code request - success", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			svc := AuthCodeRequestService{DAO: mockDaoService}

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			responder := httpmock.NewStringResponder(http.StatusOK, `{"forename":"joe","surname":"bloggs"}`)
			httpmock.RegisterResponder(http.MethodGet, "/emergency-auth-code/company/87654321/eligible-officers/987", responder)
			queueAPIResponder := httpmock.NewStringResponder(http.StatusOK, `{}`)
			httpmock.RegisterResponder(http.MethodPost, "/api/queue/authcode", queueAPIResponder)

			authCodeReq := models.AuthCodeRequestResourceDao{
				Data: models.AuthCodeRequestDataDao{
					OfficerID: "987",
				},
			}

			responseType := svc.SendAuthCodeRequest(&authCodeReq, "87654321", "email@companieshouse.gov.uk", true)
			So(responseType, ShouldEqual, Success)
		})
	})
}

func TestUnitGetAuthCodeReqDao(t *testing.T) {
	companyNumber := "87654321"
	Convey("Get Auth Code Request DAO", t, func() {
		Convey("error getting request", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockDaoService.EXPECT().GetAuthCodeRequest(gomock.Any()).Return(nil, fmt.Errorf("error"))
			svc := AuthCodeRequestService{DAO: mockDaoService}

			request, responseType := svc.GetAuthCodeReqDao("123", companyNumber)
			So(request, ShouldBeNil)
			So(responseType, ShouldEqual, Error)
		})

		Convey("request not found", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDaoService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
			mockDaoService.EXPECT().GetAuthCodeRequest(gomock.Any()).Return(nil, nil)
			svc := AuthCodeRequestService{DAO: mockDaoService}

			request, responseType := svc.GetAuthCodeReqDao("123", companyNumber)
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

			request, responseType := svc.GetAuthCodeReqDao("123", companyNumber)
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

			request, responseType := svc.GetAuthCodeReqDao("123", companyNumber)
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
