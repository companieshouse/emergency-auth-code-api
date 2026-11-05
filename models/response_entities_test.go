package models

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitNewMessageResponse(t *testing.T) {
	Convey("new message response", t, func() {
		responseResource := NewMessageResponse("test-message")
		So(responseResource.Message, ShouldEqual, "test-message")
	})
}
