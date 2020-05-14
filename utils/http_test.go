package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitWriteJSONWithStatus(t *testing.T) {
	Convey("Failure to marshal json", t, func() {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		// causes an UnsupportedTypeError
		WriteJSONWithStatus(w, r, make(chan int), 500)

		So(w.Code, ShouldEqual, http.StatusInternalServerError)
		So(w.Header().Get("Content-Type"), ShouldEqual, "application/json")
		So(w.Body.String(), ShouldEqual, "")
	})

	Convey("contents are written as json", t, func() {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		m := NewMessageResponse("successful marshalling")

		WriteJSONWithStatus(w, r, m, http.StatusOK)

		So(w.Code, ShouldEqual, http.StatusOK)
		So(w.Header().Get("Content-Type"), ShouldEqual, "application/json")
	})
}
