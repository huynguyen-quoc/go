package service

import (
	"context"
	"log"

	"github.com/huynguyen-quoc/go/server"

	"google.golang.org/grpc"
)

type GRPCOpts struct {
	MakeEndpoints          []MakeGRPCEndpointsFunc
	EnableServiceDiscovery bool
	OnInit                 func(context.Context, *grpc.Server) error
	MaxRecvMsgSize         int
	MaxSendMsgSize         int
}

// startRemiService starts the gRPC and HTTP servers using servicenanny
func (s *serviceImpl) startRemiService(ctx context.Context) error {
	var grpcServer *remiGRPCImpl
	nannyOpts := []server.HybridOption{
		server.WithMaxRecvMsgSize(s.configs.GRPC.MaxRecvMsgSize),
		server.WithMaxSendMsgSize(s.configs.GRPC.MaxSendMsgSize),
	}

	onShutdown := func(_ context.Context) {
		for i := len(s.configs.OnShutdown) - 1; i >= 0; i-- {
			s.configs.OnShutdown[i]()
		}
	}

	// Create gRPC server
	if s.configs.EnabledGRPC {
		grpcServer = s.newRPCServer()
	}

	// Construct nanny
	nanny := s.createServer(ctx, grpcServer, nannyOpts)

	// Start servers
	return s.startServerAndListen(ctx, onShutdown, nanny)
}

func (s *serviceImpl) newRPCServer() *remiGRPCImpl {
	opts := GRPCOpts{
		MakeEndpoints:  s.makeGRPCEndpoints,
		OnInit:         s.configs.GRPC.OnGRPCServerInit,
		MaxRecvMsgSize: s.configs.GRPC.MaxRecvMsgSize,
		MaxSendMsgSize: s.configs.GRPC.MaxSendMsgSize,
	}
	return &remiGRPCImpl{
		makeEndpoints: opts.MakeEndpoints,
		onInit:        opts.OnInit,
	}
}

func (s *serviceImpl) createServer(_ context.Context, grpcServer server.RPCServer, hybridOptions []server.HybridOption) server.HybridServer {
	var hybridServer server.HybridServer

	if s.configs.EnabledGRPC {
		hybridServer = server.NewServer(
			grpcServer,
			hybridOptions...,
		)
	} else {
		hybridServer = server.NewServer(
			nil,
			hybridOptions...,
		)
	}

	return hybridServer
}

func (s *serviceImpl) startServerAndListen(ctx context.Context, onShutdown func(context.Context), hybrid server.HybridServer) error {
	if err := hybrid.Start(ctx, s.configs.Name, s.configs.GRPC.Addr); err != nil {
		return err
	}

	go func() {
		select {
		case <-ctx.Done():
			hybrid.Stop()
			<-hybrid.Stopped()
		case <-hybrid.Stopped():
		}

		onShutdown(ctx)
		log.Printf("graceful stopped")
		s.wg.Done()
	}()

	return nil
}
