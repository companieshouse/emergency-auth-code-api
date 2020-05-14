package handlers

import (
	"net/http"

	"github.com/companieshouse/chs.go/authentication"
	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/emergency-auth-code-api/config"
	"github.com/companieshouse/emergency-auth-code-api/dao"
	"github.com/companieshouse/emergency-auth-code-api/service"
	"github.com/gorilla/mux"
)

var authCodeService *service.AuthCodeService

// Register defines the endpoints for the API
func Register(mainRouter *mux.Router, cfg *config.Config, svc dao.Service) {

	authCodeService = &service.AuthCodeService{
		Config: cfg,
		DAO:    svc,
	}

	userAuthInterceptor := &authentication.UserAuthenticationInterceptor{
		AllowAPIKeyUser:                true,
		RequireElevatedAPIKeyPrivilege: false,
	}

	// Create a router that requires all users to be authenticated when making requests
	appRouter := mainRouter.PathPrefix("/emergency-auth-code-service").Subrouter()
	appRouter.Use(userAuthInterceptor.UserAuthenticationIntercept)

	// Declare endpoint URIs
	appRouter.HandleFunc("/company/{company_number}/officers", GetCompanyDirectors).Methods(http.MethodGet).Name("get-company-directors")
	appRouter.Handle("/auth-code-requests", CreateAuthCodeRequest(authCodeService)).Methods(http.MethodPost).Name("create-auth-code-request")

	mainRouter.Use(log.Handler)
}