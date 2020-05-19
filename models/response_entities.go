package models

import "time"

// ResponseResource is the object returned in an error case
type ResponseResource struct {
	Message string `json:"message"`
}

// NewMessageResponse - convenience function for creating a response resource
func NewMessageResponse(message string) *ResponseResource {
	return &ResponseResource{Message: message}
}

// AuthCodeRequestResourceResponse is the entity returned in a
// successful response to creating an auth code request resource
type AuthCodeRequestResourceResponse struct {
	CompanyNumber string                       `json:"company_number"`
	CompanyName   string                       `json:"company_name"`
	UserID        string                       `json:"user_id"`
	UserEmail     string                       `json:"user_email"`
	OfficerID     string                       `json:"officer_id"`
	OfficerUraID  string                       `json:"officer_ura_id"`
	OfficerName   string                       `json:"officer_name"`
	Status        string                       `json:"status"`
	CreatedAt     *time.Time                   `json:"created_at"`
	SubmittedAt   *time.Time                   `json:"submitted_at"`
	Etag          string                       `json:"etag"`
	Kind          string                       `json:"kind"`
	Links         AuthCodeRequestResourceLinks `json:"links"`
}

// AuthCodeRequestResourceLinks is the links object of the payable resource
type AuthCodeRequestResourceLinks struct {
	Self string `json:"self"`
}
