package handlers

import (
	"net/http"

	"github.com/companieshouse/emergency-auth-code-api/service"

	"github.com/companieshouse/chs.go/authentication"
	"github.com/companieshouse/chs.go/log"
	"github.com/gorilla/mux"
)

// Register defines the endpoints for the API
func Register(mainRouter *mux.Router, directorSvc service.DirectorDatabase) {

	userAuthInterceptor := &authentication.UserAuthenticationInterceptor{
		AllowAPIKeyUser:                true,
		RequireElevatedAPIKeyPrivilege: false,
	}

	// Create a router that requires all users to be authenticated when making requests
	appRouter := mainRouter.PathPrefix("/emergency-auth-code-service").Subrouter()
	appRouter.Use(userAuthInterceptor.UserAuthenticationIntercept)

	// Declare endpoint URIs
	appRouter.Handle("/company/{company_number}/officers", GetCompanyDirectorsHandler(directorSvc)).Methods(http.MethodGet).Name("get-company-directors")
	appRouter.HandleFunc("/auth-code-requests", CreateAuthCodeRequest).Methods(http.MethodPost).Name("create-auth-code-request")

	mainRouter.Use(log.Handler)
}
