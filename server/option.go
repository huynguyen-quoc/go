package server

import (
	"google.golang.org/grpc"
	"time"
)

const (
	defaultMaxRecvMsgSize   = 1024 * 1024 * 20
	defaultMaxSendMsgSize   = 1024 * 1024 * 20
	defaultHTTPReadTimeout  = 10 * time.Second
	defaultHTTPWriteTimeout = 10 * time.Second
)

type hybridOptions struct {
	unaryInterceptor grpc.UnaryServerInterceptor
	maxRecvMsgSize   int
	maxSendMsgSize   int
}

type HybridOption func(options *hybridOptions)

// WithMaxRecvMsgSize configures hybrid gRPC server MaxRecvMsgSize
func WithMaxRecvMsgSize(size int) HybridOption {
	return func(o *hybridOptions) {
		o.maxRecvMsgSize = size
	}
}

// WithMaxRecvMsgSize configures hybrid gRPC server MaxSendMsgSize
func WithMaxSendMsgSize(size int) HybridOption {
	return func(o *hybridOptions) {
		o.maxSendMsgSize = size
	}
}

func WithUnaryInterceptor(ui grpc.UnaryServerInterceptor) HybridOption {
	return func(o *hybridOptions) {
		o.unaryInterceptor = ui
	}
}

func makeOptions(opt ...HybridOption) *hybridOptions {
	opts := defaultHybridOption()
	for _, o := range opt {
		o(opts)
	}

	return opts
}

func defaultHybridOption() *hybridOptions {
	return &hybridOptions{
		maxRecvMsgSize:   defaultMaxRecvMsgSize,
		maxSendMsgSize:   defaultMaxSendMsgSize,
	}
}
