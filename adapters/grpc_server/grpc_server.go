package grpcserver

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/mikew79/port-domain-service/ports"
	"google.golang.org/grpc"

	pb "github.com/mikew79/port-domain-service/proto"
)

// PortDomainServer - An Adpater providing the GRPC service for API access
type PortDomainServer struct {
	ds ports.PortsDomainRepositioryPort
	pb.UnimplementedPortsDomainServer
}

// NewPortDomainServer - Create a new Server with our datastore
func NewPortDomainServer(ds ports.PortsDomainRepositioryPort) *PortDomainServer {
	return &PortDomainServer{ds: ds}
}

// RunServer - Run the service
func RunServer(ctx context.Context, api pb.PortsDomainServer, port int) error {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	log.Printf("Server Running")

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

// GetPort - Returns one port from the backend datastore by Id
func (s *PortDomainServer) GetPort(ctx context.Context, port *pb.Port) (*pb.Port, error) {
	return s.ds.FindOne(ctx, port)
}

// CreatePort - Inserts A new port into the backend datastore,
func (s *PortDomainServer) CreatePort(ctx context.Context, port *pb.Port) (*pb.UpdateResponse, error) {
	log.Printf("Create Port Called")
	return s.ds.CreateAndUpdate(ctx, port)
}

// UpdatePort - Updates a port with given Id in the backend datastore, if port doe not exist it is created
func (s *PortDomainServer) UpdatePort(ctx context.Context, port *pb.Port) (*pb.UpdateResponse, error) {
	return s.ds.CreateAndUpdate(ctx, port)
}

// DeletePort - Deletes a port with given Id from the backend datastore
func (s *PortDomainServer) DeletePort(ctx context.Context, port *pb.Port) (*pb.UpdateResponse, error) {
	return s.ds.Delete(ctx, port)
}

//CreateUpdatePorts - inserts mulitple ports into the backednd datbase from a stream
func (s *PortDomainServer) CreateUpdatePorts(stream pb.PortsDomain_CreateUpdatePortsServer) error {
	var count int32 = 0
	for {
		port, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.UpdateResponse{Count: count})
		}
		if err != nil {
			return err
		}
		_, err = s.ds.CreateAndUpdate(context.TODO(), port)
		if err == nil {
			count++
		}
	}
}

// ListPorts - Returns a numebr of ports form the backend databse as a stream
func (s *PortDomainServer) ListPorts(req *pb.ListRequest, stream pb.PortsDomain_ListPortsServer) error {

	err := s.ds.List(context.TODO(), req.Count)
	if err != nil {
		log.Fatal(err)
		return err
	}

	for {
		port, err := s.ds.ListNext(context.TODO())
		if err != nil {
			if err.Error() == "no more ports" {
				break
			} else {
				return err
			}
		}
		if err := stream.Send(port); err != nil {
			return err
		}
	}
	return nil
}
