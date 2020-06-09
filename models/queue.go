package models

// QueueItem is authcode data to be sent to chs-queue-api
type QueueItem struct {
	Type          string  `json:"type"`
	Email         string  `json:"email"`
	CompanyNumber string  `json:"company_number"`
	CompanyName   string  `json:"company_name"`
	Address       Address `json:"ro_address"`
	Status        string  `json:"status"`
}

// Address is the address to which the authcode letter should be posted
type Address struct {
	POBox        string `json:"po_box,omitempty"`
	Premises     string `json:"premises,omitempty"`
	AddressLine1 string `json:"address_line_1,omitempty"`
	AddressLine2 string `json:"address_line_2,omitempty"`
	Locality     string `json:"locality,omitempty"`
	Region       string `json:"region,omitempty"`
	PostalCode   string `json:"postal_code,omitempty"`
	Country      string `json:"country,omitempty"`
}
