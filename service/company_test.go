package service

import (
	"context"
	"net/http"
	"testing"

	"github.com/companieshouse/go-session-handler/httpsession"
	"github.com/companieshouse/go-session-handler/session"
	"github.com/jarcoal/httpmock"

	. "github.com/smartystreets/goconvey/convey"
)

var companyDetailsResponse = `
{
  "company_name": "Test Company",
  "registered_office_address" : {
    "postal_code" : "CF14 3UZ",
    "address_line_2" : "Cardiff",
    "address_line_1" : "100 Crown Way"
  }
}
`

func TestUnitGetCompanyInformation(t *testing.T) {

	Convey("GetCompanyNameFromCompanyProfileAPI", t, func() {

		apiURL := "https://api.companieshouse.gov.uk"

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		ctx := context.WithValue(context.Background(), httpsession.ContextKeySession, &session.Session{})
		r := &http.Request{}
		r = r.WithContext(ctx)

		Convey("invalid request", func() {
			defer httpmock.Reset()
			httpmock.RegisterResponder(http.MethodGet, apiURL+"/company/12345678", httpmock.NewStringResponder(http.StatusTeapot, ""))

			resp, err := GetCompanyName("12345678", &http.Request{})
			So(resp, ShouldBeEmpty)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, `ch-api: got HTTP response code 418 with body: `)
		})

		Convey("it returns a serialised version of the response", func() {
			defer httpmock.Reset()
			httpmock.RegisterResponder(http.MethodGet, apiURL+"/company/12345678", httpmock.NewStringResponder(http.StatusOK, companyDetailsResponse))

			resp, err := GetCompanyName("12345678", r)

			So(err, ShouldBeNil)
			So(resp, ShouldEqual, "Test Company")
		})
	})
}
