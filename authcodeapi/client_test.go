package authcodeapi

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/companieshouse/emergency-auth-code-api/models"
	"github.com/jarcoal/httpmock"
	. "github.com/smartystreets/goconvey/convey"
)

const url = "api-url"
const path = "/api/test/authcode"
const queueAPIURL = url + path
const testRequestID = "xyz123"
const authKey = "test-key-123"

func TestUnitSendAuthCodeItem(t *testing.T) {
	AuthCodeItem := models.AuthCodeItem{}

	Convey("unexpected status returned from authcode API", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		client := NewClient(url, path, authKey)
		responder := httpmock.NewStringResponder(http.StatusNotFound, "")
		httpmock.RegisterResponder(http.MethodPost, queueAPIURL, responder)

		err := client.SendAuthCodeItem(&AuthCodeItem, testRequestID)
		So(err.Error(), ShouldEqual, "unexpected status returned from authCode API: 404")
	})

	Convey("queue API - success (OK - 200)", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		client := NewClient(url, path, authKey)
		responder := httpmock.NewStringResponder(http.StatusOK, "error")
		httpmock.RegisterResponder(http.MethodPost, queueAPIURL, responder)

		err := client.SendAuthCodeItem(&AuthCodeItem, testRequestID)
		So(err, ShouldBeNil)
	})

	Convey("queue API - success (CREATED - 201)", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		client := NewClient(url, path, authKey)
		responder := httpmock.NewStringResponder(http.StatusCreated, "error")
		httpmock.RegisterResponder(http.MethodPost, queueAPIURL, responder)

		err := client.SendAuthCodeItem(&AuthCodeItem, testRequestID)
		So(err, ShouldBeNil)
	})
}

func TestUnitSendRequest(t *testing.T) {

	Convey("send request - transport error", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		AuthCodeItem := models.AuthCodeItem{}
		client := NewClient(url, path, authKey)
		responder := httpmock.NewErrorResponder(fmt.Errorf("test-error"))
		httpmock.RegisterResponder(http.MethodPost, queueAPIURL, responder)

		response, err := client.sendRequest("POST", testRequestID, &AuthCodeItem)
		So(response, ShouldBeNil)
		So(err.Error(), ShouldContainSubstring, "test-error")
	})
}
