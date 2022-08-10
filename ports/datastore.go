package ports

import (
	"context"

	"github.com/mikew79/port-domain-service/application/entities"
	pb "github.com/mikew79/port-domain-service/proto"
)

// PortsDomainRepositioryPort - Port providing connection to a datastroe reposoitory for the service
type PortsDomainRepositioryPort interface {
	Connect(cfg entities.Configuration) error
	FindOne(ctx context.Context, port *pb.Port) (*pb.Port, error)
	List(ctx context.Context, count int32) error
	ListNext(ctx context.Context) (*pb.Port, error)
	CreateAndUpdate(ctx context.Context, port *pb.Port) (*pb.UpdateResponse, error)
	Delete(ctx context.Context, port *pb.Port) (*pb.UpdateResponse, error)
}
