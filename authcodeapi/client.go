// Package authcodeapi connects and sends requests to the authcode api flow
package authcodeapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/emergency-auth-code-api/models"
)

// Client interacts with the AuthCode API
type Client struct {
	AuthCodeAPIURL  string
	AuthCodeAPIPath string
}

// NewClient will construct a new client service struct that can be used to interact with the Client API
func NewClient(authCodeAPIURL, authCodeAPIPath string) *Client {
	return &Client{
		AuthCodeAPIURL:  authCodeAPIURL,
		AuthCodeAPIPath: authCodeAPIPath,
	}
}

// sendRequest will make a http request and unmarshal the response body into a struct
func (c *Client) sendRequest(method, authCodeRequestID string, item *models.AuthCodeItem) (*http.Response, error) {
	reqBody, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, c.AuthCodeAPIURL+c.AuthCodeAPIPath, bytes.NewReader(reqBody))

	logContext := log.Data{"request_method": method, "path": c.AuthCodeAPIPath}
	if err != nil {
		log.Error(err, logContext)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-Id", authCodeRequestID)

	resp, err := http.DefaultClient.Do(req)
	// any errors here are due to transport errors, not 4xx/5xx responses
	if err != nil {
		log.Error(err, logContext)
		return nil, err
	}

	return resp, err
}

// SendAuthCodeItem sends an item to the AuthCode API
func (c *Client) SendAuthCodeItem(item *models.AuthCodeItem, authCodeRequestID string) error {
	resp, err := c.sendRequest(http.MethodPost, authCodeRequestID, item)
	if err != nil {
		log.Error(fmt.Errorf("error sending request to authCode API: %v", err))
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		log.Error(fmt.Errorf("error closing response body from AuthCode API: %v", err))
		// No need to return err here, as sending request might have been successful
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status returned from authCode API: %v", resp.StatusCode)
	}
	return nil
}
