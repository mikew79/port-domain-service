# built the Protocol buffer files
protos:
	protoc --go_out=plugins=grpc:./service/proto - --proto_path=./proto port.proto

