package handlers

import (
	"net/http"

	"github.com/companieshouse/chs.go/authentication"
	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/emergency-auth-code-api/config"
	"github.com/gorilla/mux"
)

// Contains config and DAO for
type EmergencyAuthCodeService struct {
	Config *config.Config
}

func Register(mainRouter *mux.Router) {

	userAuthInterceptor := &authentication.UserAuthenticationInterceptor{
		AllowAPIKeyUser:                true,
		RequireElevatedAPIKeyPrivilege: false,
	}

	// Create a router that requires all users to be authenticated when making requests
	appRouter := mainRouter.PathPrefix("/emergency-auth-code-service").Subrouter()
	appRouter.Use(userAuthInterceptor.UserAuthenticationIntercept)

	// Declare endpoint URIs
	appRouter.HandleFunc("/company/{company_number}/officers", GetCompanyDirectors).Methods(http.MethodGet).Name("get-company-directors")
	appRouter.HandleFunc("/auth-code-requests", CreateAuthCodeRequest).Methods(http.MethodPost).Name("create-auth-code-request")

	mainRouter.Use(log.Handler)
}
