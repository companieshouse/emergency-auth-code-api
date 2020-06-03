package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/companieshouse/emergency-auth-code-api/models"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitWriteJSON(t *testing.T) {
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
		m := models.NewMessageResponse("successful marshalling")

		WriteJSON(w, r, m)

		So(w.Code, ShouldEqual, http.StatusOK)
		So(w.Header().Get("Content-Type"), ShouldEqual, "application/json")
	})
}

func TestUnitGetCompanyNumber(t *testing.T) {
	Convey("Get Company Number", t, func() {
		vars := map[string]string{
			"company_number": "12345",
		}
		companyNumber, err := GetValueFromVars(vars, "company_number")
		So(companyNumber, ShouldEqual, "12345")
		So(err, ShouldBeNil)
	})

	Convey("No Company Number", t, func() {
		vars := map[string]string{}
		companyNumber, err := GetValueFromVars(vars, "company_number")
		So(companyNumber, ShouldBeEmpty)
		So(err.Error(), ShouldEqual, "company_number not found in vars")
	})
}
