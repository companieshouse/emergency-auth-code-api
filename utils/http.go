package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/emergency-auth-code-api/models"
)

// WriteJSON writes the interface as a json string with status of 200.
func WriteJSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	WriteJSONWithStatus(w, r, data, http.StatusOK)
}

// WriteErrorMessage logs an error and adds it to the response, along with the supplied status
func WriteErrorMessage(w http.ResponseWriter, req *http.Request, status int, message string) {
	log.ErrorR(req, fmt.Errorf(message))
	WriteJSONWithStatus(w, req, models.NewMessageResponse(message), status)
}

// WriteResponseMessage writes a supplied message and status to the response, and logs an info message
func WriteResponseMessage(w http.ResponseWriter, req *http.Request, status int, message string) {
	log.InfoR(req, message)
	WriteJSONWithStatus(w, req, models.NewMessageResponse(message), status)
}

// WriteJSONWithStatus writes the interface as a json string with the supplied status.
func WriteJSONWithStatus(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.ErrorR(r, fmt.Errorf("error writing response: %v", err))
	}
}

// GetValueFromVars returns a specified value from the supplied request vars.
func GetValueFromVars(vars map[string]string, key string) (string, error) {
	val := vars[key]
	if val == "" {
		return "", fmt.Errorf("%s not found in vars", key)
	}
	return val, nil
}
