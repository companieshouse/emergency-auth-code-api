// Package oracle connects to the Oracle Query API to retrieve officer details
package oracle

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/companieshouse/chs.go/log"
)

var (
	// ErrFailedToReadBody is a generic error when failing to parse a response body
	ErrFailedToReadBody = errors.New("failed reading the body of the response")
	// ErrOracleAPIBadRequest is a 400
	ErrOracleAPIBadRequest = errors.New("failed request to Oracle API")
	// ErrOracleAPIInternalServer is anything in the 5xx
	ErrOracleAPIInternalServer = errors.New("got an internal server error from Oracle API")
	// ErrOracleAPINotFound is a 404
	ErrOracleAPINotFound = errors.New("not found")
	// ErrUnexpectedServerError represents anything other than a 400, 404 or 500
	ErrUnexpectedServerError = errors.New("unexpected server error")
)

// Client interacts with the Oracle API
type Client struct {
	OracleAPIURL string
}

// GetOfficers will return a list of officers for a company
func (c *Client) GetOfficers(companyNumber string, startIndex string, itemsPerPage string) (*GetOfficersResponse, error) {

	logContext := log.Data{"company_number": companyNumber}

	path := fmt.Sprintf("/emergency-auth-code/company/%s/eligible-officers?start_index=%s&items_per_page=%s", companyNumber, startIndex, itemsPerPage)

	resp, err := c.sendRequest(http.MethodGet, path)

	// deal with any http transport errors
	if err != nil {
		log.Error(err, logContext)
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	defer resp.Body.Close()

	// determine if there are unexpected 4xx/5xx errors. an error here relates to a response parsing issue
	err = c.checkResponseForError(resp)
	if err != nil {
		log.Error(err, logContext)
		return nil, err
	}

	out := &GetOfficersResponse{}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err, logContext)
		return nil, ErrFailedToReadBody
	}

	err = json.Unmarshal(b, out)
	if err != nil {
		log.Error(err, logContext)
		return nil, ErrFailedToReadBody
	}

	return out, nil
}

// GetOfficer will return a single officer transactions for a company
func (c *Client) GetOfficer(companyNumber, officerID string) (*Officer, error) {

	logContext := log.Data{"company_number": companyNumber}

	path := fmt.Sprintf("/emergency-auth-code/company/%s/eligible-officers/%s", companyNumber, officerID)

	resp, err := c.sendRequest(http.MethodGet, path)

	// deal with any http transport errors
	if err != nil {
		log.Error(err, logContext)
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	defer resp.Body.Close()

	// determine if there are unexpected 4xx/5xx errors. an error here relates to a response parsing issue
	err = c.checkResponseForError(resp)
	if err != nil {
		log.Error(err, logContext)
		return nil, err
	}

	out := &Officer{}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err, logContext)
		return nil, ErrFailedToReadBody
	}

	err = json.Unmarshal(b, out)
	if err != nil {
		log.Error(err, logContext)
		return nil, ErrFailedToReadBody
	}

	return out, nil
}

// CheckFilingHistory will return details of the companies filing history
func (c *Client) CheckFilingHistory(companyNumber string) (*CompanyFilingCheck, error) {

	logContext := log.Data{"company_number": companyNumber}

	path := fmt.Sprintf("/emergency-auth-code/company/%s/efiling-status", companyNumber)

	resp, err := c.sendRequest(http.MethodGet, path)

	// deal with any http transport errors
	if err != nil {
		log.Error(err, logContext)
		return nil, err
	}

	defer resp.Body.Close()

	// determine if there are unexpected 4xx/5xx errors. an error here relates to a response parsing issue
	err = c.checkResponseForError(resp)
	if err != nil {
		log.Error(err, logContext)
		return nil, err
	}

	out := &CompanyFilingCheck{}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err, logContext)
		return nil, ErrFailedToReadBody
	}

	err = json.Unmarshal(b, out)
	if err != nil {
		log.Error(err, logContext)
		return nil, ErrFailedToReadBody
	}

	return out, nil
}

// Generic function which inspects the http response
// Returns the response struct or an error if there was a problem reading and parsing the body
func (c *Client) checkResponseForError(r *http.Response) error {

	if r.StatusCode == 200 {
		return nil
	}

	logContext := log.Data{
		"response_status": r.StatusCode,
	}

	// parse the error response and log all output
	e := &apiErrorResponse{}
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Error(err, logContext)
		return ErrFailedToReadBody
	}

	err = json.Unmarshal(b, e)
	if err != nil {
		log.Error(err, logContext)
		return ErrFailedToReadBody
	}

	d := log.Data{
		"status":  e.Status,
		"message": e.Message,
		"path":    e.Path,
	}

	log.Error(errors.New("error response from Oracle API query: status code returned = "+strconv.Itoa(r.StatusCode)), d)

	switch r.StatusCode {
	case http.StatusBadRequest:
		return ErrOracleAPIBadRequest
	case http.StatusInternalServerError:
		return ErrOracleAPIInternalServer
	default:
		return ErrUnexpectedServerError
	}
}

// sendRequest will make a http request and unmarshal the response body into a struct
func (c *Client) sendRequest(method, path string) (*http.Response, error) {
	url := c.OracleAPIURL + path
	req, err := http.NewRequest(method, url, nil)

	logContext := log.Data{"request_method": method, "path": path}
	if err != nil {
		log.Error(err, logContext)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	// any errors here are due to transport errors, not 4xx/5xx responses
	if err != nil {
		log.Error(err, logContext)
		return nil, err
	}

	return resp, err
}

// NewClient will construct a new client service struct that can be used to interact with the Client API
func NewClient(oracleAPIURL string) *Client {
	return &Client{
		OracleAPIURL: oracleAPIURL,
	}
}
