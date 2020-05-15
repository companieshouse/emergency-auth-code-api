package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/emergency-auth-code-api/models"
	"github.com/companieshouse/emergency-auth-code-api/service"
	"github.com/companieshouse/emergency-auth-code-api/utils"
	"github.com/gorilla/mux"
)

// GetCompanyOfficersHandler returns a list of valid company officers that can apply for an auth code
func GetCompanyOfficersHandler(svc *service.OfficersService) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		companyNumber := vars["company_number"]
		if companyNumber == "" {
			log.ErrorR(req, fmt.Errorf("no company number provided in request"))
			m := models.NewMessageResponse(fmt.Sprintf("no company number provided in request"))
			utils.WriteJSONWithStatus(w, req, m, http.StatusBadRequest)
			return
		}

		// return list of officers from officers service interface
		companyOfficers, err := svc.GetListOfCompanyOfficers(companyNumber)
		if err != nil {
			log.ErrorR(req, fmt.Errorf("error receiving data from officer database: %v", err))
			m := models.NewMessageResponse(fmt.Sprintf("error receiving data from officer database: %v", err))
			utils.WriteJSONWithStatus(w, req, m, http.StatusInternalServerError)
			return
		}

		// return 404 if no officers found
		if companyOfficers.TotalCount == 0 {
			m := models.NewMessageResponse(fmt.Sprintf("no officer found for company: %s", companyNumber))
			utils.WriteJSONWithStatus(w, req, m, http.StatusNotFound)
			return
		}

		// prepare response to request
		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(&companyOfficers); err != nil {
			log.ErrorR(req, fmt.Errorf("error writing response: %v", err))
			m := models.NewMessageResponse("error writing response")
			utils.WriteJSONWithStatus(w, req, m, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
