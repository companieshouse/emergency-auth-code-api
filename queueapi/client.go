// Package queueapi connects and sends requests to the Queue API
package queueapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/emergency-auth-code-api/models"
)

// Client interacts with the Queue API
type Client struct {
	QueueAPIURL string
}

// NewClient will construct a new client service struct that can be used to interact with the Client API
func NewClient(queueAPIURL string) *Client {
	return &Client{
		QueueAPIURL: queueAPIURL,
	}
}

// sendRequest will make a http request and unmarshal the response body into a struct
func (c *Client) sendRequest(method, path string, item *models.QueueItem) (*http.Response, error) {
	url := c.QueueAPIURL + path

	reqBody, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(reqBody))

	logContext := log.Data{"request_method": method, "path": path}
	if err != nil {
		log.Error(err, logContext)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	// any errors here are due to transport errors, not 4xx/5xx responses
	if err != nil {
		log.Error(err, logContext)
		return nil, err
	}

	return resp, err
}

// SendQueueItem sends an item to the Queue API
func (c *Client) SendQueueItem(item *models.QueueItem) error {

	path := "/api/queue/authcode"

	resp, err := c.sendRequest(http.MethodPost, path, item)
	if err != nil {
		log.Error(fmt.Errorf("error sending request to queue API: %v", err))
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		log.Error(fmt.Errorf("error closing response body from Queue API: %v", err))
		// No need to return err, as sending request was successful
	}

	return nil
}
