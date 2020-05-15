package models

// CompanyOfficers contains a list of valid company officers for the company.
type CompanyOfficers struct {
	Items      []Items `json:"items"`
	TotalCount int     `json:"total_count"`
}

// Items contains the details for a specific officer
type Items struct {
	ID        string `json:"id"`
	Forename1 string `json:"forename_1"`
	Forename2 string `json:"forename_2"`
	Surname   string `json:"surname"`
}
