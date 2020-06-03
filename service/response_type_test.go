package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitResponseType(t *testing.T) {
	Convey("Successful Get Response Type", t, func() {
		So(NotFound.String(), ShouldEqual, "not-found")
	})
}
