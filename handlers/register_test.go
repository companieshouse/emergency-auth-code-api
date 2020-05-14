package handlers

import (
	"testing"

	"github.com/companieshouse/emergency-auth-code-api/service"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitRegisterRoutes(t *testing.T) {
	Convey("Register routes", t, func() {
		router := mux.NewRouter()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockService := service.NewMockDirectorDatabase(mockCtrl)
		Register(router, mockService)

		So(router.GetRoute("get-company-directors"), ShouldNotBeNil)
		So(router.GetRoute("create-auth-code-request"), ShouldNotBeNil)
	})
}
