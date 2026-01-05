package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/companieshouse/emergency-auth-code-api/config"
	"github.com/companieshouse/emergency-auth-code-api/models"
	"github.com/companieshouse/filing-notification-sender/util"
)

const eacReceivedAppID = "emergency-auth-code-api.emergency_auth_code_request_received"
const eacFilingDescription = "Emergency Auth Code Request"
const eacMessageType = "emergency_auth_code_request_received"

func SendEmail(emailAddress string) error {
	cfg, err := config.Get()
	if err != nil {
		err = fmt.Errorf("error getting config for kafka message production: [%v]", err)
		return err
	}

	// Populate email details
	dataFieldMessage := models.DataField{
		FilingDescription: eacFilingDescription,
		To:                emailAddress,
		Subject:           fmt.Sprintf("Confirmation of your company authentication code request"),
		CHSURL:            cfg.CHSURL,
	}

	dataBytes, err := json.Marshal(dataFieldMessage)
	if err != nil {
		err = fmt.Errorf("error marshalling dataFieldMessage for emailSend: [%v]", err)
		return err
	}
	messageID := "<emergency-auth-code-request." + emailAddress + strconv.Itoa(util.Random(0, 100000)) + "@companieshouse.gov.uk>"
	emailSend := models.EmailSend{
		AppID:        eacReceivedAppID,
		MessageID:    messageID,
		MessageType:  eacMessageType,
		Data:         string(dataBytes),
		EmailAddress: emailAddress,
		CreatedAt:    time.Now().String(),
	}

	// Build email API request
	emailSendBytes, err := json.Marshal(emailSend)
	if err != nil {
		err = fmt.Errorf("error marshalling emailSend: [%v]", err)
		return err
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/send-email", cfg.ChsKafkaApiURL),
		bytes.NewBuffer(emailSendBytes))
	if err != nil {
		err = fmt.Errorf("error creating emailSend http request: [%v]", err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(cfg.APIKey, "")

	// Send email request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("error sending email: [%v]", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("wrong status code from kafka api when sending email: [%v]", resp.StatusCode)
		return err
	}

	return nil
}
