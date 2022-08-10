package cmdlineconfig

import (
	"flag"

	"github.com/mikew79/port-domain-service/application/entities"
)

// CommandLineArgs - Adapter for parsing Service configuration from the command line
type CommandLineArgs struct {
}

// GetConfiguration - Returns the application configuration entity filled with config options from the command line or default values
func (cfg *CommandLineArgs) GetConfiguration() (entities.Configuration, error) {
	// Command line arguments for configuring the service
	var mongoURI string = ""
	var dbName string = ""
	var portNumber int = 0
	flag.StringVar(&mongoURI, "mongouri", "mongodb://mongodb:27017", "The URI of the mongodb databse server")
	//flag.StringVar(&mongoUri, "mongouri", "mongodb://localhost:27017", "The URI of the mongodb databse server")
	flag.StringVar(&dbName, "dbname", "domainPortsDb", "The name to use for the mongodb Database, a new one will be created if it does not exist")
	flag.IntVar(&portNumber, "port", 7000, "The Port number to host the gRPC service on")
	flag.Parse()

	return entities.Configuration{MongoURI: mongoURI, DbName: dbName, PortNumber: portNumber}, nil
}
