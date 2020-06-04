package handlers

import (
	"net/http"

	"github.com/companieshouse/emergency-auth-code-api/service"
)

// UpdateAuthCodeRequest updates an auth code request for a specified auth-code-request ID
func UpdateAuthCodeRequest(authCodeReqSvc *service.AuthCodeRequestService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		// TODO

		return

	})
}
