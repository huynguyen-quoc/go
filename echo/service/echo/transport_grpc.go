package echo

import (
	"github.com/huynguyen-quoc/go/common/server"
	"github.com/huynguyen-quoc/go/common/server/metadata"
	grpcTransport "github.com/huynguyen-quoc/go/common/server/transport/grpc"
	v1 "github.com/huynguyen-quoc/go/echo/pb/v1"
	"golang.org/x/net/context"
)

// MakeGRPCServer makes a set of endpoints available as a gRPC.
func MakeGRPCServer(ctx context.Context, s server.Runner, endpoints *Endpoints, mdGetter metadata.Getter) v1.EchoServiceServer {
	return makeGRPCServer(ctx, s, endpoints, mdGetter)
}

// makeGRPCServer makes a set of endpoints available as a gRPC GammaServer.
func makeGRPCServer(ctx context.Context, s server.Runner, endpoints *Endpoints, mdGetter metadata.Getter) v1.EchoServiceServer {
	return &grpcServer{
		echo: grpcTransport.NewTransport(
			ctx,
			endpoints.EchoEndpoint,
		),
	}
}

type grpcServer struct {
	echo grpcTransport.Handler
}

func (g grpcServer) Echo(ctx context.Context, request *v1.EchoRequest) (*v1.EchoResponse, error) {
	_, resp, err := g.echo.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*v1.EchoResponse), nil
}
