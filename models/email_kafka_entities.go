package models

// EmailSend represents the avro schema which can be found in the chs-kafka-schemas repo
type EmailSend struct {
	AppID        string `avro:"app_id"`
	MessageID    string `avro:"message_id"`
	MessageType  string `avro:"message_type"`
	Data         string `avro:"data"`
	EmailAddress string `avro:"email_address"`
	CreatedAt    string `avro:"created_at"`
}

// DataField represents the data that will be sent to the email consumer and eventually displayed in the email
type DataField struct {
	FilingDescription string `json:"filing_description"`
	To                string `json:"to"`
	Subject           string `json:"subject"`
	CHSURL            string `json:"chs_url"`
}
