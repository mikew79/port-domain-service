package portdomainserver

import (
	"context"
	"log"

	cfg "github.com/mikew79/port-domain-service/adapters/cmdlineconfig"
	server "github.com/mikew79/port-domain-service/adapters/grpc_server"
	mongo "github.com/mikew79/port-domain-service/adapters/mongodb_datastore"
	"github.com/mikew79/port-domain-service/application/entities"
)

//PortDomainServer - The core port domain server Application
type PortDomainServer struct {
	config     entities.Configuration
	repository mongo.MongoDbDataStore
}

//InitPortDomainService - Initialise the port domain server
func (pds *PortDomainServer) InitPortDomainService() error {

	// Get the configuration from the commandline loader
	var configLoader cfg.CommandLineArgs
	config, err := configLoader.GetConfiguration()
	if err != nil {
		return err
	}
	pds.config = config

	//Connect to our backend database
	err = pds.repository.Connect(pds.config)
	if err != nil {
		log.Fatalf("Error Connecting to repository: %s", err)
		return err
	}
	return nil
}

// Go - Run the port domain server
func (pds *PortDomainServer) Go() error {
	// run the GRPC Service
	api := server.NewPortDomainServer(&pds.repository)
	err := server.RunServer(context.Background(), api, pds.config.PortNumber)
	if err != nil {
		log.Fatalf("Error Running GRPC Server: %s", err)
		return err
	}
	return nil
}
