package grpc

import (
	"context"
	"errors"

	"github.com/huynguyen-quoc/go/server"

	"github.com/huynguyen-quoc/go/server/endpoint"
	serverError "github.com/huynguyen-quoc/go/server/errors"
	serverMetadata "github.com/huynguyen-quoc/go/server/metadata"
	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
)

// Handler which should be called from the gRPC binding of the service
// implementation. The incoming request parameter, and returned response
// parameter, are both gRPC types, not user-domain.
type Handler interface {
	ServeGRPC(ctx context.Context, request interface{}) (context.Context, interface{}, error)
	ServeGRPCStream(ctx context.Context, request interface{}) (context.Context, interface{}, error)
}

// Server wraps an endpoint and implements grpc.Handler.
type Server struct {
	ctx         context.Context
	endpoint    endpoint.Endpoint
	before      []RequestFunc
	after       []ResponseFunc
	sendHeaders func(context.Context) error
}

// NewServer constructs a new server *for a single endpoint*, which implements wraps the provided
// endpoint and implements the Handler interface. Consumers should write
// bindings that adapt the concrete gRPC methods from their compiled protobuf
// definitions to individual handlers. Request and response objects are from the
// caller business domain, not gRPC request and reply types.
func NewTransport(
	ctx context.Context,
	e endpoint.Endpoint,
	options ...ServerOption,
) *Server {
	s := &Server{
		ctx:         ctx,
		endpoint:    e,
		sendHeaders: sendHeaders,
	}
	for _, option := range options {
		option(s)
	}
	return s
}

// ServerOption sets an optional parameter for servers.
type ServerOption func(*Server)

// ServerBefore functions are executed on the gRPC request object before the
// request is decoded.
func ServerBefore(before ...RequestFunc) ServerOption {
	return func(s *Server) { s.before = before }
}

// ServerAfter functions are executed on the gRPC response writer after the
// endpoint is invoked, but before anything is written to the client.
func ServerAfter(after ...ResponseFunc) ServerOption {
	return func(s *Server) { s.after = after }
}

// ServeGRPC implements the Handler interface.
func (s Server) ServeGRPC(ctx context.Context, req interface{}) (context.Context, interface{}, error) {
	ctx = s.ctx

	// Retrieve gRPC ServerMetadata.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}

	// Use grpc ServerMetadata directly as our ServerMetadata
	ctx = convertMetaCtx(ctx, md)

	for _, f := range s.before {
		ctx = f(ctx, &md)
	}

	encodeDTO, ok := req.(server.EncodeDTO)
	if !ok {
		return ctx, nil, errors.New("request convert to DTO error")
	}

	request, err := encodeDTO.ToDTO(ctx)

	if err != nil {
		return ctx, nil, err
	}

	// Run endpoint + middleware
	_, response, err := s.endpoint(ctx, request)
	if err != nil {
		// encode custom errors
		if rErr, ok := err.(*serverError.Response); ok {
			return ctx, nil, rErr.ToRPCError()
		}
		return ctx, nil, err
	}

	for _, f := range s.after {
		ctx = f(ctx, &md)
	}

	encodePB, ok := response.(server.EncodePB)
	if !ok {
		return ctx, nil, errors.New("request convert to DTO error")
	}
	// Encode response
	grpcResp, err := encodePB.ToPB(ctx)
	if err != nil {
		return ctx, nil, err
	}

	// Add headers from outgoing ServerMetadata
	if err := s.sendHeaders(ctx); err != nil {
		return ctx, nil, err
	}

	return ctx, grpcResp, nil
}

func sendHeaders(ctx context.Context) error {
	md := serverMetadata.FromOutgoingContext(ctx)
	// Don't call SendHeader if there are none
	if len(md) == 0 {
		return nil
	}
	return grpc.SendHeader(ctx, headersToGRPC(serverMetadata.FromOutgoingContext(ctx)))
}

// ServeGRPCStream implements the Handler interface for grpc streams
func (s Server) ServeGRPCStream(ctx context.Context, req interface{}) (context.Context, interface{}, error) {
	ctx = s.ctx

	// Retrieve gRPC ServerMetadata.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}

	// Use grpc ServerMetadata directly as our ServerMetadata
	ctx = convertMetaCtx(ctx, md)

	for _, f := range s.before {
		ctx = f(ctx, &md)
	}

	// Decode request
	encodeDTO, ok := req.(server.EncodeDTO)
	if !ok {
		return ctx, nil, errors.New("request convert to DTO error")
	}

	request, err := encodeDTO.ToDTO(ctx)

	// Run endpoint + middleware
	_, response, err := s.endpoint(ctx, request)

	if err != nil {
		// encode custom errors
		if rErr, ok := err.(*serverError.Response); ok {
			return ctx, nil, rErr.ToRPCError()
		}
		return ctx, nil, err
	}

	for _, f := range s.after {
		ctx = f(ctx, &md)
	}

	return ctx, response, nil
}
