# emergency-auth-code-api

[![GoDoc](https://godoc.org/github.com/companieshouse/emergency-auth-code-api?status.svg)](https://godoc.org/github.com/companieshouse/emergency-auth-code-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/companieshouse/emergency-auth-code-api)](https://goreportcard.com/report/github.com/companieshouse/emergency-auth-code-api)

Temporary - Server side application for the emergency auth code solution.


## Requirements

In order to run this service locally, you will need:

- [Go](https://golang.org/doc/install)
- [Git](https://git-scm.com/downloads)
- [MongoDB](https://www.mongodb.com/)
 

## Getting Started

1. Clone this repository: `go get github.com/companieshouse/emergency-auth-code-api`
1. Build the executable: `make build`


## Configuration

Variable                            | Default | Description
:-----------------------------------|:-------:|:-----------
`BIND_ADDR`                         | `-`     | The host:port to bind to
`MONGODB_URL`                       | `-`     | The mongo DB connection string
`MONGO_AUTHCODE_DATABASE`           | `-`     | Authcode mongo database
`MONGO_AUTHCODE_COLLECTION`         | `-`     | Authcode mongo collection
`MONGO_AUTHCODE_REQUEST_DATABASE`   | `-`     | Authcode Request mongo database
`MONGO_AUTHCODE_REQUEST_COLLECTION` | `-`     | Authcode Request mongo collection
`ORACLE_QUERY_API_URL`              | `-`     | URL of the Oracle Query API
`QUEUE_API_LOCAL_URL`               | `-`     | URL of the Queue API


## Endpoints

Method   | Path                                                                         | Description
:--------|:-----------------------------------------------------------------------------|:-----------
**GET**  | `/emergency-auth-code-service/healthcheck`                                                               | Standard healthcheck endpoint
**GET**  | `emergency-auth-code-service/company/{company_number}/officers`              | Get list of eligible officers
**GET**  | `emergency-auth-code-service/company/{company_number}/officers/{officer_id}` | Get officer details
**POST** | `emergency-auth-code-service/auth-code-requests`                             | Create auth code request
**GET**  | `emergency-auth-code-service/auth-code-requests/{auth_code_request_id}`      | Get auth code request
**PUT**  | `emergency-auth-code-service/auth-code-requests/{auth_code_request_id}`      | Update auth code request
