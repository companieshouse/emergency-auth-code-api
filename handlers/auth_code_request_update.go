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

const submitted = "submitted"

// UpdateAuthCodeRequest updates an auth code request for a specified auth-code-request ID
func UpdateAuthCodeRequest(authCodeSvc *service.AuthCodeService, authCodeReqSvc *service.AuthCodeRequestService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		var request models.AuthCodeRequest
		err := json.NewDecoder(req.Body).Decode(&request)

		// request body failed to get decoded
		if err != nil {
			writeResponseWithMessage(w, req, http.StatusBadRequest, "failed to read request body")
			return
		}

		userDetails := req.Context().Value(authentication.ContextKeyUserDetails)
		if userDetails == nil {
			writeResponseWithMessage(w, req, http.StatusBadRequest, "user details not in request context")
			return
		}

		// Check for a auth-code-request ID in the request
		vars := mux.Vars(req)
		authCodeRequestID := vars["auth_code_request_id"]
		if authCodeRequestID == "" {
			writeResponseWithMessage(w, req, http.StatusBadRequest, "auth code request ID missing from request")
			return
		}

		if request.CompanyNumber == "" {
			writeResponseWithMessage(w, req, http.StatusBadRequest, "company number missing from request")
			return
		}

		if request.OfficerID == "" && request.Status != submitted {
			writeResponseWithMessage(w, req, http.StatusBadRequest, "no valid changes supplied")
			return
		}

		var response models.AuthCodeRequestResourceResponse

		authCodeReqDao, authCodeReqStatus := authCodeReqSvc.GetAuthCodeReqDao(authCodeRequestID, request.CompanyNumber)
		if authCodeReqStatus != service.Success {
			writeResponseWithMessage(w, req, http.StatusInternalServerError, "error reading auth code request")
			return
		}

		if authCodeReqDao.Data.Status == submitted {
			writeResponseWithMessage(w, req, http.StatusBadRequest, "request already submitted")
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

			officerUpdateResponse, responseType := authCodeReqSvc.UpdateAuthCodeRequestOfficer(
				authCodeReqDao,
				authCodeRequestID,
				officer,
			)

			if responseType != service.Success {
				switch responseType {
				case service.Error:
					writeResponseWithMessage(w, req, http.StatusInternalServerError, "error updating officer details in authcode request")
				case service.InvalidData:
					writeResponseWithMessage(w, req, http.StatusBadRequest, "error updating officer details in authcode request")
				default:
					writeResponseWithMessage(w, req, http.StatusInternalServerError, "error updating officer details in authcode request")
				}
				return
			}

			response = *officerUpdateResponse

			log.InfoR(req, "officer details updated in authcode request", log.Data{"company_number": request.CompanyNumber})
		}

		if request.Status == submitted {

			if authCodeReqDao.Data.OfficerID == "" {
				writeResponseWithMessage(w, req, http.StatusBadRequest, "officer details not supplied")
				return
			}

			companyHasAuthCode, err := authCodeSvc.CheckAuthCodeExists(request.CompanyNumber)
			if err != nil {
				log.ErrorR(req, fmt.Errorf("error retrieving Auth Code from DB: %v", err))
				writeResponseWithMessage(w, req, http.StatusInternalServerError, "error retrieving Auth Code from DB")
				return
			}

			responseType := authCodeReqSvc.SendAuthCodeRequest(
				authCodeReqDao,
				request.CompanyNumber,
				userDetails.(authentication.AuthUserDetails).Email,
				companyHasAuthCode,
			)

			if responseType == service.NotFound {
				writeResponseWithMessage(w, req, http.StatusNotFound, "officer not found")
				return
			}

			if responseType != service.Success {
				writeResponseWithMessage(w, req, http.StatusInternalServerError, "error sending queue item")
				return
			}

			authCodeStatusResponse, authCodeStatusResponseType := authCodeReqSvc.UpdateAuthCodeRequestStatus(authCodeReqDao, authCodeRequestID, request.Status, companyHasAuthCode)

			if authCodeStatusResponseType != service.Success {
				writeResponseWithMessage(w, req, http.StatusInternalServerError, "error updating status")
				return
			}

			response = *authCodeStatusResponse
		}

		utils.WriteJSONWithStatus(w, req, response, http.StatusCreated)

	})
}

func writeResponseWithMessage(w http.ResponseWriter, req *http.Request, status int, message string) {
	log.ErrorR(req, fmt.Errorf(message))
	m := models.NewMessageResponse(message)
	utils.WriteJSONWithStatus(w, req, m, status)
}
