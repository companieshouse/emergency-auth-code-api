package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/companieshouse/chs.go/authentication"
	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/emergency-auth-code-api/models"
	"github.com/companieshouse/emergency-auth-code-api/service"
	"github.com/companieshouse/emergency-auth-code-api/transformers"
	"github.com/companieshouse/emergency-auth-code-api/utils"
)

// CreateAuthCodeRequest creates the auth code request for a specific officer ID
func CreateAuthCodeRequest(authCodeReqSvc *service.AuthCodeRequestService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var request models.AuthCodeRequest
		err := json.NewDecoder(req.Body).Decode(&request)

		// request body failed to get decoded
		if err != nil {
			log.ErrorR(req, fmt.Errorf("invalid request"))
			m := models.NewMessageResponse("failed to read request body")
			utils.WriteJSONWithStatus(w, req, m, http.StatusBadRequest)
			return
		}

		userDetails := req.Context().Value(authentication.ContextKeyUserDetails)
		if userDetails == nil {
			log.ErrorR(req, fmt.Errorf("user details not in context"))
			m := models.NewMessageResponse("user details not in request context")
			utils.WriteJSONWithStatus(w, req, m, http.StatusBadRequest)
			return
		}

		if request.CompanyNumber == "" {
			errorMessage := "company number missing from request"
			log.ErrorR(req, fmt.Errorf(errorMessage))
			m := models.NewMessageResponse(errorMessage)
			utils.WriteJSONWithStatus(w, req, m, http.StatusBadRequest)
			return
		}

		request.CreatedBy = userDetails.(authentication.AuthUserDetails)

		if request.OfficerID != "" {
			// retrieve details for officer from oracle-query-api
			officer, officerResponse, err := service.GetOfficerDetails(request.CompanyNumber, request.OfficerID)
			if err != nil {
				log.ErrorR(req, fmt.Errorf("error calling Oracle API to get officer: %v", err))
				m := models.NewMessageResponse("there was a problem communicating with the Oracle API")
				utils.WriteJSONWithStatus(w, req, m, http.StatusInternalServerError)
				return
			}
			if officerResponse == service.NotFound {
				m := models.NewMessageResponse("No officer found")
				utils.WriteJSONWithStatus(w, req, m, http.StatusNotFound)
				return
			}

			request.OfficerUraID = officer.UsualResidentialAddress.ID
			request.OfficerForename = officer.Forename
			request.OfficerSurname = officer.Surname
		} else {
			// check if any eligible officers exist for specified company
			companyIsEligible, err := service.CheckOfficers(request.CompanyNumber)
			if err != nil {
				utils.WriteErrorMessage(w, req, http.StatusInternalServerError, "there was a problem communicating with the Oracle API")
				return
			}
			if !companyIsEligible {
				utils.WriteResponseMessage(w, req, http.StatusNotFound, "corporate body has no eligible officers")
				return
			}
		}

		hasFiledWithinPeriod, err := service.CheckCompanyFilingHistory(request.CompanyNumber)
		if err != nil {
			log.ErrorR(req, fmt.Errorf("error calling Oracle API to check filing history: %v", err))
			m := models.NewMessageResponse("there was a problem communicating with the Oracle API")
			utils.WriteJSONWithStatus(w, req, m, http.StatusInternalServerError)
			return
		}
		if hasFiledWithinPeriod {
			log.Info(fmt.Sprintf("company has had a filing within a recent period: %v", request.CompanyNumber))
			m := models.NewMessageResponse("the company has had a filing within a recent period")
			utils.WriteJSONWithStatus(w, req, m, http.StatusForbidden)
			return
		}

		model := transformers.AuthCodeResourceRequestToDB(&request)

		companyName, err := service.GetCompanyName(request.CompanyNumber, req)
		if err != nil {
			log.ErrorR(req, fmt.Errorf("error getting company name: [%v]", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		model.Data.CompanyName = companyName

		err = authCodeReqSvc.CreateAuthCodeRequest(model)
		if err != nil {
			log.ErrorR(req, fmt.Errorf("error creating Auth Code Request: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		utils.WriteJSONWithStatus(w, req, transformers.AuthCodeRequestResourceDaoToResponse(model), http.StatusCreated)
	})
}
