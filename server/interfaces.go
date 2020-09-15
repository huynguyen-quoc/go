package server

import (
	"context"
	"net/http"

	"google.golang.org/grpc"
)


// EncodePB wraps methods for encode some object to protobuf object
type EncodePB interface {
	ToPB(ctx context.Context) (interface{}, error)
}

// EncodeDTO wraps methods for encode some object to DTO object
type EncodeDTO interface {
	ToDTO(ctx context.Context) (interface{}, error)
}

// RPCServer wraps methods for general gRPC server
type RPCServer interface {
	// Start will be called before Hybrid starts RPC server
	Start(context.Context, *grpc.Server) error

	// Stop will be called before Hybrid shuts down RPC server
	Stop(context.Context)
}

// HTTPServer wraps methods for general HTTP server
type HTTPServer interface {
	// Start will be called before Hybrid starts HTTP server
	Start(context.Context, *http.Server) error

	// Stop will be called before Hybrid shuts down HTTP server
	Stop(context.Context)
}

// HybridServer wraps methods for server which can run gRPC server as main server and others server
type HybridServer interface {
	// Start starts a gRPC server
	Start(ctx context.Context, serviceName string, addr string) error

	// GetServer returns the RPC server managed by Hybrid
	// Must be called after Run
	GetServer() RPCServer

	// Stopped signals if has stopped the RPC server
	Stopped() <-chan struct{}

	// Stop manually stops
	Stop()
}
