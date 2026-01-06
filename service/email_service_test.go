package service

import (
	"fmt"
	"github.com/companieshouse/emergency-auth-code-api/config"
	"github.com/jarcoal/httpmock"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestUnitSendEmail(t *testing.T) {
	// Build test config
	cfg, _ := config.Get()
	cfg.CHSURL = "http://local.test"
	cfg.ChsKafkaApiURL = "http://local.test.chs.kafka"
	cfg.APIKey = "testApiKey"

	Convey("error sending email", t, func() {
		res := SendEmail("test@test.com")

		So(res.Error(), ShouldContainSubstring, "error sending email")
	})

	Convey("wrong status code from kafka api", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		responder := httpmock.NewStringResponder(http.StatusInternalServerError, "")
		httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s/send-email", cfg.ChsKafkaApiURL), responder)

		res := SendEmail("test@test.com")

		// Assert send-email endpoint was hit
		timesHttpHit := httpmock.GetCallCountInfo()
		timesEmailEndpointHit := timesHttpHit[fmt.Sprintf("POST %s/send-email", cfg.ChsKafkaApiURL)]
		So(timesEmailEndpointHit, ShouldEqual, 1)

		// Assert error from sending email
		So(res.Error(), ShouldContainSubstring, "wrong status code from kafka api")
	})

	Convey("successfully send email", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		responder := httpmock.NewStringResponder(http.StatusOK, ``)
		httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s/send-email", cfg.ChsKafkaApiURL), responder)

		res := SendEmail("test@test.com")

		// Assert send-email endpoint was hit
		timesHttpHit := httpmock.GetCallCountInfo()
		timesEmailEndpointHit := timesHttpHit[fmt.Sprintf("POST %s/send-email", cfg.ChsKafkaApiURL)]
		So(timesEmailEndpointHit, ShouldEqual, 1)

		// Assert no errors from sending email
		So(res, ShouldBeNil)
	})

}
