package mongodbdatastore

import (
	"context"
	"errors"
	"log"

	"github.com/mikew79/port-domain-service/application/entities"
	pb "github.com/mikew79/port-domain-service/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDbDataStore - An adapter to provide a connection to a mongoDB backend data repository for the Service
type MongoDbDataStore struct {
	client           *mongo.Client
	db               *mongo.Database
	results          *mongo.Cursor
	resultsAvailable bool
}

// Connect - Connect to the database
func (ds *MongoDbDataStore) Connect(cfg entities.Configuration) error {
	// connect to our mongo db backend
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("Failed To Connect to Database: %s", err)
		return err
	}

	// get the database
	db := client.Database(cfg.DbName)

	ds.client = client
	ds.db = db
	return nil
}

// FindOne - Find and return one port from the datastore, returns an error if nothing is found
func (ds *MongoDbDataStore) FindOne(ctx context.Context, port *pb.Port) (*pb.Port, error) {

	collection := ds.db.Collection("ports")
	var foundPort pb.Port
	filter := bson.M{"portid": port.Id}
	err := collection.FindOne(ctx, filter).Decode(&foundPort)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return &pb.Port{}, err
	}
	return &foundPort, nil
}

// List - Returns a list of all the ports from the database, limited by the count argument. if return is nil ListNext can be used to obtain the results
func (ds *MongoDbDataStore) List(ctx context.Context, count int32) error {

	collection := ds.db.Collection("ports")
	opts := options.Find()
	if count > 0 {
		opts.SetLimit(int64(count))
	}

	results, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
		return err
	}
	ds.results = results
	ds.resultsAvailable = true
	return nil
}

// ListNext - Returns the next available port from a listing
func (ds *MongoDbDataStore) ListNext(ctx context.Context) (*pb.Port, error) {
	if !ds.resultsAvailable {
		return &pb.Port{}, errors.New("no results available")
	}

	if !ds.results.Next(ctx) {
		return &pb.Port{}, errors.New("no more ports")
	}

	var foundPort pb.Port

	if err := ds.results.Decode(&foundPort); err != nil {
		log.Fatal(err)
		return &pb.Port{}, err
	}

	return &foundPort, nil
}

// CreateAndUpdate - Updates a port entry if one exists, if not a new entry is created
func (ds *MongoDbDataStore) CreateAndUpdate(ctx context.Context, port *pb.Port) (*pb.UpdateResponse, error) {
	collection := ds.db.Collection("ports")

	opts := options.UpdateOptions{}
	opts.SetUpsert(true)

	filter := bson.M{"portid": port.Id}
	update := bson.M{"$set": port}
	result, err := collection.UpdateOne(ctx, filter, update, &opts)
	if err != nil {
		log.Fatal(err)
		return &pb.UpdateResponse{}, err
	}
	return &pb.UpdateResponse{Count: int32(result.UpsertedCount + result.ModifiedCount)}, nil
}

// Delete - Deletes a port from the database with a given Port ID if one exists
func (ds *MongoDbDataStore) Delete(ctx context.Context, port *pb.Port) (*pb.UpdateResponse, error) {
	collection := ds.db.Collection("ports")

	filter := bson.M{"portid": port.Id}
	opts := options.DeleteOptions{}
	result, err := collection.DeleteOne(ctx, filter, &opts)
	if err != nil {
		log.Fatal(err)
		return &pb.UpdateResponse{}, err
	}
	return &pb.UpdateResponse{Count: int32(result.DeletedCount)}, nil
}
