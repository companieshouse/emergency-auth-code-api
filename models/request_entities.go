package models

import "github.com/companieshouse/chs.go/authentication"

// AuthCodeRequest is the model that should be sent when creating a new payable request. It will contain a list of
// transactions along with their id and amount.
type AuthCodeRequest struct {
	CompanyNumber string                         `json:"company_number" validate:"required"`
	CreatedBy     authentication.AuthUserDetails `json:",omitempty" validate:"required"`
	OfficerID     string                         `json:"officer_id" validate:"required"`
	Status        string                         `json:"status"`
}
