package oracle

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetOfficers(t *testing.T) {
	companyNumber := "87654321"
	url := "api-url/emergency-auth-code/company/87654321/eligible-officers"

	Convey("Officers not found", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		client := NewClient("api-url")
		responder := httpmock.NewStringResponder(http.StatusNotFound, "")
		httpmock.RegisterResponder(http.MethodGet, url, responder)

		resp, err := client.GetOfficers(companyNumber)
		So(resp, ShouldBeNil)
		So(err, ShouldBeNil)
	})

	Convey("Failure to read response", t, func() {
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

	Convey("Error response - bad request", t, func() {
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

	Convey("Error response - internal server error", t, func() {
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

	Convey("Error response - unexpected error", t, func() {
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

	Convey("Bad response", t, func() {
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

	Convey("Successful response", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		client := NewClient("api-url")
		responder := httpmock.NewStringResponder(http.StatusOK, `{"total_results":3}`)
		httpmock.RegisterResponder(http.MethodGet, url, responder)

		resp, err := client.GetOfficers(companyNumber)
		So(err, ShouldBeNil)
		So(resp.TotalResults, ShouldEqual, 3)
	})
}
