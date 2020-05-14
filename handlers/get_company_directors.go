package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/companieshouse/emergency-auth-code-api/service"

	"github.com/companieshouse/emergency-auth-code-api/utils"

	"github.com/companieshouse/chs.go/log"
)

// GetCompanyDirectorsHandler returns a list of valid company directors that can apply for an auth code
func GetCompanyDirectorsHandler(directorSvc service.DirectorDatabase) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		companyNumber := vars["company_number"]
		if companyNumber == "" {
			log.ErrorR(req, fmt.Errorf("no company number provided in request"))
			m := utils.NewMessageResponse(fmt.Sprintf("no company number provided in request"))
			utils.WriteJSONWithStatus(w, req, m, http.StatusBadRequest)
			return
		}

		// return list of directors from DirectorDatabase interface
		companyOfficers, err := directorSvc.GetCompanyDirectors(companyNumber)
		if err != nil {
			log.ErrorR(req, fmt.Errorf("error receiving data from oracle database: %v", err))
			m := utils.NewMessageResponse(fmt.Sprintf("error receiving data from director database: %v", err))
			utils.WriteJSONWithStatus(w, req, m, http.StatusInternalServerError)
			return
		}

		// return 404 if no directors found
		if companyOfficers.TotalCount == 0 {
			m := utils.NewMessageResponse(fmt.Sprintf("no directors found for company: %s", companyNumber))
			utils.WriteJSONWithStatus(w, req, m, http.StatusNotFound)
			return
		}

		// prepare response to request
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(&companyOfficers); err != nil {
			log.ErrorR(req, fmt.Errorf("error writing response: %v", err))
			m := utils.NewMessageResponse("error writing response")
			utils.WriteJSONWithStatus(w, req, m, http.StatusInternalServerError)
			return
		}
	})
}
