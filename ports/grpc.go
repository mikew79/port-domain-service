package ports

// GRPCPort - Port providing access to the gRPC API service
type GRPCPort interface {
	RunServer() error
}
