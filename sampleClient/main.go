package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/mikew79/port-domain-service/sampleClient/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var getPort bool = false
	var createPort bool = false
	var updatePort bool = false
	var deletePort bool = false
	var listPorts bool = false
	var bulkAdd bool = false
	var portNumber int = 7000

	var portId = ""

	var portName string = ""
	var portCity string = ""
	var portCountry string = ""
	var portAlias string = ""
	var portRegions string = ""
	var portCoordinates string = ""
	var portProvince string = ""
	var portTimezone string = ""
	var portUnlocs string = ""
	var portCode string = ""

	var listCount int = 0
	var jsonFile string = ""

	flag.BoolVar(&getPort, "get", false, "Get a port by given ID")
	flag.BoolVar(&createPort, "create", false, "create a new port with the given data")
	flag.BoolVar(&updatePort, "update", false, "update a port of given ID with given data")
	flag.BoolVar(&deletePort, "delete", false, "Deleete a port by given ID")
	flag.BoolVar(&listPorts, "list", false, "list all ports ")
	flag.BoolVar(&bulkAdd, "stream", false, "bulk create or update porst from json file")

	flag.IntVar(&portNumber, "port", 7000, "The Port number of the hosted gRPC service")
	flag.IntVar(&listCount, "count", 0, "The number of ports to list")

	flag.StringVar(&jsonFile, "file", "", "The JSON file with the ports data to stream")

	flag.StringVar(&portId, "id", "", "The ID of the port to transact with")

	flag.StringVar(&portName, "name", "", "The name of the port object")
	flag.StringVar(&portCity, "city", "", "The city for this port object")
	flag.StringVar(&portCountry, "country", "", "The country of the port object")
	flag.StringVar(&portAlias, "alias", "", "Aliases for the port object")
	flag.StringVar(&portRegions, "regions", "", "Regions of the port object")
	flag.StringVar(&portCoordinates, "coordinates", "", "The Coordinates for this port object")
	flag.StringVar(&portProvince, "province", "", "The province of this port object")
	flag.StringVar(&portTimezone, "timezone", "", "The timezone for this port object")
	flag.StringVar(&portUnlocs, "unlocs", "", "The unlocs for this port object")
	flag.StringVar(&portCode, "code", "", "The code for this port object")

	flag.Parse()

	if (getPort || deletePort || updatePort) && portId == "" {
		fmt.Println("Argument -id must be specified for get and delete requests")
		return
	}

	portAliasSlice := strings.Split(portAlias, ",")
	portRegionSlice := strings.Split(portRegions, ",")
	portUnlocsSlice := strings.Split(portUnlocs, ",")

	// parse the cooodinates
	coords := strings.Split(portCoordinates, ",")
	var portCoordsSlice []float64
	if len(coords) == 2 {
		portCoordsSlice = make([]float64, 0, 2)
		val, _ := strconv.ParseFloat(coords[0], 64)
		portCoordsSlice = append(portCoordsSlice, val)
		val, _ = strconv.ParseFloat(coords[1], 64)
		portCoordsSlice = append(portCoordsSlice, val)
	}

	conn, err := grpc.Dial(fmt.Sprintf(":%d", portNumber), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Failed to connect to server")
		return
	}

	defer conn.Close()

	apiClient := api.NewPortDomainClient(conn)

	if getPort {
		result, err := apiClient.GetPort(portId)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result)

	} else if createPort {
		result, err := apiClient.CreatePort(portId, portName, portCity, portCountry, portAliasSlice, portRegionSlice, portCoordsSlice, portProvince, portTimezone, portUnlocsSlice, portCode)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result)

	} else if updatePort {
		result, err := apiClient.UpdatePort(portId, portName, portCity, portCountry, portAliasSlice, portRegionSlice, portCoordsSlice, portProvince, portTimezone, portUnlocsSlice, portCode)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result)
	} else if deletePort {
		result, err := apiClient.DeletePort(portId)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result)
	} else if listPorts {
		err := apiClient.ListPorts(int32(listCount))
		if err != nil {
			fmt.Println(err)
		}
	} else if bulkAdd {
		fmt.Println("Hello1")
		err := apiClient.CreateFromJson(jsonFile)
		if err != nil {
			fmt.Println(err)
		}
	}
}
