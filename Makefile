help:
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@echo '    help               		Show this help screen.'
	@echo '    runserver              	Run the port domain server locally.'
	@echo '    runsampleclient-get    	Runs the sample client to retreive a sample port object.'
	@echo '    runsampleclient-insert   Runs the sample client to insert a sample port object.'
	@echo '    runsampleclient-update   Runs the sample client to update a sample port object.'
	@echo '    runsampleclient-delete   Runs the sample client to delete a sample port object.'
	@echo '    runsampleclient-list     Runs the sample client to list the first 10 ports in the datastore'
	@echo '    runsampleclient-stream   Runs the sample client to stream the sample data file to the server.'
	@echo '    rundockerserver          Spin up the production docker containers for the server using docker-compose'
	@echo '    lint                		Run golint.'
	@echo '    test             		Run the gRPC server unit tests.'
	@echo '    protos               	Build the protobuf files.'
	@echo ''
runserver: # Run the port domain server for debugging on the local machine
	go run cmd/portdomainserver/main.go -port=7000 -mongouri="mongodb://localhost:27017" -dbname="domainPortsDb"
runsampleclient-get: # Run the sample client making a get request
	go run cmd/sampleClient/main.go -get -port=7000 -id="AEJM"
runsampleclient-insert: # Run the sample client making an instert request
	go run cmd/sampleClient/main.go -create -port=7000 -id="AEJM" -name="Ajman" -city="Ajman" -country="United Arab Emirates" -alias="alias1,alias2" -regions="region1,region2" -coordinates="55.5136433,25.4052165" -province="Ajman" -timezone="Asia/Dubai" -unlocs="AEAJM" -code=52000
runsampleclient-update: # Run the sample client to update the name of a port created using runsampleclient-insert
	go run cmd/sampleClient/main.go -update -port=7000 -id="AEJM" -name="Ajman1" -city="Ajman" -country="United Arab Emirates" -alias="alias1,alias2" -regions="region1,region2" -coordinates="55.5136433,25.4052165" -province="Ajman" -timezone="Asia/Dubai" -unlocs="AEAJM" -code=52000
runsampleclient-delete: # run the sample client to delete an item created by the runsampleclient-insert command
	go run cmd/sampleClient/main.go -delete -port=7000 -id="AEJM"
runsampleclient-list: # run the sample cient, requsting thew first 10 ports, change count to get more, a count of 0 will return all ports
	go run cmd/sampleClient/main.go -list -port=7000 -count=10
runsampleclient-stream: # run the sample client and stream all the porst from the sample data file into the datab repository
	go run cmd/sampleClient/main.go -stream -port=7000 --file="data/ports.json"
rundockerserver: # Build and spin up the server in a docker container with mongo database container
	docker-compose up --build --detach
lint:	# run the go linter on the source
	golint ./...
test: # run the unit test for the core gRPC server
	go test ./adapters/grpc_server
protos: # build the protobuf files
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/ports.proto