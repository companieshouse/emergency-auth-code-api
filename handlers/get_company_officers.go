package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/emergency-auth-code-api/models"
	"github.com/companieshouse/emergency-auth-code-api/service"
	"github.com/companieshouse/emergency-auth-code-api/utils"
	"github.com/gorilla/mux"
)

// GetCompanyOfficers returns a list of valid company officers who may apply for an auth code
func GetCompanyOfficers(w http.ResponseWriter, req *http.Request) {

	// Check for a company number in request
	vars := mux.Vars(req)
	companyNumber, err := utils.GetCompanyNumberFromVars(vars)
	if err != nil {
		log.ErrorR(req, err)
		m := models.NewMessageResponse("company number is not in request context")
		utils.WriteJSONWithStatus(w, req, m, http.StatusBadRequest)
		return
	}

	companyNumber = strings.ToUpper(companyNumber)

	companyOfficers, responseType, err := service.GetOfficers(companyNumber)
	if err != nil {
		log.ErrorR(req, fmt.Errorf("error calling Oracle API to get officers: %v", err))
		switch responseType {
		case service.InvalidData:
			m := models.NewMessageResponse("failed to read officers data")
			utils.WriteJSONWithStatus(w, req, m, http.StatusBadRequest)
			return
		case service.Error:
		default:
			m := models.NewMessageResponse("there was a problem communicating with the Oracle API")
			utils.WriteJSONWithStatus(w, req, m, http.StatusInternalServerError)
			return
		}
	}

	if responseType == service.NotFound {
		m := models.NewMessageResponse("No officers found")
		utils.WriteJSONWithStatus(w, req, m, http.StatusNotFound)
		return
	}

	utils.WriteJSON(w, req, companyOfficers)
}
