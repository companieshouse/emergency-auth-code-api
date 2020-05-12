package handlers

import "net/http"

// CreateAuthCodeRequest creates the auth code request for a specific officer ID
func CreateAuthCodeRequest(w http.ResponseWriter, req *http.Request) {

	// TODO :- Add logic to create Authorization Code Request

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
