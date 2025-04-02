package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/companieshouse/chs.go/authentication"
	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/emergency-auth-code-api/models"
	"github.com/companieshouse/emergency-auth-code-api/service"
	"github.com/companieshouse/emergency-auth-code-api/transformers"
	"github.com/companieshouse/emergency-auth-code-api/utils"
)

// CreateAuthCodeRequest creates the auth code request for a specific officer ID
func CreateAuthCodeRequest(authCodeReqSvc *service.AuthCodeRequestService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var request models.AuthCodeRequest
		err := json.NewDecoder(req.Body).Decode(&request)

		fmt.Fprint(os.Stdout, "[debug] func CreateAuthCodeRequest(): 1...")

		// request body failed to get decoded
		if err != nil {
			log.ErrorR(req, fmt.Errorf("invalid request"))
			m := models.NewMessageResponse("failed to read request body")
			utils.WriteJSONWithStatus(w, req, m, http.StatusBadRequest)
			return
		}

		fmt.Fprint(os.Stdout, "[debug] func CreateAuthCodeRequest(): 2...")

		userDetails := req.Context().Value(authentication.ContextKeyUserDetails)
		if userDetails == nil {
			log.ErrorR(req, fmt.Errorf("user details not in context"))
			m := models.NewMessageResponse("user details not in request context")
			utils.WriteJSONWithStatus(w, req, m, http.StatusBadRequest)
			return
		}

		fmt.Fprint(os.Stdout, "[debug] func CreateAuthCodeRequest(): 3...")

		if request.CompanyNumber == "" {
			errorMessage := "company number missing from request"
			log.ErrorR(req, fmt.Errorf(errorMessage))
			m := models.NewMessageResponse(errorMessage)
			utils.WriteJSONWithStatus(w, req, m, http.StatusBadRequest)
			return
		}

		fmt.Fprint(os.Stdout, "[debug] func CreateAuthCodeRequest(): 4...")

		createdBy := userDetails.(authentication.AuthUserDetails)

		fmt.Fprint(os.Stdout, "[debug] func CreateAuthCodeRequest(): 5...")

		validCorporateBody, err := validateCorporateBody(req, authCodeReqSvc, request.CompanyNumber, createdBy.Email)

		fmt.Fprint(os.Stdout, "[debug] func CreateAuthCodeRequest(): 6...")

		if err != nil {
			utils.WriteErrorMessage(w, req, http.StatusInternalServerError, "error checking corporate body")
			return
		}

		fmt.Fprint(os.Stdout, "[debug] func CreateAuthCodeRequest(): 7...")

		if !validCorporateBody {
			utils.WriteResponseMessage(w, req, http.StatusForbidden, "request not permitted for corporate body")
			return
		}

		fmt.Fprint(os.Stdout, "[debug] func CreateAuthCodeRequest(): 8...")

		request.CreatedBy = createdBy

		fmt.Fprint(os.Stdout, "[debug] func CreateAuthCodeRequest(): 9...")

		if request.OfficerID != "" {
			fmt.Fprint(os.Stdout, "[debug] func CreateAuthCodeRequest(): 9.1.1...")
			// retrieve details for officer from oracle-query-api
			officer, officerResponse, err := service.GetOfficerDetails(request.CompanyNumber, request.OfficerID)
			fmt.Fprint(os.Stdout, "[debug] func CreateAuthCodeRequest(): 9.1.2...")
			if err != nil {
				log.ErrorR(req, fmt.Errorf("error calling Oracle API to get officer: %v", err))
				m := models.NewMessageResponse("there was a problem communicating with the Oracle API")
				utils.WriteJSONWithStatus(w, req, m, http.StatusInternalServerError)
				return
			}
			fmt.Fprint(os.Stdout, "[debug] func CreateAuthCodeRequest(): 9.1.3...")
			if officerResponse == service.NotFound {
				m := models.NewMessageResponse("No officer found")
				utils.WriteJSONWithStatus(w, req, m, http.StatusNotFound)
				return
			}
			fmt.Fprint(os.Stdout, "[debug] func CreateAuthCodeRequest(): 9.1.4...")
			request.OfficerUraID = officer.UsualResidentialAddress.ID
			request.OfficerForename = officer.Forename
			request.OfficerSurname = officer.Surname
			fmt.Fprint(os.Stdout, "[debug] func CreateAuthCodeRequest(): 9.1.1...")
		} else {
			fmt.Fprint(os.Stdout, "[debug] func CreateAuthCodeRequest(): 9.2.1...")
			// check if any eligible officers exist for specified company
			companyIsEligible, err := service.CheckOfficers(request.CompanyNumber)
			fmt.Fprint(os.Stdout, "[debug] func CreateAuthCodeRequest(): 9.2.2...")
			if err != nil {
				utils.WriteErrorMessage(w, req, http.StatusInternalServerError, "there was a problem communicating with the Oracle API")
				return
			}
			fmt.Fprint(os.Stdout, "[debug] func CreateAuthCodeRequest(): 9.2.3...")
			if !companyIsEligible {
				utils.WriteResponseMessage(w, req, http.StatusNotFound, "corporate body has no eligible officers")
				return
			}
		}
		fmt.Fprint(os.Stdout, "[debug] func CreateAuthCodeRequest(): 10...")
		model := transformers.AuthCodeResourceRequestToDB(&request)
		fmt.Fprint(os.Stdout, "[debug] func CreateAuthCodeRequest(): 11...")

		fmt.Fprint(os.Stdout, "[debug] func CreateAuthCodeRequest(): 12...")
		companyName, err := service.GetCompanyName(request.CompanyNumber, req)
		fmt.Fprint(os.Stdout, "[debug] func CreateAuthCodeRequest(): 13...")
		if err != nil {
			log.ErrorR(req, fmt.Errorf("error getting company name: [%v]", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Fprint(os.Stdout, "[debug] func CreateAuthCodeRequest(): 14...")
		model.Data.CompanyName = companyName

		fmt.Fprint(os.Stdout, "[debug] func CreateAuthCodeRequest(): 15...")
		err = authCodeReqSvc.CreateAuthCodeRequest(model)
		if err != nil {
			log.ErrorR(req, fmt.Errorf("error creating Auth Code Request: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Fprint(os.Stdout, "[debug] func CreateAuthCodeRequest(): 16...")
		utils.WriteJSONWithStatus(w, req, transformers.AuthCodeRequestResourceDaoToResponse(model), http.StatusCreated)
	})
}

func validateCorporateBody(req *http.Request, authCodeReqSvc *service.AuthCodeRequestService, companyNumber string, email string) (bool, error) {

	// Check whether multiple submissions have been made for company
	corpBodyMultipleRequests, err := authCodeReqSvc.CheckMultipleCorporateBodySubmissions(companyNumber)
	if corpBodyMultipleRequests {
		log.InfoR(req, "Request already submitted for company number "+companyNumber)
		return false, err
	}
	if err != nil {
		return false, err
	}

	// Check whether user has made too many requests
	userExceededRequests, err := authCodeReqSvc.CheckMultipleUserSubmissions(email)
	if userExceededRequests {
		log.InfoR(req, "requests exceeded for user "+email)
		return false, err
	}
	if err != nil {
		return false, err
	}

	// Check whether company has made recent filings
	hasFiledWithinPeriod, err := service.CheckCompanyFilingHistory(companyNumber)
	if hasFiledWithinPeriod {
		log.InfoR(req, "Recent filings found for company number "+companyNumber)
		return false, err
	}
	// if err != nil {
	// 	return false, err
	// }

	return true, nil
}
