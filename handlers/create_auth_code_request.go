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
func CreateAuthCodeRequest(authCodeSvc *service.AuthCodeService, authCodeReqSvc *service.AuthCodeRequestService) http.Handler {
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

		model := transformers.AuthCodeResourceRequestToDB(&request)

		companyHasAuthCode, err := authCodeSvc.CheckAuthCodeExists(request.CompanyNumber) // TODO move this to the PUT/update when implemented
		if err != nil {
			log.ErrorR(req, fmt.Errorf("error retrieving Auth Code from DB: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var letterType string
		if companyHasAuthCode {
			letterType = "reminder"
		} else {
			letterType = "apply"
		}
		model.Data.Type = letterType

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
