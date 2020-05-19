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
