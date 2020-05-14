package service

import (
	"errors"
	"testing"

	"github.com/companieshouse/emergency-auth-code-api/mocks"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitCheckAuthCode(t *testing.T) {

	Convey("Check Auth Code Service", t, func() {
		Convey("Error checking authcode in DB", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDaoService := mocks.NewMockAuthcodeDAOService(mockCtrl)
			mockDaoService.EXPECT().CompanyHasAuthCode(gomock.Any()).Return(false, errors.New("error"))
			svc := AuthCodeService{DAO: mockDaoService}

			_, err := svc.CheckAuthCodeExists("87654321")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "error checking DB for auth code: [error]")
		})

		Convey("Company does not have authcode", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDaoService := mocks.NewMockAuthcodeDAOService(mockCtrl)
			mockDaoService.EXPECT().CompanyHasAuthCode(gomock.Any()).Return(false, nil)
			svc := AuthCodeService{DAO: mockDaoService}

			companyHasAuthCode, err := svc.CheckAuthCodeExists("87654321")
			So(err, ShouldBeNil)
			So(companyHasAuthCode, ShouldBeFalse)
		})

		Convey("Company has authcode", func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDaoService := mocks.NewMockAuthcodeDAOService(mockCtrl)
			mockDaoService.EXPECT().CompanyHasAuthCode(gomock.Any()).Return(true, nil)
			svc := AuthCodeService{DAO: mockDaoService}

			companyHasAuthCode, err := svc.CheckAuthCodeExists("87654321")
			So(err, ShouldBeNil)
			So(companyHasAuthCode, ShouldBeTrue)
		})
	})
}
