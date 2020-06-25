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

	mainRouter.HandleFunc("/emergency-auth-code-service/healthcheck", healthCheck).Methods(http.MethodGet).Name("healthcheck")

	// Create a router that requires all users to be authenticated when making requests
	appRouter := mainRouter.PathPrefix("/emergency-auth-code-service").Subrouter()
	appRouter.Use(userAuthInterceptor.UserAuthenticationIntercept)

	// Declare endpoint URIs
	appRouter.HandleFunc("/company/{company_number}/officers", GetCompanyOfficers).Queries("start_index", "{start_index:([0-9]+)?}", "items_per_page", "{items_per_page:([0-9]+)?}").Methods(http.MethodGet).Name("get-company-officers")
	appRouter.HandleFunc("/company/{company_number}/officers/{officer_id}", GetCompanyOfficer).Methods(http.MethodGet).Name("get-company-officer")
	appRouter.Handle("/auth-code-requests", CreateAuthCodeRequest(authCodeRequestService)).Methods(http.MethodPost).Name("create-auth-code-request")
	appRouter.Handle("/auth-code-requests/{auth_code_request_id}", GetAuthCodeRequest(authCodeRequestService)).Methods(http.MethodGet).Name("get-auth-code-request")
	appRouter.Handle("/auth-code-requests/{auth_code_request_id}", UpdateAuthCodeRequest(authCodeService, authCodeRequestService)).Methods(http.MethodPut).Name("update-auth-code-request")

	mainRouter.Use(log.Handler)
}

func healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
