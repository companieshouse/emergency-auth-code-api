package handlers

import "net/http"

// GetCompanyDirectors returns a list of valid company directors that can apply for an auth code
func GetCompanyDirectors(w http.ResponseWriter, req *http.Request) {

	// TODO :- Add logic to return company directors

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
