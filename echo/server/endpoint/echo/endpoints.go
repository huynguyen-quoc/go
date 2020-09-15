package echo

import (
	"github.com/huynguyen-quoc/go/common/constant"
	"github.com/huynguyen-quoc/go/common/server"
	"github.com/huynguyen-quoc/go/common/server/endpoint"
	v1 "github.com/huynguyen-quoc/go/echo/pb/v1"
	"github.com/huynguyen-quoc/go/echo/service/echo"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type EndpointEchoHandler struct {
	Server server.Runner
}

func (h EndpointEchoHandler) makeServerEchoEndpoints(service echo.Service) *echo.Endpoints {
	var echoEndpoint endpoint.Endpoint
	userLoginEndpointMeta := endpoint.Meta{
		Name:        "echo",
		Type:        "server",
		Service:     "Echo",
		Idempotency: constant.IdempotencyIdempotent,
	}

	handler := echo.ServiceHandler{Service: service}
	{
		echoEndpoint = h.Server.Endpoint(userLoginEndpointMeta, handler.MakeEchoEndpoint())
	}

	return &echo.Endpoints{
		EchoEndpoint: echoEndpoint,
	}
}

// RegisterEchoHandlers
func (h EndpointEchoHandler) RegisterEchoHandlers(service echo.Service) {
	endpoints := h.makeServerEchoEndpoints(service)
	h.Server.RegisterGRPCEndpoints(makeEchoGRPCEndpoints(h.Server, []*echo.Endpoints{endpoints}))
	//h.server.RegisterHTTPEndpoints(makeEchoHTTPEndpoints(h.server, []*echo.Endpoints{endpoints}))
}

func makeEchoGRPCEndpoints(s server.Runner, endpoints []*echo.Endpoints) server.MakeGRPCEndpointsFunc {
	return func(ctx context.Context, rpcServer *grpc.Server) error {
		for _, endpoints := range endpoints {
			srv := echo.MakeGRPCServer(ctx, s, endpoints, server.Config.MDGetter)
			v1.RegisterEchoServiceServer(rpcServer, srv)
		}
		return nil
	}
}


//func makeEchoHTTPEndpoints(server server.Runner, endpoints []*echo.Endpoints) gkserver.MakeHTTPEndpointsFunc {
//	return func(ctx context.Context, r *mux.Router) error {
//		for _, endpoints := range endpointsList {
//			r = paysiaries.MakeGrabKitHTTPHandler(ctx, s, r, endpoints, grabkit.Config.MDGetter)
//		}
//		// Swagger docs
//
//		return nil
//	}
//}
