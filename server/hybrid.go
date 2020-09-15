package server

import (
	"github.com/huynguyen-quoc/go/server/interrupt"
	"golang.org/x/net/context"
	"log"
	"os"
	"sync"
	"syscall"
)

var (
	sysSignalsToCapture = []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL}
)

// support both gRPC server and http Server
type hybridServerImpl struct {
	options          *hybridOptions
	interruptHandler interrupt.Handler

	// gRPC configuration
	rpcServer             RPCServer
	rpcServerShutdownChan chan struct{}

	wg       sync.WaitGroup
	stopChan chan struct{}

	cleanUpOnce sync.Once
}

func NewServer(rpcServer RPCServer, opts ...HybridOption) HybridServer {
	options := makeOptions(opts...)

	return &hybridServerImpl{
		rpcServer:             rpcServer,
		rpcServerShutdownChan: make(chan struct{}),
		stopChan:              make(chan struct{}),
		interruptHandler:      interrupt.NewInterruptHandler(),
		options:               options,
	}
}

func (h *hybridServerImpl) Start(ctx context.Context, serviceName string, addr string) error {
	h.interruptHandler.Start(sysSignalsToCapture...)
	go func() {
		defer h.Stop()
		defer close(h.rpcServerShutdownChan)
		h.interruptHandler.Wait()
	}()

	// if rpc server is not nil, start rpc server
	if h.rpcServer != nil {
		if err := h.startRPCServer(ctx, serviceName, addr); err != nil {
			log.Printf("Failed to start gRPC server. err=[%+v]", err)
			return err
		}
	}
	return nil
}

func (h *hybridServerImpl) GetServer() RPCServer {
	return h.rpcServer
}

func (h *hybridServerImpl) Stopped() <-chan struct{} {
	return h.stopChan
}

func (h *hybridServerImpl) Stop() {
	defer func() {
		h.wg.Wait()
		h.cleanUp()
	}()
}

func (h *hybridServerImpl) cleanUp() {
	h.cleanUpOnce.Do(func() {
		log.Printf("roi vao day")
		close(h.stopChan)
	})
}
