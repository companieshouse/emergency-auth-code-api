package dao

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/emergency-auth-code-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func getMongoClient(mongoDBURL string) *mongo.Client {
	if client != nil {
		return client
	}

	ctx := context.Background()

	clientOptions := options.Client().ApplyURI(mongoDBURL)
	client, err := mongo.Connect(ctx, clientOptions)

	// Assume the caller of this func cannot handle the case where there is no database connection
	// so the service must crash here as it cannot continue.
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// Check we can connect to the mongodb instance. Failure here should result in a crash.
	pingContext, cancel := context.WithDeadline(ctx, time.Now().Add(5*time.Second))
	err = client.Ping(pingContext, nil)
	if err != nil {
		log.Error(errors.New("ping to mongodb timed out. please check the connection to mongodb and that it is running"))
		os.Exit(1)
	}
	defer cancel()

	log.Info("connected to mongodb successfully")

	return client
}

// MongoDatabaseInterface is an interface that describes the mongodb driver
type MongoDatabaseInterface interface {
	Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
}

func getMongoDatabase(mongoDBURL, databaseName string) MongoDatabaseInterface {
	return getMongoClient(mongoDBURL).Database(databaseName)
}

// MongoService is an implementation of the Service interface using MongoDB as the backend driver.
type MongoService struct {
	db             MongoDatabaseInterface
	CollectionName string
}

// CompanyHasAuthCode checks whether a company has an active auth code
func (m *MongoService) CompanyHasAuthCode(companyNumber string) (bool, error) {
	collection := m.db.Collection(m.CollectionName)
	dbResourceCount, err := collection.CountDocuments(context.Background(), bson.M{"_id": companyNumber, "is_active": true})
	if err != nil {
		return false, err
	}
	if dbResourceCount > 0 {
		return true, nil
	}
	return false, nil
}

// InsertAuthCodeRequest inserts an auth code request into the db
func (m *MongoService) InsertAuthCodeRequest(dao *models.AuthCodeRequestResourceDao) error {
	collection := m.db.Collection(m.CollectionName)
	_, err := collection.InsertOne(context.Background(), dao)
	return err
}

// UpdateAuthCodeRequestOfficer updates an authcode request with officer details
func (m *MongoService) UpdateAuthCodeRequestOfficer(dao *models.AuthCodeRequestResourceDao) error {
	collection := m.db.Collection(m.CollectionName)

	filter := bson.M{"_id": dao.ID}
	update := bson.M{
		"$set": bson.M{
			"data.officer_id":       dao.Data.OfficerID,
			"data.officer_forename": dao.Data.OfficerForename,
			"data.officer_surname":  dao.Data.OfficerSurname,
			"data.officer_ura_id":   dao.Data.OfficerUraID,
		},
	}

	_, err := collection.UpdateOne(context.Background(), filter, update)

	return err
}

// UpdateAuthCodeRequestStatus updates an authcode request with status details
func (m *MongoService) UpdateAuthCodeRequestStatus(dao *models.AuthCodeRequestResourceDao) error {
	collection := m.db.Collection(m.CollectionName)

	filter := bson.M{"_id": dao.ID}
	update := bson.M{
		"$set": bson.M{
			"data.status":       dao.Data.Status,
			"data.type":         dao.Data.Type,
			"data.submitted_at": dao.Data.SubmittedAt,
		},
	}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	return err
}

// GetAuthCodeRequest returns an auth code request from the db
func (m *MongoService) GetAuthCodeRequest(authCodeRequestID string) (*models.AuthCodeRequestResourceDao, error) {
	var resource models.AuthCodeRequestResourceDao

	collection := m.db.Collection(m.CollectionName)
	dbResource := collection.FindOne(context.Background(), bson.M{"_id": authCodeRequestID})

	err := dbResource.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Info("no auth-code-request found for id " + authCodeRequestID)
			return nil, nil
		}
		log.Error(err)
		return nil, err
	}

	err = dbResource.Decode(&resource)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &resource, nil
}

// CheckMultipleCorporateBodySubmissions checks for multiple company submitted requests.
// A maximum of one request every 3 days is permitted per company.
func (m *MongoService) CheckMultipleCorporateBodySubmissions(companyNumber string) (bool, error) {

	collection := m.db.Collection(m.CollectionName)
	dbResource := collection.FindOne(
		context.Background(),
		bson.M{
			"data.company_number": companyNumber,
			"data.status":         "submitted",
			"data.submitted_at":   bson.M{"$gt": time.Now().AddDate(0, 0, -3)},
		},
	)

	err := dbResource.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// CheckMultipleUserSubmissions checks whether a user has submitted multiple requests.
// A maximum of 3 user requests in a 24 hour period are permitted.
func (m *MongoService) CheckMultipleUserSubmissions(email string) (bool, error) {

	collection := m.db.Collection(m.CollectionName)
	submissionCount, err := collection.CountDocuments(
		context.Background(),
		bson.M{
			"data.created_by.user_email": email,
			"data.status":                "submitted",
			"data.submitted_at":          bson.M{"$gt": time.Now().AddDate(0, 0, -1)},
		},
	)

	if err != nil {
		return false, err
	}

	return submissionCount >= 3, nil
}
