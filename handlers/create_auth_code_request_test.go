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

func serveGetPaymentDetailsHandler(
	ctx context.Context,
	t *testing.T,
	reqBody *models.AuthCodeRequest,
	daoSvc dao.Service) (*httptest.ResponseRecorder, *models.ResponseResource) {

	svc := &service.AuthCodeService{}

	if daoSvc != nil {
		svc.DAO = daoSvc
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

	h := CreateAuthCodeRequest(svc)
	req := httptest.NewRequest(http.MethodPost, "/", body).WithContext(ctx)
	res := httptest.NewRecorder()

	h.ServeHTTP(res, req.WithContext(ctx))

	if res.Body.Len() > 0 {
		var responseBody models.ResponseResource
		err := json.NewDecoder(res.Body).Decode(&responseBody)
		if err != nil {
			t.Errorf("failed to read response body")
		}

		return res, &responseBody
	}

	return res, nil
}

func TestUnitCreateAuthCodeRequestHandler(t *testing.T) {
	Convey("CreateAuthCodeRequestHandler tests", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		Convey("authcode resource must be in context", func() {
			res, body := serveGetPaymentDetailsHandler(context.Background(), t, nil, nil)

			So(res.Code, ShouldEqual, http.StatusBadRequest)
			So(body.Message, ShouldEqual, "failed to read request body")
		})

		Convey("company number missing from request", func() {
			res, body := serveGetPaymentDetailsHandler(context.Background(), t, &models.AuthCodeRequest{}, nil)
			So(res.Code, ShouldEqual, http.StatusBadRequest)
			So(body.Message, ShouldEqual, "company number missing from request")
		})

		Convey("error checking DB for authcode", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// stub the DB lookup
			mockService := mocks.NewMockService(mockCtrl)
			mockService.EXPECT().CompanyHasAuthCode(gomock.Any()).Return(false, fmt.Errorf("error"))

			res, body := serveGetPaymentDetailsHandler(context.Background(), t, &models.AuthCodeRequest{CompanyNumber: "87654321"}, mockService)
			So(res.Code, ShouldEqual, http.StatusInternalServerError)
			So(body, ShouldBeNil)
		})

		Convey("successful Authcode Reminder", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// stub the DB lookup
			mockService := mocks.NewMockService(mockCtrl)
			mockService.EXPECT().CompanyHasAuthCode(gomock.Any()).Return(true, nil)

			res, body := serveGetPaymentDetailsHandler(context.Background(), t, &models.AuthCodeRequest{CompanyNumber: "87654321"}, mockService)
			So(res.Code, ShouldEqual, http.StatusCreated)
			So(body, ShouldBeNil) // TODO check body.Message when implemented
		})

		Convey("successful Authcode Request", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// stub the DB lookup
			mockService := mocks.NewMockService(mockCtrl)
			mockService.EXPECT().CompanyHasAuthCode(gomock.Any()).Return(false, nil)

			res, body := serveGetPaymentDetailsHandler(context.Background(), t, &models.AuthCodeRequest{CompanyNumber: "87654321"}, mockService)
			So(res.Code, ShouldEqual, http.StatusCreated)
			So(body, ShouldBeNil) // TODO check body.Message when implemented
		})
	})
}
