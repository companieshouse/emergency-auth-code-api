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
var authCodeRequestService *service.AuthCodeRequestService

// Register defines the endpoints for the API
func Register(mainRouter *mux.Router, cfg *config.Config, authCodeDao dao.AuthcodeDAOService, authCodeRequestDao dao.AuthcodeRequestDAOService) {

	authCodeService = &service.AuthCodeService{
		Config: cfg,
		DAO:    authCodeDao,
	}

	authCodeRequestService = &service.AuthCodeRequestService{
		Config: cfg,
		DAO:    authCodeRequestDao,
	}

	userAuthInterceptor := &authentication.UserAuthenticationInterceptor{
		AllowAPIKeyUser:                true,
		RequireElevatedAPIKeyPrivilege: false,
	}

	// Create a router that requires all users to be authenticated when making requests
	appRouter := mainRouter.PathPrefix("/emergency-auth-code-service").Subrouter()
	appRouter.Use(userAuthInterceptor.UserAuthenticationIntercept)

	// Declare endpoint URIs
	appRouter.HandleFunc("/company/{company_number}/officers", GetCompanyOfficers).Methods(http.MethodGet).Name("get-company-officers")
	appRouter.HandleFunc("/company/{company_number}/officers/{officer_id}", GetCompanyOfficer).Methods(http.MethodGet).Name("get-company-officer")
	appRouter.Handle("/auth-code-requests", CreateAuthCodeRequest(authCodeService, authCodeRequestService)).Methods(http.MethodPost).Name("create-auth-code-request")
	appRouter.Handle("/auth-code-requests/{auth_code_request_id}", GetAuthCodeRequest(authCodeRequestService)).Methods(http.MethodGet).Name("get-auth-code-request")

	mainRouter.Use(log.Handler)
}
