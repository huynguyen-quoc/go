package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	"github.com/huynguyen-quoc/go/investment/config"
)

type Grpc struct {}

func (g Grpc) Run(ctx context.Context, config *config.AppConfig) error {
	address := fmt.Sprintf(":%d", config.Grpc.Port)
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
	log.Println("starting gRPC server...")
	return server.Serve(listen)
}

