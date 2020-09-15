package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/huynguyen-quoc/go/investment/config"
)

type Http struct {
	Config *config.AppConfig
}

func (g Http) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	_ = []grpc.DialOption{grpc.WithInsecure()}

	address := fmt.Sprintf(":%d", g.Config.Http.Port)
	srv := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_ = srv.Shutdown(ctx)
	}()

	log.Printf("Starting HTTP server [%s]...\n", address)
	return srv.ListenAndServe()
}
