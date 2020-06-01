package service

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetOfficers(t *testing.T) {
	companyNumber := "87654321"
	url := "/emergency-auth-code/company/87654321/eligible-officers"

	Convey("No response", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		resp, respType, err := GetOfficers(companyNumber)
		So(resp, ShouldBeNil)
		So(respType, ShouldEqual, Error)
		So(err, ShouldBeError, `Get "/emergency-auth-code/company/87654321/eligible-officers": no responder found`)
	})

	Convey("Empty response", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		responder := httpmock.NewStringResponder(http.StatusNotFound, "")
		httpmock.RegisterResponder(http.MethodGet, url, responder)

		resp, respType, err := GetOfficers(companyNumber)
		So(resp, ShouldBeNil)
		So(respType, ShouldEqual, NotFound)
		So(err, ShouldBeNil)
	})

	Convey("Successful response", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		responder := httpmock.NewStringResponder(http.StatusOK, `{"total_results":3}`)
		httpmock.RegisterResponder(http.MethodGet, url, responder)

		resp, respType, err := GetOfficers(companyNumber)
		So(resp.TotalResults, ShouldEqual, 3)
		So(respType, ShouldEqual, Success)
		So(err, ShouldBeNil)
	})
}
