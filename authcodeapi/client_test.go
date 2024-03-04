package authcodeapi

import (
	"net/http"
	"testing"

	"github.com/companieshouse/emergency-auth-code-api/models"
	"github.com/jarcoal/httpmock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitSendAuthCodeItem(t *testing.T) {
	url := "api-url"
	path := "/api/test/authcode"
	queueAPIURL := url + path
	AuthCodeItem := models.AuthCodeItem{}
	testRequestID := "xyz123"
	authKey := "test-key-123"

	Convey("unexpected status returned from authcode API", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		client := NewClient(url, path, authKey)
		responder := httpmock.NewStringResponder(http.StatusNotFound, "")
		httpmock.RegisterResponder(http.MethodPost, queueAPIURL, responder)

		err := client.SendAuthCodeItem(&AuthCodeItem, testRequestID)
		So(err.Error(), ShouldEqual, "unexpected status returned from authCode API: 404")
	})

	Convey("queue API - success", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		client := NewClient(url, path, authKey)
		responder := httpmock.NewStringResponder(http.StatusOK, "error")
		httpmock.RegisterResponder(http.MethodPost, queueAPIURL, responder)

		err := client.SendAuthCodeItem(&AuthCodeItem, testRequestID)
		So(err, ShouldBeNil)
	})
}
