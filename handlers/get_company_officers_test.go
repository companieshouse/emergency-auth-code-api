package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jarcoal/httpmock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitGetCompanyOfficers(t *testing.T) {
	Convey("GetCompanyOfficer tests", t, func() {
		url := "/emergency-auth-code/company/87654321/eligible-officers"

		Convey("company number missing from request", func() {
			req, _ := http.NewRequest("GET", "url", nil)
			w := httptest.NewRecorder()

			req = mux.SetURLVars(req, map[string]string{"company_number": ""})

			GetCompanyOfficers(w, req)
			So(w.Code, ShouldEqual, http.StatusBadRequest)
		})

		Convey("response error", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			responder := httpmock.NewStringResponder(http.StatusInternalServerError, "")
			httpmock.RegisterResponder(http.MethodGet, url, responder)

			req, _ := http.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			req = mux.SetURLVars(req, map[string]string{"company_number": "87654321"})

			GetCompanyOfficers(w, req)
			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})

		Convey("no officers found", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			responder := httpmock.NewStringResponder(http.StatusNotFound, "")
			httpmock.RegisterResponder(http.MethodGet, url, responder)

			req, _ := http.NewRequest("GET", "url", nil)
			w := httptest.NewRecorder()

			req = mux.SetURLVars(req, map[string]string{"company_number": "87654321"})

			GetCompanyOfficers(w, req)
			So(w.Code, ShouldEqual, http.StatusNotFound)
		})

		Convey("Success - officers found", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			responder := httpmock.NewStringResponder(http.StatusOK, `{"total_results":3}`)
			httpmock.RegisterResponder(http.MethodGet, url, responder)

			req, _ := http.NewRequest("GET", "url", nil)
			w := httptest.NewRecorder()

			req = mux.SetURLVars(req, map[string]string{"company_number": "87654321"})

			GetCompanyOfficers(w, req)
			So(w.Code, ShouldEqual, http.StatusOK)
		})
	})
}
