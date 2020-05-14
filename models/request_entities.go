package models

// AuthCodeRequest is the model that should be sent when creating a new payable request. It will contain a list of
// transactions along with their id and amount.
type AuthCodeRequest struct {
	CompanyNumber string `json:"company_number" validate:"required"`
	OfficerID     string `json:"officer_id" validate:"required"`
	Status        string `json:"status"`
}
