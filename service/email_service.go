package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/companieshouse/chs.go/avro"
	"github.com/companieshouse/chs.go/avro/schema"
	"github.com/companieshouse/chs.go/kafka/producer"
	"github.com/companieshouse/emergency-auth-code-api/config"
	"github.com/companieshouse/emergency-auth-code-api/models"
	"github.com/companieshouse/filing-notification-sender/util"
)

const eacReceivedAppID = "emergency-auth-code-api.emergency_auth_code_request_received"
const eacFilingDescription = "Emergency Auth Code Request"
const eacMessageType = "emergency_auth_code_request_received"

// ProducerTopic is the topic to which the email-send kafka message is sent
const ProducerTopic = "email-send"

// ProducerSchemaName is the schema which will be used to send the email-send kafka message with
const ProducerSchemaName = "email-send"

// SendEmailKafkaMessage sends a kafka message to the email-sender to send an email
func SendEmailKafkaMessage(emailAddress string) error {
	cfg, err := config.Get()
	if err != nil {
		err = fmt.Errorf("error getting config for kafka message production: [%v]", err)
		return err
	}

	// Get a producer
	kafkaProducer, err := producer.New(&producer.Config{Acks: &producer.WaitForAll, BrokerAddrs: cfg.BrokerAddr})
	if err != nil {
		err = fmt.Errorf("error creating kafka producer: [%v]", err)
		return err
	}
	emailSendSchema, err := schema.Get(cfg.SchemaRegistryURL, ProducerSchemaName)
	if err != nil {
		err = fmt.Errorf("error getting schema from schema registry: [%v]", err)
		return err
	}
	producerSchema := &avro.Schema{
		Definition: emailSendSchema,
	}

	// Prepare a message with the avro schema
	message, err := prepareKafkaMessage(emailAddress, *producerSchema)
	if err != nil {
		err = fmt.Errorf("error preparing kafka message with schema: [%v]", err)
		return err
	}

	// Send the message
	partition, offset, err := kafkaProducer.Send(message)
	if err != nil {
		err = fmt.Errorf("failed to send message in partition: %d at offset %d", partition, offset)
		return err
	}
	return nil
}

// prepareKafkaMessage generates the kafka message that is to be sent
func prepareKafkaMessage(emailAddress string, emailSendSchema avro.Schema) (*producer.Message, error) {
	cfg, err := config.Get()
	if err != nil {
		err = fmt.Errorf("error getting config: [%v]", err)
		return nil, err
	}

	// Set dataField to be used in the avro schema.
	dataFieldMessage := models.DataField{
		FilingDescription: eacFilingDescription,
		To:                emailAddress,
		Subject:           fmt.Sprintf("Confirmation of your company authentication code request"),
		CHSURL:            cfg.CHSURL,
	}

	dataBytes, err := json.Marshal(dataFieldMessage)
	if err != nil {
		err = fmt.Errorf("error marshalling dataFieldMessage: [%v]", err)
		return nil, err
	}

	messageID := "<emergency-auth-code-request." + emailAddress + strconv.Itoa(util.Random(0, 100000)) + "@companieshouse.gov.uk>"

	emailSendMessage := models.EmailSend{
		AppID:        eacReceivedAppID,
		MessageID:    messageID,
		MessageType:  eacMessageType,
		Data:         string(dataBytes),
		EmailAddress: emailAddress,
		CreatedAt:    time.Now().String(),
	}

	messageBytes, err := emailSendSchema.Marshal(emailSendMessage)
	if err != nil {
		err = fmt.Errorf("error marshalling email send message: [%v]", err)
		return nil, err
	}

	producerMessage := &producer.Message{
		Value: messageBytes,
		Topic: ProducerTopic,
	}
	return producerMessage, nil
}
