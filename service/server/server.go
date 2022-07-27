package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"

	pb "github.com/mikew79/port-domain-service/proto"
)

type portDomainServer struct {
	db *mongo.Database
	pb.UnimplementedPortsDomainServer
}

// Create a new Server
func NewPortDomainServer(db *mongo.Database) *portDomainServer {
	return &portDomainServer{db: db}
}

// Run the service
func RunServer(ctx context.Context, api pb.PortsDomainServer, port int) error {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	// register our service
	server := grpc.NewServer()
	pb.RegisterPortsDomainServer(server, api)

	// ensure agraceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a CTRL+C, handle it
			log.Println("Shutting Down gRPC Server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	log.Println("Starting gRPC Server...")
	return server.Serve(listen)
}

func (s *portDomainServer) GetPort(ctx context.Context, port *pb.Port) (*pb.Port, error) {
	collection := s.db.Collection("ports")
	var foundPort pb.Port
	filter := bson.M{"portid": port.Id}
	err := collection.FindOne(ctx, filter).Decode(&foundPort)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return &pb.Port{}, err
	}
	return &foundPort, nil
}

func (s *portDomainServer) CreatePort(ctx context.Context, port *pb.Port) (*pb.UpdateResponse, error) {
	collection := s.db.Collection("ports")

	opts := options.UpdateOptions{}
	opts.SetUpsert(true)

	filter := bson.M{"portid": port.Id}
	update := bson.M{"$set": port}
	result, err := collection.UpdateOne(ctx, filter, update, &opts)
	if err != nil {
		log.Fatal(err)
		return &pb.UpdateResponse{}, err
	}
	return &pb.UpdateResponse{Count: int32(result.UpsertedCount + result.ModifiedCount)}, nil
}

func (s *portDomainServer) UpdatePort(ctx context.Context, port *pb.Port) (*pb.UpdateResponse, error) {
	return s.CreatePort(ctx, port) // Create port is an upsert method, so we can use the same method
}

func (s *portDomainServer) DeletePort(ctx context.Context, port *pb.Port) (*pb.UpdateResponse, error) {
	collection := s.db.Collection("ports")

	filter := bson.M{"portid": port.Id}
	opts := options.DeleteOptions{}
	result, err := collection.DeleteOne(ctx, filter, &opts)
	if err != nil {
		log.Fatal(err)
		return &pb.UpdateResponse{}, err
	}
	return &pb.UpdateResponse{Count: int32(result.DeletedCount)}, nil
}

func (s *portDomainServer) CreateUpdatePorts(stream pb.PortsDomain_CreateUpdatePortsServer) error {
	collection := s.db.Collection("ports")

	opts := options.UpdateOptions{}
	opts.SetUpsert(true)
	var count int32 = 0
	for {
		port, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.UpdateResponse{Count: count})
		}
		if err != nil {
			return err
		}
		filter := bson.M{"portid": port.Id}
		update := bson.M{"$set": port}
		result, err := collection.UpdateOne(context.TODO(), filter, update, &opts)
		if err != nil {
			log.Fatal(err)
			continue
		}
		count += int32(result.UpsertedCount)
		count += int32(result.ModifiedCount)
	}
}

func (s *portDomainServer) ListPorts(req *pb.ListRequest, stream pb.PortsDomain_ListPortsServer) error {
	collection := s.db.Collection("ports")
	opts := options.Find()
	if req.Count > 0 {
		opts.SetLimit(int64(req.Count))
	}

	results, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
		return err
	}

	for results.Next(context.TODO()) {
		var port pb.Port

		if err := results.Decode(&port); err != nil {
			log.Fatal(err)
			return err
		}

		if err := stream.Send(&port); err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}
