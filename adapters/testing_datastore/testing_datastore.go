package testingdatastore

import (
	"context"
	"errors"

	"github.com/mikew79/port-domain-service/application/entities"
	pb "github.com/mikew79/port-domain-service/proto"
)

// TestingDataStore - An adapter to proide a connection to an in memory data store for testing for the Service
type TestingDataStore struct {
	store     []*pb.Port
	listCount int32
	listIndex int32
}

// Connect - Usually called to connect to the data store, but for testing just initialises the datastore
func (ds *TestingDataStore) Connect(cfg entities.Configuration) error {
	ds.listCount = -1
	ds.listIndex = -1
	return nil
}

// FindOne - Find and return one port from the datastore, returns an error if nothing is found
func (ds *TestingDataStore) FindOne(ctx context.Context, port *pb.Port) (*pb.Port, error) {
	if len(ds.store) <= 0 {
		return &pb.Port{}, errors.New("no documents in data store")
	}

	for _, p := range ds.store {
		if p.Id == port.Id {
			return p, nil
		}
	}
	return &pb.Port{}, errors.New("document not found")
}

// List - Returns a list of all the ports from the datastore, limited by the count argument. if return is nil ListNext can be used to obtain the results
func (ds *TestingDataStore) List(ctx context.Context, count int32) error {

	if len(ds.store) <= 0 {
		return errors.New("no documents in data store")
	}
	if count == 0 {
		ds.listCount = int32(len(ds.store))
	} else {
		ds.listCount = count
	}

	ds.listIndex = -1
	return nil
}

// ListNext - Returns the next available port from a listing
func (ds *TestingDataStore) ListNext(ctx context.Context) (*pb.Port, error) {
	if ds.listCount < 0 {
		return &pb.Port{}, errors.New("list not initialised")
	}

	ds.listIndex++
	if ds.listIndex < (ds.listCount - 1) {
		return ds.store[ds.listIndex], nil
	}
	return &pb.Port{}, errors.New("no more ports")
}

// CreateAndUpdate - Updates a port entry if one exists, if not a new entry is created
func (ds *TestingDataStore) CreateAndUpdate(ctx context.Context, port *pb.Port) (*pb.UpdateResponse, error) {
	for _, p := range ds.store {
		if p.Id == port.Id {
			p.Name = port.Name
			p.City = port.City
			p.Country = port.Country
			p.Alias = port.Alias
			p.Regions = port.Regions
			p.Coordinates = port.Coordinates
			p.Province = port.Province
			p.Timezone = port.Timezone
			p.Unlocs = port.Unlocs
			p.Code = port.Code
			return &pb.UpdateResponse{Count: 1}, nil
		}
	}
	ds.store = append(ds.store, port)
	return &pb.UpdateResponse{Count: 1}, nil
}

// Delete - Deletes a port from the datastore with a given Port ID if one exists
func (ds *TestingDataStore) Delete(ctx context.Context, port *pb.Port) (*pb.UpdateResponse, error) {
	if len(ds.store) <= 0 {
		return &pb.UpdateResponse{}, errors.New("no documents in data store")
	}
	index := -1
	for i, p := range ds.store {
		if p.Id == port.Id {
			index = i
			break
		}
	}
	if index >= 0 {
		ds.store = removeIndex(ds.store, index)
		return &pb.UpdateResponse{Count: 1}, nil
	}
	return &pb.UpdateResponse{Count: 0}, nil
}

func removeIndex(store []*pb.Port, index int) []*pb.Port {
	return append(store[:index], store[index+1:]...)
}
