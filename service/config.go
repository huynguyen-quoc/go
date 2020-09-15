package service

import (
	"context"

	"google.golang.org/grpc"
)

type MakeGRPCEndpointsFunc func(ctx context.Context, rpcServer *grpc.Server) error

var (
	// DefaultGRPCAddr is the port for the gRPC server
	DefaultGRPCAddr = ":8087"
)

type GrpcConfig struct {
	Addr string
	// OnGRPCServerInit allows the user to add more routes to the GRPC router
	OnGRPCServerInit MakeGRPCEndpointsFunc `json:"-"`
	MaxRecvMsgSize   int
	MaxSendMsgSize   int
}

// Options ...
type Configs struct {
	Name        string
	EnabledGRPC bool
	GRPC        *GrpcConfig

	OnShutdown []func() `json:"-"`
}

func NewConfigs() Configs {
	return Configs{
		Name:        "Go service",
		EnabledGRPC: true,
		GRPC: &GrpcConfig{
			Addr: DefaultGRPCAddr,
		},
	}
}

type Config func(o *Configs)

// WithName
func WithName(name string) Config {
	return func(o *Configs) {
		o.Name = name
	}
}

// WithGRPC ...
func WithGRPC(enabled bool) Config {
	return func(o *Configs) {
		o.EnabledGRPC = enabled
	}
}

func WithGRPCAddr(address string) Config {
	return func(o *Configs) {
		if o.GRPC == nil {
			o.GRPC = &GrpcConfig{}
		}
		o.GRPC.Addr = address
	}
}

func WithGRPCMaxSendMsgSize(size int) Config {
	return func(o *Configs) {
		if o.GRPC == nil {
			o.GRPC = &GrpcConfig{}
		}
		o.GRPC.MaxSendMsgSize = size
	}
}

func WithGRPCMaxRecvMsgSize(size int) Config {
	return func(o *Configs) {
		if o.GRPC == nil {
			o.GRPC = &GrpcConfig{}
		}
		o.GRPC.MaxRecvMsgSize = size
	}
}
