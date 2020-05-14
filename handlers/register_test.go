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

		mockOfficerService := mocks.NewMockOfficerDAOService(mockCtrl)
		mockAuthcodeService := mocks.NewMockAuthcodeDAOService(mockCtrl)
		Register(router, &config.Config{}, mockAuthcodeService, mockOfficerService)

		So(router.GetRoute("get-company-officers"), ShouldNotBeNil)
		So(router.GetRoute("create-auth-code-request"), ShouldNotBeNil)
	})
}
