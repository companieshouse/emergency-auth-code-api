package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/emergency-auth-code-api/models"
	"github.com/companieshouse/emergency-auth-code-api/service"
	"github.com/companieshouse/emergency-auth-code-api/utils"
)

// CreateAuthCodeRequest creates the auth code request for a specific officer ID
func CreateAuthCodeRequest(svc *service.AuthCodeService) http.Handler {
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

		if request.CompanyNumber == "" {
			errorMessage := "company number missing from request"
			log.ErrorR(req, fmt.Errorf(errorMessage))
			m := models.NewMessageResponse(errorMessage)
			utils.WriteJSONWithStatus(w, req, m, http.StatusBadRequest)
			return
		}

		companyHasAuthCode, err := svc.CheckAuthCodeExists(request.CompanyNumber)
		if err != nil {
			log.ErrorR(req, fmt.Errorf("error retrieving Auth Code from DB: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Info(fmt.Sprintf("companyHasAuthCode: [%v]", companyHasAuthCode))

		// TODO :- Add logic to create Authorization Code Request

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	})
}
