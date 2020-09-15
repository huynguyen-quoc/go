package service

import (
	"context"

	"github.com/huynguyen-quoc/go/server/endpoint"
)

type Runner interface {
	Start(context.Context) error
	Run(context.Context) error
	Stop() error
	Endpoint(endpoint.Meta, endpoint.Endpoint) endpoint.Endpoint
	RegisterGRPCEndpoints(fn MakeGRPCEndpointsFunc)
	Init(configs ...Config)
}
