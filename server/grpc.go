package server

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func (h *hybridServerImpl) startRPCServer(ctx context.Context, serviceName string, address string) error {
	// Create RPC server and initialize custom server
	rpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(h.options.unaryInterceptor),
		grpc.MaxRecvMsgSize(h.options.maxRecvMsgSize),
		grpc.MaxSendMsgSize(h.options.maxSendMsgSize),
	)

	err := h.rpcServer.Start(ctx, rpcServer)
	if err != nil {
		return err
	}

	log.Printf("[RPC] starts listening to %s for service=[%s]", address, serviceName)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf("[RPC] failed to listen: %v", err)
		return err
	}

	// register health check gRPC
	healthCheckServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(rpcServer, healthCheckServer)

	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		log.Printf("[RPC] Server starts serving")

		if err := rpcServer.Serve(listener); err != nil {
			log.Printf("serving error: %s", err)
		}
		log.Printf("[RPC] Server stopped serving")
	}()

	return nil
}
