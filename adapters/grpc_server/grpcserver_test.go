package grpcserver

import (
	"context"
	"testing"

	tds "github.com/mikew79/port-domain-service/adapters/testing_datastore"
	"github.com/mikew79/port-domain-service/application/entities"
	pb "github.com/mikew79/port-domain-service/proto"
)

type UpdateTests struct {
	id    string
	input pb.Port
	want  pb.UpdateResponse
}

func TestCoreServer(t *testing.T) {

	// create our in memory backend
	var repository tds.TestingDataStore
	repository.Connect(entities.Configuration{})

	// connect to our dummy backend
	api := NewPortDomainServer(&repository)

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

	// Sample objects used for table driven testing
	portUpdate1 := pb.Port{Id: "TESTID",
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

	portUpdate2 := pb.Port{Id: "ABCD",
		Name:        "Port Name Updated",
		City:        "Port City",
		Country:     "Port Country",
		Alias:       []string{"Alias1", "Alias2"},
		Regions:     []string{"Region 1", "region 2"},
		Coordinates: []float64{55.5136433, 25.4052165},
		Province:    "Province",
		Timezone:    "GMT",
		Unlocs:      []string{"Unlocs1", "Unlocs2"},
		Code:        "12345"}

	// table driven tests for update
	// NOTE for illustrative puposes the second test in this table will fail with the current code, it is designed to show an issue in
	// that it is not posible to detect an update failure vs a creation, in reality the code should be fixed to make this test pass
	// A simple fix would be to modify the Updtae response to include a created and updated field, bgrpcccccccccccccccccccu tI wanted to leave this as is, to illustrate the additonal covergae
	// of adding the testing table portion, over the previous version, which did not detect this failure
	tests := []UpdateTests{
		{id: "Id exists", input: portUpdate1, want: pb.UpdateResponse{Count: 1}},         // ensure nothing is updated if id does not exist
		{id: "ID does not exist", input: portUpdate2, want: pb.UpdateResponse{Count: 0}}, // ensure nothing is updated if id does not exist
	}

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

	// Table driven testing for checking update of ports
	for _, test := range tests {
		res, err = api.UpdatePort(context.Background(), &test.input)
		if err != nil {
			t.Errorf("Error Updating Port with Service: %s", err)
		}

		if res.Count != test.want.Count {
			t.Errorf("UpdatePort(..) = %d: want %d - TestID: %s", res.Count, test.want.Count, test.id)
		}
	}

	res, err = api.DeletePort(context.Background(), &port)
	if err != nil {
		t.Errorf("Error Deleting Port with Service: %s", err)
	}

	if res.Count != 1 {
		t.Errorf("DeletePort(..) = %d; want 1", res.Count)
	}
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
