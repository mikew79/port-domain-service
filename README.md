# Ports Domain Server
 A gRPC micoservice hosted in Docker environment with Docker Mongo Db service for the databse backend. A Sample client is provided. Written in Golang.

## Summary
___
This repository contains 1 microservice, running the ports domain service, which provides a host docker container to run this in. Along with this is a docker-compose file that will spin up both a Mongodb container and the micorservice hosted in the container.

A sample client is then provided to allow you to test the microservice, from the command line. 

## Getting Started
___

## Makefile
The makefile contains scripts to test, and deploy the Ports Domain service, for more details run `make help` to list all options available

## Docker Compose

Docker Compose is the easiest way of starting the micorservice and a database.

* First build the protocol buffer definitions, by running `make protos` in the root of the repository
* Then run the two docker containers using `docker-compose up --detatch` in the root of the repository


## Using the client

The sample client is provided in the `/cmd/sampleClient` folder rof the repository.

This can be run using `go run /cmd/sampleClient/main.go` from the client folder.

There are a number of command line arguments avialbe for the client

Firstly set the relvant argument for the action you wish to carry out
```
get         - Gets a port by given ID
create      - Create a new port with the given data
update      - Update a port of given ID with given data
delete      - Deleete a port by given ID
list        - List all ports in the database
stream      - Bulk create or update porst from json file
```

A number of other arguments allow tou to control how the base commands function

```
count=<number>      - This when used with the list command will limit how many response are retunred
file=<filename>     - This is required when using the stream option and is the JSON file with the ports data to stream
id=<portID>     - This is required for any DRUD operation and is The ID of the port to transact with, e.g.AEAJM
port=<port>   - This is the port number of the ports domain service default is 7000
```

When adding or updating a record you can pass the port data using the following arguments
```
	name=<value>            - The name of the port object
	city=<value>            - The city for this port object
	country=<value>         - The country of the port object
	alias=<value>           - Aliases for the port object, A comma separated list of values
	regions=<value>         - Regions of the port object, A comma separated list of values
	coordinates=<value>     - The Coordinates for this port object, A comma separated list of longitude,lattitude
	province=<value>        - The province of this port object
	timezone=<value>        - The timezone for this port object
	unlocs=<value>          - The unlocs for this port object, A comma separated list of values
	code=<value>            - The code for this port object
```

