package service

import (
	"context"
	"google.golang.org/grpc"
)

type remiGRPCImpl struct {
	makeEndpoints []MakeGRPCEndpointsFunc
	onInit        func(context.Context, *grpc.Server) error
}

func (r remiGRPCImpl) Start(ctx context.Context, server *grpc.Server) error {
	if r.makeEndpoints != nil {
		for _, fn := range r.makeEndpoints {
			if err := fn(ctx, server); err != nil {
				return err
			}
		}
	}

	if r.onInit != nil {
		if err := r.onInit(ctx, server); err != nil {
			return err
		}
	}

	return nil
}

func (r remiGRPCImpl) Stop(ctx context.Context) {
	// will do nothing
}
