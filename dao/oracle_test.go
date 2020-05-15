package dao

import (
	"errors"
	"fmt"
	"testing"

	"github.com/companieshouse/emergency-auth-code-api/config"

	"github.com/jarcoal/httpmock"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"
)

const companyNumber = "12345678"

func createMockOracleService() (OracleService, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	return OracleService{db: db}, mock
}

func mockedRows() *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"id", "forename1", "forename2", "surname"}).
		AddRow("1", "forenametest1", "forenametest2", "surnametest1").
		AddRow("2", "forenametest3", "forenametest4", "surnametest2")

	return rows
}

func TestUnitGetCompanyOfficers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	cfg, _ := config.Get()
	cfg.OfficerDatabaseURL = "databaseurl"
	cfg.OfficerDatabaseUsername = "databaseusername"
	cfg.OfficerDatabasePassword = "databasepassword"

	Convey("Failed to run query for company officers on database", t, func() {
		// create mock sql driver
		mockOracleService, sqlMock := createMockOracleService()
		defer mockOracleService.db.Close()

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", cfg.OfficerDatabaseUsername+"/"+cfg.OfficerDatabasePassword+"@//"+cfg.OfficerDatabaseURL, nil)

		sqlMock.ExpectQuery(fmt.Sprintf(query, companyNumber)).WillReturnError(fmt.Errorf("error"))
		companyOfficers, err := mockOracleService.GetCompanyOfficers(companyNumber)

		So(companyOfficers, ShouldEqual, nil)
		So(err, ShouldBeError, errors.New("failed to run query for company officers on database"))
	})

	Convey("Error reading rows returned from database query", t, func() {
		// create mock sql driver
		mockOracleService, sqlMock := createMockOracleService()
		defer mockOracleService.db.Close()

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", cfg.OfficerDatabaseUsername+"/"+cfg.OfficerDatabasePassword+"@//"+cfg.OfficerDatabaseURL, nil)

		// Deliberately return too many variables to throw sql scanning error
		sqlMock.ExpectQuery(fmt.Sprintf(query, companyNumber)).WillReturnRows(sqlmock.
			NewRows([]string{"id", "forename1", "forename2", "surname", "error"}).
			AddRow("id", "forename", "forename2", "surname", "error"))
		companyOfficers, err := mockOracleService.GetCompanyOfficers(companyNumber)

		So(companyOfficers, ShouldEqual, nil)
		So(err, ShouldBeError, errors.New("error reading rows returned from database query"))
	})

	Convey("Successfully return officers for given company", t, func() {
		// create mock sql driver
		mockOracleService, sqlMock := createMockOracleService()
		defer mockOracleService.db.Close()

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", cfg.OfficerDatabaseUsername+"/"+cfg.OfficerDatabasePassword+"@//"+cfg.OfficerDatabaseURL, nil)

		sqlMock.ExpectQuery(fmt.Sprintf(query, companyNumber)).WillReturnRows(mockedRows())
		companyOfficers, err := mockOracleService.GetCompanyOfficers(companyNumber)

		So(companyOfficers.TotalCount, ShouldEqual, 2)
		So(companyOfficers.Items[0].ID, ShouldEqual, "1")
		So(companyOfficers.Items[1].ID, ShouldEqual, "2")
		So(companyOfficers.Items[0].Forename1, ShouldEqual, "forenametest1")
		So(companyOfficers.Items[1].Forename1, ShouldEqual, "forenametest3")
		So(companyOfficers.Items[0].Forename2, ShouldEqual, "forenametest2")
		So(companyOfficers.Items[1].Forename2, ShouldEqual, "forenametest4")
		So(companyOfficers.Items[0].Surname, ShouldEqual, "surnametest1")
		So(companyOfficers.Items[1].Surname, ShouldEqual, "surnametest2")
		So(err, ShouldEqual, nil)
	})
}
