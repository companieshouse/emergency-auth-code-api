package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/companieshouse/chs.go/authentication"
	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/emergency-auth-code-api/models"
	"github.com/companieshouse/emergency-auth-code-api/service"
	"github.com/companieshouse/emergency-auth-code-api/utils"
	"github.com/gorilla/mux"
)

// handleEmailKafkaMessage allows us to mock the call to sendEmailKafkaMessage for unit tests
var handleEmailKafkaMessage = service.SendEmailKafkaMessage

const submitted = "submitted"

// UpdateAuthCodeRequest updates an auth code request for a specified auth-code-request ID
func UpdateAuthCodeRequest(authCodeSvc *service.AuthCodeService, authCodeReqSvc *service.AuthCodeRequestService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		var request models.AuthCodeRequest
		err := json.NewDecoder(req.Body).Decode(&request)

		// request body failed to get decoded
		if err != nil {
			utils.WriteErrorMessage(w, req, http.StatusBadRequest, "failed to read request body")
			return
		}

		userDetails := req.Context().Value(authentication.ContextKeyUserDetails)
		if userDetails == nil {
			utils.WriteErrorMessage(w, req, http.StatusBadRequest, "user details not in request context")
			return
		}

		// Check for a auth-code-request ID in the request
		vars := mux.Vars(req)
		authCodeRequestID := vars["auth_code_request_id"]
		if authCodeRequestID == "" {
			utils.WriteErrorMessage(w, req, http.StatusBadRequest, "auth code request ID missing from request")
			return
		}

		if request.CompanyNumber == "" {
			utils.WriteErrorMessage(w, req, http.StatusBadRequest, "company number missing from request")
			return
		}

		if request.OfficerID == "" && request.Status != submitted {
			utils.WriteErrorMessage(w, req, http.StatusBadRequest, "no valid changes supplied")
			return
		}

		authCodeReqDao, authCodeReqStatus := authCodeReqSvc.GetAuthCodeReqDao(authCodeRequestID, request.CompanyNumber)
		if authCodeReqStatus != service.Success {
			utils.WriteErrorMessage(w, req, http.StatusInternalServerError, "error reading auth code request")
			return
		}

		if authCodeReqDao.Data.Status == submitted {
			utils.WriteErrorMessage(w, req, http.StatusBadRequest, "request already submitted")
			return
		}

		// Update officer details in Request if supplied
		if request.OfficerID != "" {

			// retrieve details for officer from oracle-query-api
			officer, officerResponse, err := service.GetOfficerDetails(request.CompanyNumber, request.OfficerID)
			if err != nil {
				log.ErrorR(req, fmt.Errorf("error calling Oracle API to get officer: %v", err))
				m := models.NewMessageResponse("there was a problem communicating with the Oracle API")
				utils.WriteJSONWithStatus(w, req, m, http.StatusInternalServerError)
				return
			}
			if officerResponse == service.NotFound {
				m := models.NewMessageResponse("No officer found")
				utils.WriteJSONWithStatus(w, req, m, http.StatusNotFound)
				return
			}

			responseType := authCodeReqSvc.UpdateAuthCodeRequestOfficer(
				authCodeReqDao,
				authCodeRequestID,
				officer,
			)

			if responseType != service.Success {
				switch responseType {
				case service.Error:
					utils.WriteErrorMessage(w, req, http.StatusInternalServerError, "error updating officer details in authcode request")
				case service.InvalidData:
					utils.WriteErrorMessage(w, req, http.StatusBadRequest, "error updating officer details in authcode request")
				default:
					utils.WriteErrorMessage(w, req, http.StatusInternalServerError, "error updating officer details in authcode request")
				}
				return
			}

			log.InfoR(req, "officer details updated in authcode request", log.Data{"company_number": request.CompanyNumber})
		}

		if request.Status == submitted {

			if authCodeReqDao.Data.OfficerID == "" {
				utils.WriteErrorMessage(w, req, http.StatusBadRequest, "officer details not supplied")
				return
			}

			companyHasAuthCode, err := authCodeSvc.CheckAuthCodeExists(request.CompanyNumber)
			if err != nil {
				log.ErrorR(req, fmt.Errorf("error retrieving Auth Code from DB: %v", err))
				utils.WriteErrorMessage(w, req, http.StatusInternalServerError, "error retrieving Auth Code from DB")
				return
			}

			responseType := authCodeReqSvc.SendAuthCodeRequest(
				authCodeReqDao,
				request.CompanyNumber,
				userDetails.(authentication.AuthUserDetails).Email,
				companyHasAuthCode,
			)

			if responseType == service.NotFound {
				utils.WriteErrorMessage(w, req, http.StatusNotFound, "officer not found")
				return
			}

			if responseType != service.Success {
				utils.WriteErrorMessage(w, req, http.StatusInternalServerError, "error sending queue item")
				return
			}

			err = sendConfirmationEmail(userDetails.(authentication.AuthUserDetails).Email, req)
			if err != nil {
				log.ErrorR(req, err)
			}

			authCodeStatusResponseType := authCodeReqSvc.UpdateAuthCodeRequestStatusSubmitted(authCodeReqDao, authCodeRequestID, companyHasAuthCode)

			if authCodeStatusResponseType != service.Success {
				utils.WriteErrorMessage(w, req, http.StatusInternalServerError, "error updating status")
				return
			}

			log.InfoR(req, "status updated in authcode request; queue item submitted.", log.Data{"company_number": request.CompanyNumber})

		}

		response, responseType := authCodeReqSvc.GetAuthCodeRequest(authCodeRequestID)

		if responseType != http.StatusOK {
			utils.WriteErrorMessage(w, req, http.StatusInternalServerError, "error reading authcode request")
			return
		}

		utils.WriteJSONWithStatus(w, req, response, http.StatusOK)

	})
}

func sendConfirmationEmail(emailAddress string, r *http.Request) error {
	// Send confirmation email
	if err := handleEmailKafkaMessage(emailAddress); err != nil {
		return fmt.Errorf("error sending email to %s: %v", emailAddress, err)
	}

	log.InfoR(r, "confirmation email sent to customer", log.Data{
		"email_address": emailAddress,
	})

	return nil
}
