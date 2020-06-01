package config

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetConfig(t *testing.T) {
	Convey("Successful Get Config", t, func() {
		os.Setenv("BIND_ADDR", "123")
		config, err := Get()
		So(err, ShouldBeNil)
		So(config, ShouldNotBeNil)
		So(config.BindAddr, ShouldEqual, "123")
	})

	Convey("Config already defined", t, func() {
		os.Setenv("BIND_ADDR", "123")
		_, _ = Get()
		os.Setenv("BIND_ADDR", "456")
		config, err := Get()
		So(err, ShouldBeNil)
		So(config, ShouldNotBeNil)
		So(config.BindAddr, ShouldEqual, "123")
	})

}
