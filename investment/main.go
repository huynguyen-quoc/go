package main

import (
	"fmt"
	"github.com/huynguyen-quoc/go/investment/protocol/grpc"
	"github.com/huynguyen-quoc/go/investment/protocol/rest"
	"github.com/huynguyen-quoc/go/investment/server"
	"os"
)


func main() {

	grpcServer := grpc.Grpc{}

	httpServer := rest.Http{}


	serverRun := &server.Server{
		GrpcRunner: grpcServer,
		HttpRunner: httpServer,
	}
	if err := serverRun.RunServer(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

