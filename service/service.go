package service

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"context"

	"log"

	"github.com/huynguyen-quoc/go/server/endpoint"
)

const (
	gracefulShutdownPeriod = 5 * time.Second
)

type serviceImpl struct {
	makeGRPCEndpoints []MakeGRPCEndpointsFunc
	configs           Configs
	exit              chan chan error
	wg                sync.WaitGroup

	registeredGRPCAddr string
}

func NewService(configs ...Config) Runner {
	c := NewConfigs()
	for _, config := range configs {
		config(&c)
	}
	return &serviceImpl{
		configs:           c,
		exit:              make(chan chan error),
		makeGRPCEndpoints: []MakeGRPCEndpointsFunc{},
	}
}

func (s *serviceImpl) Start(ctx context.Context) error {
	s.wg.Add(1)

	ctx, cancelFunc := context.WithCancel(ctx)

	// ServiceNanny starts both the gRPC and HTTP servers
	if err := s.startRemiService(ctx); err != nil {
		cancelFunc()
		return err
	}

	// graceful shutdown
	go func() {
		waitAll := make(chan struct{})

		// Only exit when requested
		errCh := <-s.exit

		cancelFunc()

		// Wait for all servers to exit gracefully
		go func() {
			s.wg.Wait()
			waitAll <- struct{}{}
			log.Printf( "graceful shutdown")
		}()

		// Force stop if server exceeds the graceful shutdown time limit
		go func() {
			<-time.After(gracefulShutdownPeriod)
			waitAll <- struct{}{}
			log.Printf("forced shutdown")
		}()

		<-waitAll

		errCh <- nil
	}()

	return nil
}

func (s *serviceImpl) Run(ctx context.Context) error {
	if err := s.Start(ctx); err != nil {
		return err
	}
	// Interrupt / stop handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-c:
		log.Printf( "receive shut down [%+v]", sig)
	case msg := <-ctx.Done():
		log.Printf( "context finished [%s]", msg)
	}

	return s.Stop()
}

func (s *serviceImpl) Stop() error {
	ch := make(chan error)
	s.exit <- ch
	return <-ch
}

func (s *serviceImpl) Endpoint(_ endpoint.Meta, endpoint endpoint.Endpoint) endpoint.Endpoint {
	return endpoint
}

func (s *serviceImpl) RegisterGRPCEndpoints(fn MakeGRPCEndpointsFunc) {
	s.makeGRPCEndpoints = append(s.makeGRPCEndpoints, fn)
}

func (s *serviceImpl) Init(configs ...Config) {
	for _, opt := range configs {
		opt(&s.configs)
	}
}
