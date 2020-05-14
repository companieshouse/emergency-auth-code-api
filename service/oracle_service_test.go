package service

//import (
//	"fmt"
//	"testing"
//
//	"github.com/companieshouse/emergency-auth-code-api/config"
//
//	"github.com/jarcoal/httpmock"
//
//	"github.com/DATA-DOG/go-sqlmock"
//	"github.com/golang/mock/gomock"
//
//	. "github.com/smartystreets/goconvey/convey"
//)
//
//const companyNumber = "12345678"
//
//func createMockOracleService(cfg *config.Config) (OracleService, sqlmock.Sqlmock) {
//	db, mock, _ := sqlmock.New()
//	return OracleService{DatabaseDriver: db}, mock
//}
//
//func TestUnitGetCompanyDirectors(t *testing.T) {
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//	cfg, _ := config.Get()
//	cfg.DirectorDatabaseUrl = "databaseurl"
//	cfg.DirectorDatabaseUsername = "databaseusername"
//	cfg.DirectorDatabasePassword = "databasepassword"
//
//	Convey("Successfully return directors for given company", t, func() {
//		// create mock sql driver
//		mockOracleService, sqlMock := createMockOracleService(cfg)
//		defer mockOracleService.DatabaseDriver.Close()
//
//		httpmock.Activate()
//		defer httpmock.DeactivateAndReset()
//		httpmock.RegisterResponder("GET", cfg.DirectorDatabaseUsername+"/"+cfg.DirectorDatabasePassword+"@//"+cfg.DirectorDatabaseUrl, nil)
//
//		sqlMock.ExpectQuery(fmt.Sprintf(QUERY, companyNumber)).WillReturnRows(nil)
//		companyOfficers, _ := mockOracleService.GetCompanyDirectors(companyNumber)
//
//		So(companyOfficers, ShouldNotBeNil)
//	})
//}
