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
	companyNumber, err := utils.GetValueFromVars(vars, "company_number")
	if err != nil {
		log.ErrorR(req, err)
		m := models.NewMessageResponse("company number is not in request context")
		utils.WriteJSONWithStatus(w, req, m, http.StatusBadRequest)
		return
	}

	startIndex := req.FormValue("start_index")
	itemsPerPage := req.FormValue("items_per_page")

	companyNumber = strings.ToUpper(companyNumber)

	companyOfficers, responseType, err := service.GetOfficers(companyNumber, startIndex, itemsPerPage)
	if err != nil {
		log.ErrorR(req, fmt.Errorf("error calling Oracle API to get officers: %v", err))
		m := models.NewMessageResponse("there was a problem communicating with the Oracle API")
		utils.WriteJSONWithStatus(w, req, m, http.StatusInternalServerError)
		return
	}

	if responseType == service.NotFound {
		m := models.NewMessageResponse("No officers found")
		utils.WriteJSONWithStatus(w, req, m, http.StatusNotFound)
		return
	}

	utils.WriteJSON(w, req, companyOfficers)
}

// GetCompanyOfficer returns a single company officer who may apply for an auth code
func GetCompanyOfficer(w http.ResponseWriter, req *http.Request) {

	// Check for a company number in request
	vars := mux.Vars(req)

	companyNumber, err := utils.GetValueFromVars(vars, "company_number")
	if err != nil {
		log.ErrorR(req, err)
		m := models.NewMessageResponse("company number not in request context")
		utils.WriteJSONWithStatus(w, req, m, http.StatusBadRequest)
		return
	}
	companyNumber = strings.ToUpper(companyNumber)

	// Check for Officer ID in request
	officerID, err := utils.GetValueFromVars(vars, "officer_id")
	if err != nil {
		log.ErrorR(req, err)
		m := models.NewMessageResponse("officer ID not in request context")
		utils.WriteJSONWithStatus(w, req, m, http.StatusBadRequest)
		return
	}

	companyOfficer, responseType, err := service.GetOfficer(companyNumber, officerID)
	if err != nil {
		log.ErrorR(req, fmt.Errorf("error calling Oracle API to get officer: %v", err))
		m := models.NewMessageResponse("there was a problem communicating with the Oracle API")
		utils.WriteJSONWithStatus(w, req, m, http.StatusInternalServerError)
		return
	}

	if responseType == service.NotFound {
		m := models.NewMessageResponse("No officer found")
		utils.WriteJSONWithStatus(w, req, m, http.StatusNotFound)
		return
	}

	utils.WriteJSON(w, req, companyOfficer)
}
