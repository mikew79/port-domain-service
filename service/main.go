package main

import (
	"context"
	"flag"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mikew79/port-domain-service/service/server"
)

type Configuration struct {
	mongoUri   string
	dbName     string
	portNumber int
}

func main() {

	var config Configuration

	// Command line arguments for configuring the service
	flag.StringVar(&config.mongoUri, "mongouri", "mongodb://mongodb:27017", "The URI of the mongodb databse server")
	flag.StringVar(&config.dbName, "dbname", "domainPortsDb", "The name to use for the mongodb Database, a new one will be created if it does not exist")
	flag.IntVar(&config.portNumber, "port", 7000, "The Port number to host the gRPC service on")
	flag.Parse()

	// connect to our mongo db backend
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.mongoUri))
	if err != nil {
		log.Fatalf("Failed To Connect to Database: %s", err)
	}

	// get the database
	db := client.Database(config.dbName)

	// run our service
	api := server.NewPortDomainServer(db)
	server.RunServer(context.Background(), api, config.portNumber)
}
