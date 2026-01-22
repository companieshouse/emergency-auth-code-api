package models

// EmailSend represents the json request body expected by chs-kafka-api
type EmailSend struct {
	AppID        string `json:"app_id"`
	MessageID    string `json:"message_id"`
	MessageType  string `json:"message_type"`
	Data         string `json:"json_data"`
	EmailAddress string `json:"email_address"`
	CreatedAt    string `json:"created_at"`
}

// DataField represents the data that will eventually be displayed in the email
type DataField struct {
	FilingDescription string `json:"filing_description"`
	To                string `json:"to"`
	Subject           string `json:"subject"`
	CHSURL            string `json:"chs_url"`
}
