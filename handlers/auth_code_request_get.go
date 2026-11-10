package handlers

import (
	"fmt"
	"net/http"

	"github.com/companieshouse/chs.go/authentication"
	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/emergency-auth-code-api/models"
	"github.com/companieshouse/emergency-auth-code-api/service"
	"github.com/companieshouse/emergency-auth-code-api/utils"
	"github.com/gorilla/mux"
)

// GetAuthCodeRequest returns an auth code request for a specified auth-code-request ID
func GetAuthCodeRequest(authCodeReqSvc *service.AuthCodeRequestService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Check for an auth-code-request ID in the request
		vars := mux.Vars(req)
		authCodeRequestID := vars["auth_code_request_id"]
		if authCodeRequestID == "" {
			log.ErrorR(req, fmt.Errorf("no auth code request id in request"))
			m := models.NewMessageResponse("no auth code request id in request")
			utils.WriteJSONWithStatus(w, req, m, http.StatusBadRequest)
			return
		}
		userDetails := req.Context().Value(authentication.ContextKeyUserDetails)
		if userDetails == nil {
			log.ErrorR(req, fmt.Errorf("user details not in context"))
			m := models.NewMessageResponse("user details not in request context")
			utils.WriteJSONWithStatus(w, req, m, http.StatusUnauthorized)
			return
		}
		createdBy := userDetails.(authentication.AuthUserDetails)

		// Get the auth code request from the ID in request
		authCodeRequest, responseType := authCodeReqSvc.GetAuthCodeRequest(authCodeRequestID, createdBy)
		if responseType != http.StatusOK {
			w.WriteHeader(responseType)
			return
		}

		utils.WriteJSONWithStatus(w, req, authCodeRequest, responseType)
	})
}
