package service

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitGetOfficers(t *testing.T) {
	companyNumber := "87654321"

	Convey("Get Officer List", t, func() {

		url := "/emergency-auth-code/company/" + companyNumber + "/eligible-officers"

		Convey("No response", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			resp, respType, err := GetOfficers(companyNumber)
			So(resp, ShouldBeNil)
			So(respType, ShouldEqual, Error)
			So(err, ShouldBeError, "Get "+url+": no responder found")
		})

		Convey("Empty response", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			responder := httpmock.NewStringResponder(http.StatusNotFound, "")
			httpmock.RegisterResponder(http.MethodGet, url, responder)

			resp, respType, err := GetOfficers(companyNumber)
			So(resp, ShouldBeNil)
			So(respType, ShouldEqual, NotFound)
			So(err, ShouldBeNil)
		})

		Convey("Successful response", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			responder := httpmock.NewStringResponder(http.StatusOK, `{"total_results":3}`)
			httpmock.RegisterResponder(http.MethodGet, url, responder)

			resp, respType, err := GetOfficers(companyNumber)
			So(resp.TotalResults, ShouldEqual, 3)
			So(respType, ShouldEqual, Success)
			So(err, ShouldBeNil)
		})
	})

	Convey("Get Single Officer ", t, func() {

		officerID := "12345"
		url := "/emergency-auth-code/company/" + companyNumber + "/eligible-officers/" + officerID

		Convey("No response", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			resp, respType, err := GetOfficer(companyNumber, officerID)
			So(resp, ShouldBeNil)
			So(respType, ShouldEqual, Error)
			So(err, ShouldBeError, "Get "+url+": no responder found")
		})

		Convey("Empty response", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			responder := httpmock.NewStringResponder(http.StatusNotFound, "")
			httpmock.RegisterResponder(http.MethodGet, url, responder)

			resp, respType, err := GetOfficer(companyNumber, officerID)
			So(resp, ShouldBeNil)
			So(respType, ShouldEqual, NotFound)
			So(err, ShouldBeNil)
		})

		Convey("Successful response", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			responder := httpmock.NewStringResponder(http.StatusOK, `{"occupation":"bricklayer"}`)
			httpmock.RegisterResponder(http.MethodGet, url, responder)

			resp, respType, err := GetOfficer(companyNumber, officerID)
			So(resp, ShouldNotBeNil)
			So(resp.Occupation, ShouldEqual, "bricklayer")
			So(respType, ShouldEqual, Success)
			So(err, ShouldBeNil)
		})
	})
}
