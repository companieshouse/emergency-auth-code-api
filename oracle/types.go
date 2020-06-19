package oracle

// GetOfficersResponse returns the output of a get request for company officers
type GetOfficersResponse struct {
	ItemsPerPage int       `json:"items_per_page"`
	StartIndex   int       `json:"start_index"`
	TotalResults int       `json:"total_results"`
	Items        []Officer `json:"items"`
}

// Officer is a valid company officer
type Officer struct {
	ID                      string      `json:"id"`
	Forename                string      `json:"forename"`
	Surname                 string      `json:"surname"`
	OfficerRole             string      `json:"officer_role"`
	DateOfBirth             DateOfBirth `json:"date_of_birth"`
	AppointedOn             string      `json:"appointed_on"`
	Nationality             string      `json:"nationality"`
	CountryOfResidence      string      `json:"country_of_residence"`
	Occupation              string      `json:"occupation"`
	UsualResidentialAddress Address     `json:"usual_residential_address"`
}

// CompanyFilingCheck returns if a filing has happened against the company in a period determined by the oracle-query-api
type CompanyFilingCheck struct {
	EFilingFoundInPeriod bool `json:"efiling_found_in_period"`
}

// DateOfBirth is the month and year of an officer's date of birth
type DateOfBirth struct {
	Month string `json:"month"`
	Year  string `json:"year"`
}

// Address contains an officer's address
type Address struct {
	ID           string `json:"id"`
	PoBox        string `json:"po_box"`
	Premises     string `json:"premises"`
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	Locality     string `json:"locality"`
	Region       string `json:"region"`
	Postcode     string `json:"postcode"`
	Country      string `json:"country"`
}

// apiErrorResponse is the generic struct used to unmarshal the body of responses which have errored
type apiErrorResponse struct {
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"`
	Error     string `json:"string"`
	Message   string `json:"message"`
	Path      string `json:"path"`
}
