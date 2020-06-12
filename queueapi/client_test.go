package queueapi

import (
	"net/http"
	"testing"

	"github.com/companieshouse/emergency-auth-code-api/models"
	"github.com/jarcoal/httpmock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitSendQueueItem(t *testing.T) {
	url := "api-url"
	queueAPIURL := url + "/api/queue/authcode"
	queueItem := models.QueueItem{}

	Convey("unexpected status returned from queue API", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		client := NewClient(url)
		responder := httpmock.NewStringResponder(http.StatusNotFound, "")
		httpmock.RegisterResponder(http.MethodPost, queueAPIURL, responder)

		err := client.SendQueueItem(&queueItem)
		So(err.Error(), ShouldEqual, "unexpected status returned from queue API: 404")
	})

	Convey("queue API - success", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		client := NewClient(url)
		responder := httpmock.NewStringResponder(http.StatusOK, "error")
		httpmock.RegisterResponder(http.MethodPost, queueAPIURL, responder)

		err := client.SendQueueItem(&queueItem)
		So(err, ShouldBeNil)
	})
}
