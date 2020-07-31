package transformers

import (
	"testing"

	"github.com/companieshouse/emergency-auth-code-api/models"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitAuthCodeResourceRequestToDB(t *testing.T) {

	Convey("Auth code details from REST call are transformed and ID and Etag are generated", t, func() {
		req := &models.AuthCodeRequest{
			CompanyNumber: "12345678",
			OfficerID:     "87654321",
		}
		dao := AuthCodeResourceRequestToDB(req)

		So(dao.Data.OfficerID, ShouldEqual, "87654321")
		So(dao.Data.CompanyNumber, ShouldEqual, "12345678")
		So(dao.ID, ShouldHaveLength, 15)
		So(dao.Data.Etag, ShouldNotBeNil)
	})
}

func TestUnitAuthCodeRequestResourceDaoToResponse(t *testing.T) {

	Convey("Auth code details from DB are transformed back to REST call", t, func() {
		req := &models.AuthCodeRequestResourceDao{
			Data: models.AuthCodeRequestDataDao{
				CompanyNumber: "12345678",
				CompanyName:   "test",
				OfficerID:     "87654321",
			},
		}

		response := AuthCodeRequestResourceDaoToResponse(req)

		So(response.CompanyNumber, ShouldEqual, "12345678")
		So(response.CompanyName, ShouldEqual, "test")
		So(response.OfficerID, ShouldEqual, "87654321")
	})
}
