package oracle

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitGetOfficers(t *testing.T) {
	companyNumber := "87654321"

	Convey("Get Officer List", t, func() {

		url := "api-url/emergency-auth-code/company/" + companyNumber + "/eligible-officers"

		Convey("Officers not found", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			client := NewClient("api-url")
			responder := httpmock.NewStringResponder(http.StatusNotFound, "")
			httpmock.RegisterResponder(http.MethodGet, url, responder)

			resp, err := client.GetOfficers(companyNumber)
			So(resp, ShouldBeNil)
			So(err, ShouldBeNil)
		})

		Convey("Failure to read response", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			client := NewClient("api-url")
			responder := httpmock.NewStringResponder(http.StatusInternalServerError, "")
			httpmock.RegisterResponder(http.MethodGet, url, responder)

			resp, err := client.GetOfficers(companyNumber)
			So(resp, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err, ShouldBeError, ErrFailedToReadBody)
		})

		Convey("Error response - bad request", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			client := NewClient("api-url")
			responder := httpmock.NewStringResponder(http.StatusBadRequest, `{"httpStatusCode" : 500}`)
			httpmock.RegisterResponder(http.MethodGet, url, responder)

			resp, err := client.GetOfficers(companyNumber)
			So(resp, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err, ShouldBeError, ErrOracleAPIBadRequest)
		})

		Convey("Error response - internal server error", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			client := NewClient("api-url")
			responder := httpmock.NewStringResponder(http.StatusInternalServerError, `{"httpStatusCode" : 500}`)
			httpmock.RegisterResponder(http.MethodGet, url, responder)

			resp, err := client.GetOfficers(companyNumber)
			So(resp, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err, ShouldBeError, ErrOracleAPIInternalServer)
		})

		Convey("Error response - unexpected error", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			client := NewClient("api-url")
			responder := httpmock.NewStringResponder(http.StatusTeapot, `{"httpStatusCode" : 500}`)
			httpmock.RegisterResponder(http.MethodGet, url, responder)

			resp, err := client.GetOfficers(companyNumber)
			So(resp, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err, ShouldBeError, ErrUnexpectedServerError)
		})

		Convey("Bad response", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			client := NewClient("api-url")
			responder := httpmock.NewStringResponder(http.StatusOK, "")
			httpmock.RegisterResponder(http.MethodGet, url, responder)

			resp, err := client.GetOfficers(companyNumber)
			So(resp, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err, ShouldBeError, ErrFailedToReadBody)
		})

		Convey("Successful response", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			client := NewClient("api-url")
			responder := httpmock.NewStringResponder(http.StatusOK, `{"total_results":3}`)
			httpmock.RegisterResponder(http.MethodGet, url, responder)

			resp, err := client.GetOfficers(companyNumber)
			So(err, ShouldBeNil)
			So(resp.TotalResults, ShouldEqual, 3)
		})
	})

	Convey("Get Single Officer", t, func() {

		url := "api-url/emergency-auth-code/company/" + companyNumber + "/eligible-officers/" + "54321"
		officerID := "54321"

		Convey("Officer not found", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			client := NewClient("api-url")
			responder := httpmock.NewStringResponder(http.StatusNotFound, "")
			httpmock.RegisterResponder(http.MethodGet, url, responder)

			resp, err := client.GetOfficer(companyNumber, officerID)
			So(resp, ShouldBeNil)
			So(err, ShouldBeNil)
		})

		Convey("Failure to read response", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			client := NewClient("api-url")
			responder := httpmock.NewStringResponder(http.StatusInternalServerError, "")
			httpmock.RegisterResponder(http.MethodGet, url, responder)

			resp, err := client.GetOfficer(companyNumber, officerID)
			So(resp, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err, ShouldBeError, ErrFailedToReadBody)
		})

		Convey("Error response - bad request", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			client := NewClient("api-url")
			responder := httpmock.NewStringResponder(http.StatusBadRequest, `{"httpStatusCode" : 500}`)
			httpmock.RegisterResponder(http.MethodGet, url, responder)

			resp, err := client.GetOfficer(companyNumber, officerID)
			So(resp, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err, ShouldBeError, ErrOracleAPIBadRequest)
		})

		Convey("Error response - internal server error", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			client := NewClient("api-url")
			responder := httpmock.NewStringResponder(http.StatusInternalServerError, `{"httpStatusCode" : 500}`)
			httpmock.RegisterResponder(http.MethodGet, url, responder)

			resp, err := client.GetOfficer(companyNumber, officerID)
			So(resp, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err, ShouldBeError, ErrOracleAPIInternalServer)
		})

		Convey("Error response - unexpected error", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			client := NewClient("api-url")
			responder := httpmock.NewStringResponder(http.StatusTeapot, `{"httpStatusCode" : 500}`)
			httpmock.RegisterResponder(http.MethodGet, url, responder)

			resp, err := client.GetOfficer(companyNumber, officerID)
			So(resp, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err, ShouldBeError, ErrUnexpectedServerError)
		})

		Convey("Bad response", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			client := NewClient("api-url")
			responder := httpmock.NewStringResponder(http.StatusOK, "")
			httpmock.RegisterResponder(http.MethodGet, url, responder)

			resp, err := client.GetOfficer(companyNumber, officerID)
			So(resp, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err, ShouldBeError, ErrFailedToReadBody)
		})

		Convey("Successful response", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			client := NewClient("api-url")
			responder := httpmock.NewStringResponder(http.StatusOK, `{"occupation":"bricklayer"}`)
			httpmock.RegisterResponder(http.MethodGet, url, responder)

			resp, err := client.GetOfficer(companyNumber, officerID)
			So(err, ShouldBeNil)
			So(resp.Occupation, ShouldEqual, "bricklayer")
		})
	})

	Convey("Check Filing History", t, func() {

		url := "api-url/emergency-auth-code/company/" + companyNumber + "/efiling-status"

		Convey("Failure to read response", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			client := NewClient("api-url")
			responder := httpmock.NewStringResponder(http.StatusInternalServerError, "")
			httpmock.RegisterResponder(http.MethodGet, url, responder)

			resp, err := client.CheckFilingHistory(companyNumber)
			So(resp, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err, ShouldBeError, ErrFailedToReadBody)
		})

		Convey("Bad response", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			client := NewClient("api-url")
			responder := httpmock.NewStringResponder(http.StatusOK, "")
			httpmock.RegisterResponder(http.MethodGet, url, responder)

			resp, err := client.CheckFilingHistory(companyNumber)
			So(resp, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err, ShouldBeError, ErrFailedToReadBody)
		})

		Convey("Successful response", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			client := NewClient("api-url")
			responder := httpmock.NewStringResponder(http.StatusOK, `{"efiling_found_in_period":false}`)
			httpmock.RegisterResponder(http.MethodGet, url, responder)

			resp, err := client.CheckFilingHistory(companyNumber)
			So(err, ShouldBeNil)
			So(resp.EFilingFoundInPeriod, ShouldBeFalse)
		})
	})
}
