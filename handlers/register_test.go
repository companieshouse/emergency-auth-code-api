package handlers

import (
	"testing"

	"github.com/companieshouse/emergency-auth-code-api/config"
	"github.com/companieshouse/emergency-auth-code-api/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitRegisterRoutes(t *testing.T) {
	Convey("Register routes", t, func() {
		router := mux.NewRouter()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockAuthcodeService := mocks.NewMockAuthcodeDAOService(mockCtrl)
		mockAuthcodeRequestService := mocks.NewMockAuthcodeRequestDAOService(mockCtrl)
		Register(router, &config.Config{}, mockAuthcodeService, mockAuthcodeRequestService)

		So(router.GetRoute("get-company-directors"), ShouldNotBeNil)
		So(router.GetRoute("create-auth-code-request"), ShouldNotBeNil)
	})
}
