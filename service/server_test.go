package main

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	pb "github.com/mikew79/port-domain-service/proto"
	"github.com/mikew79/port-domain-service/service/server"
)

func TestCoreServer(t *testing.T) {
	// connect to our mongo db backend
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Errorf("Failed To Connect to Database: %s", err)
	}

	db := client.Database("testDb")
	api := server.NewPortDomainServer(db)

	port := pb.Port{Id: "TESTID",
		Name:        "Port Name",
		City:        "Port City",
		Country:     "Port Country",
		Alias:       []string{"Alias1", "Alias2"},
		Regions:     []string{"Region 1", "region 2"},
		Coordinates: []float64{55.5136433, 25.4052165},
		Province:    "Province",
		Timezone:    "GMT",
		Unlocs:      []string{"Unlocs1", "Unlocs2"},
		Code:        "12345"}

	res, err := api.CreatePort(context.Background(), &port)
	if err != nil {
		t.Errorf("Error Creating Port with Service: %s", err)
	}

	if res.Count != 1 {
		t.Errorf("CreatePort(..) = %d; want 1", res.Count)
	}

	retPort, err := api.GetPort(context.Background(), &port)
	if err != nil {
		t.Errorf("Error Creating Port with Service: %s", err)
	}

	if retPort == nil {
		t.Errorf("GetPort(..) =  nil want Port")
	} else {
		if !comparePorts(*retPort, port) {
			t.Errorf("GetPort(..) = %v  want %v", *retPort, port)
		}
	}

	res, err = api.DeletePort(context.Background(), &port)
	if err != nil {
		t.Errorf("Error Deleting Port with Service: %s", err)
	}

	if res.Count != 1 {
		t.Errorf("DeletePort(..) = %d; want 1", res.Count)
	}

	// clean up after
	client.Database("testDb").Drop(context.TODO())
}

func comparePorts(port1 pb.Port, port2 pb.Port) bool {
	if port1.Id != port2.Id {
		return false
	}

	if port1.Name != port2.Name {
		return false
	}

	if port1.City != port2.City {
		return false
	}

	if port1.Name != port2.Name {
		return false
	}
	if port1.Country != port2.Country {
		return false
	}
	if len(port1.Alias) != len(port2.Alias) {
		return false
	}
	for i := 0; i < len(port1.Alias); i++ {
		if port1.Alias[i] != port2.Alias[i] {
			return false
		}
	}

	if len(port1.Regions) != len(port2.Regions) {
		return false
	}
	for i := 0; i < len(port1.Regions); i++ {
		if port1.Regions[i] != port2.Regions[i] {
			return false
		}
	}

	if len(port1.Coordinates) != len(port2.Coordinates) {
		return false
	}
	for i := 0; i < len(port1.Coordinates); i++ {
		if port1.Coordinates[i] != port2.Coordinates[i] {
			return false
		}
	}

	if port1.Province != port2.Province {
		return false
	}
	if port1.Timezone != port2.Timezone {
		return false
	}

	if len(port1.Unlocs) != len(port2.Unlocs) {
		return false
	}
	for i := 0; i < len(port1.Unlocs); i++ {
		if port1.Unlocs[i] != port2.Unlocs[i] {
			return false
		}
	}

	return port1.Code == port2.Code

}
