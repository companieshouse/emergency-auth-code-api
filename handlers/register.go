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
var officerService *service.OfficersService

// Register defines the endpoints for the API
func Register(mainRouter *mux.Router, cfg *config.Config, authCodeDao dao.AuthcodeDAOService, officerDao dao.OfficerDAOService) {

	authCodeService = &service.AuthCodeService{
		Config: cfg,
		DAO:    authCodeDao,
	}

	officerService = &service.OfficersService{
		Config: cfg,
		DAO:    officerDao,
	}

	userAuthInterceptor := &authentication.UserAuthenticationInterceptor{
		AllowAPIKeyUser:                true,
		RequireElevatedAPIKeyPrivilege: false,
	}

	// Create a router that requires all users to be authenticated when making requests
	appRouter := mainRouter.PathPrefix("/emergency-auth-code-service").Subrouter()
	appRouter.Use(userAuthInterceptor.UserAuthenticationIntercept)

	// Declare endpoint URIs
	appRouter.Handle("/company/{company_number}/officers", GetCompanyOfficersHandler(officerService)).Methods(http.MethodGet).Name("get-company-officers")
	appRouter.Handle("/auth-code-requests", CreateAuthCodeRequest(authCodeService)).Methods(http.MethodPost).Name("create-auth-code-request")

	mainRouter.Use(log.Handler)
}
