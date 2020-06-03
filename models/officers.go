package models

// OfficerListResponse is the Officer List to be returned
type OfficerListResponse struct {
	ItemsPerPage int       `json:"items_per_page"`
	StartIndex   int       `json:"start_index"`
	TotalResults int       `json:"total_results"`
	Items        []Officer `json:"items"`
}

// Officer is a single officer to be returned
type Officer struct {
	ID                 string      `json:"id"`
	Name               string      `json:"name"`
	OfficerRole        string      `json:"officer_role"`
	DateOfBirth        DateOfBirth `json:"date_of_birth"`
	AppointedOn        string      `json:"appointed_on"`
	Nationality        string      `json:"nationality"`
	CountryOfResidence string      `json:"country_of_residence"`
	Occupation         string      `json:"occupation"`
}

// DateOfBirth is the Date Of Birth of an officer
type DateOfBirth struct {
	Month string `json:"month"`
	Year  string `json:"year"`
}
