package grpc

import (
	"context"
	"fmt"
	v1 "github.com/huynguyen-quoc/go/investment/pb/v1"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	"github.com/huynguyen-quoc/go/investment/config"
)

type Grpc struct {
	v1Api  v1.InvestmentServiceServer
	config *config.AppConfig
}

func (g *Grpc) Run(ctx context.Context) error {
	address := fmt.Sprintf(":%d", g.config.Grpc.Port)
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	// register service
	server := grpc.NewServer()


	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	log.Printf("starting gRPC server [%s]...\n", address)
	return server.Serve(listen)
}

func (g *Grpc) WithConfig(config *config.AppConfig) *Grpc {
	g.config = config
	return g
}

func (g *Grpc) WithInvestmentServiceServer(v1Api v1.InvestmentServiceServer) *Grpc {
	g.v1Api = v1Api
	return g
}
