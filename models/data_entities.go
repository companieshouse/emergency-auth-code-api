package models

import (
	"time"
)

// AuthCodeRequestResourceDao is the persisted resource for auth code requests
type AuthCodeRequestResourceDao struct {
	ID   string                 `bson:"_id"`
	Data AuthCodeRequestDataDao `bson:"data"`
}

// AuthCodeRequestDataDao is the data of an auth code request resource
type AuthCodeRequestDataDao struct {
	CompanyNumber   string       `bson:"company_number"`
	CompanyName     string       `bson:"company_name"`
	OfficerID       string       `bson:"officer_id"`
	OfficerUraID    string       `bson:"officer_ura_id"`
	OfficerForename string       `bson:"officer_forename"`
	OfficerSurname  string       `bson:"officer_surname"`
	Status          string       `bson:"status"`
	CreatedAt       *time.Time   `bson:"created_at"`
	SubmittedAt     *time.Time   `bson:"submittedAt"`
	Kind            string       `bson:"kind"`
	Etag            string       `bson:"etag"`
	CreatedBy       CreatedByDao `bson:"created_by"`
	Type            string
	Links           AuthCodeResourceLinksDao `bson:"links"`
}

// CreatedByDao is the object relating to who created the resource
type CreatedByDao struct {
	Email    string `bson:"user_email"`
	ID       string `bson:"user_id"`
	Forename string `bson:"forename"`
	Surname  string `bson:"surname"`
}

// AuthCodeResourceLinksDao is the links object of the auth code resource
type AuthCodeResourceLinksDao struct {
	Self string `bson:"self"`
}
