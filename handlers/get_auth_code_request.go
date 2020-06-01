package handlers

import (
	"fmt"
	"net/http"

	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/emergency-auth-code-api/models"
	"github.com/companieshouse/emergency-auth-code-api/service"
	"github.com/companieshouse/emergency-auth-code-api/utils"
	"github.com/gorilla/mux"
)

// GetAuthCodeRequest returns an auth code request for a specified auth-code-request ID
func GetAuthCodeRequest(authCodeReqSvc *service.AuthCodeRequestService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Check for a auth-code-request ID in the request
		vars := mux.Vars(req)
		authCodeRequestId := vars["auth_code_request_id"]
		if authCodeRequestId == "" {
			log.ErrorR(req, fmt.Errorf("no auth code request id in request"))
			m := models.NewMessageResponse("no auth code request id in request")
			utils.WriteJSONWithStatus(w, req, m, http.StatusBadRequest)
			return
		}

		// Get the auth code request from the ID in request
		authCodeRequest, responseType := authCodeReqSvc.GetAuthCodeRequest(authCodeRequestId)
		if responseType != http.StatusOK {
			w.WriteHeader(responseType)
			return
		}

		utils.WriteJSONWithStatus(w, req, authCodeRequest, responseType)
	})
}
