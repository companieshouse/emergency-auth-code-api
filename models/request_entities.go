package models

import "github.com/companieshouse/chs.go/authentication"

// AuthCodeRequest is the data received when creating a new Auth Code Request.
type AuthCodeRequest struct {
	CompanyNumber   string                         `json:"company_number" validate:"required"`
	CompanyName     string                         `json:"company_name" validate:"required"`
	CreatedBy       authentication.AuthUserDetails `json:",omitempty" validate:"required"`
	OfficerID       string                         `json:"officer_id" validate:"required"`
	Status          string                         `json:"status"`
	OfficerUraID    string                         `json:"officer_ura_id"`
	OfficerForename string                         `json:"officer_forename"`
	OfficerSurname  string                         `json:"officer_surname"`
}
