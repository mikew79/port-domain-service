package api

import (
	"context"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/bcicen/jstream"
	"google.golang.org/grpc"

	pb "github.com/mikew79/port-domain-service/proto"
)

type portsDomainClient struct {
	client pb.PortsDomainClient
}

func NewPortDomainClient(conn *grpc.ClientConn) *portsDomainClient {
	return &portsDomainClient{client: pb.NewPortsDomainClient(conn)}
}

func (c *portsDomainClient) CreatePort(id string, name string, city string, country string, alias []string, region []string, coordinates []float64, province string, timezone string, unlocs []string, code string) (string, error) {
	port := pb.Port{Id: id, Name: name, City: city, Country: country, Alias: alias, Regions: region, Coordinates: coordinates, Province: province, Timezone: timezone, Unlocs: unlocs, Code: code}

	response, err := c.client.CreatePort(context.Background(), &port)
	if err != nil {
		fmt.Printf("Error when calling CreatePort: %s\n", err)
		return "", err
	}

	return fmt.Sprintf("Response from server: %d\n", response.Count), nil
}

func (c *portsDomainClient) UpdatePort(id string, name string, city string, country string, alias []string, region []string, coordinates []float64, province string, timezone string, unlocs []string, code string) (string, error) {
	port := pb.Port{Id: id, Name: name, City: city, Country: country, Alias: alias, Regions: region, Coordinates: coordinates, Province: province, Timezone: timezone, Unlocs: unlocs, Code: code}

	response, err := c.client.UpdatePort(context.Background(), &port)
	if err != nil {
		fmt.Printf("Error when calling UpdatePort: %s\n", err)
		return "", err
	}
	return fmt.Sprintf("Response from server: %d\n", response.Count), nil
}

func (c *portsDomainClient) GetPort(id string) (string, error) {
	port := pb.Port{Id: id}

	response, err := c.client.GetPort(context.Background(), &port)
	if err != nil {
		fmt.Printf("Error when calling GetPort: %s\n", err)
		return "", err
	}
	return fmt.Sprintf("Response from server: %v\n", response), nil
}

func (c *portsDomainClient) DeletePort(id string) (string, error) {
	port := pb.Port{Id: id}
	response, err := c.client.DeletePort(context.Background(), &port)
	if err != nil {
		fmt.Printf("Error when calling DeletePort: %s\n", err)
		return "", err
	}
	resp := fmt.Sprintf("Response from server: %d\n", response.Count)
	return resp, nil
}

func (c *portsDomainClient) ListPorts(count int32) error {
	stream, err := c.client.ListPorts(context.Background(), &pb.ListRequest{Count: count})
	if err != nil {
		fmt.Printf("Error when calling ListPorts: %s\n", err)
		return err
	}
	for {
		port, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Printf("%v Steam Error = %v", c, err)
		}

		fmt.Println(port)
	}
	return nil
}

func (c *portsDomainClient) CreateFromJson(jsonFile string) error {
	if _, err := os.Stat(jsonFile); err != nil {
		return err // file does not exist
	}

	// Open our file to stream ports data from
	file, err := os.Open(jsonFile)
	if err != nil {
		fmt.Printf("error opening file : %v", err)
		return err
	}
	// Opn the stream to send our ports
	stream, err := c.client.CreateUpdatePorts(context.Background())
	if err != nil {
		fmt.Printf("%v.RecordRoute(_) = _, %v", c, err)
		return err
	}

	decoder := jstream.NewDecoder(file, 1).EmitKV()
	for port := range decoder.Stream() {
		kvtest := port.Value.(jstream.KV)
		dPort := toDomainPort(kvtest.Key, kvtest.Value)

		if err := stream.Send(&dPort); err != nil {
			fmt.Printf("%v.Send(%v) = %v", stream, dPort, err)
		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}

	if err := file.Close(); err != nil {
		fmt.Printf("error closing file : %v", err)
	}
	fmt.Println(reply)
	return nil
}

// This helper method converts the json port data stream from the input file to
func toDomainPort(id string, jsonPort interface{}) pb.Port {
	var port pb.Port = pb.Port{Id: id}
	iter := reflect.ValueOf(jsonPort).MapRange()
	for iter.Next() {
		key := iter.Key().Interface()
		value := iter.Value().Interface()
		switch key {
		case "name":
			port.Name = value.(string)
		case "city":
			port.City = value.(string)
		case "country":
			port.Country = value.(string)
		case "alias":
			port.Alias = reflectValueToStringSlice(value.([]interface{}))
		case "regions":
			port.Regions = reflectValueToStringSlice(value.([]interface{}))
		case "coordinates":
			port.Coordinates = reflectValueToFloat64Slice(value.([]interface{}))
		case "province":
			port.Province = value.(string)
		case "timezone":
			port.Timezone = value.(string)
		case "unlocs":
			port.Unlocs = reflectValueToStringSlice(value.([]interface{}))
		case "code":
			port.Code = value.(string)
		}
	}
	fmt.Println(port)
	return port
}

//convert the various strign lists to valid slices
func reflectValueToStringSlice(value []interface{}) []string {
	strValues := make([]string, 0, len(value))
	for _, val := range value {
		strValues = append(strValues, val.(string))
	}
	return strValues
}

// convert our streamed location values to the float 64 slice
func reflectValueToFloat64Slice(value []interface{}) []float64 {
	floatValues := make([]float64, 0, len(value))
	for _, val := range value {
		floatValues = append(floatValues, val.(float64))
	}
	return floatValues
}
