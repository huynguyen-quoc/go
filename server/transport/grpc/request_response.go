package grpc

import (
	"context"
	"encoding/base64"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	binHdrSuffix = "-bin"
)

// RequestFunc may take information from an gRPC request and put it into a
// request context. In Servers, BeforeFuncs are executed prior to invoking the
// endpoint. In Clients, BeforeFuncs are executed after creating the request
// but prior to invoking the gRPC client.
type RequestFunc func(context.Context, *metadata.MD) context.Context

// ResponseFunc may take information from a request context and use it to
// manipulate the gRPC metadata header. ResponseFuncs are only executed in
// servers, after invoking the endpoint but prior to writing a response.
type ResponseFunc func(context.Context, *metadata.MD) context.Context

// SetResponseHeader returns a ResponseFunc that sets the specified metadata
// key-value pair.
func SetResponseHeader(key, val string) ResponseFunc {
	return func(ctx context.Context, md *metadata.MD) context.Context {
		key, val = EncodeKeyValue(key, val)
		(*md)[key] = append((*md)[key], val)
		return ctx
	}
}

// SetResponseHeadersIfAbsent returns a ServerResponseFunc that sets the given header.
func SetResponseHeadersIfAbsent(key string, vals []string) ResponseFunc {
	return func(ctx context.Context, _ *metadata.MD) context.Context {
		md, _ := metadata.FromOutgoingContext(ctx)
		v, ok := md[key]
		if !ok || len(v) == 0 {
			md[key] = vals
		}
		// Note we can't *send* the headers here because grpc.SendHeader can only be called once. Instead, add them to the
		// context so they can be sent later from the transport.
		return metadata.NewOutgoingContext(ctx, md)
	}
}

// SetRequestHeader returns a RequestFunc that sets the specified metadata
// key-value pair.
func SetRequestHeader(key, val string) RequestFunc {
	return func(ctx context.Context, md *metadata.MD) context.Context {
		key, val = EncodeKeyValue(key, val)
		(*md)[key] = append((*md)[key], val)
		return ctx
	}
}

// EncodeKeyValue sanitizes a key-value pair for use in gRPC metadata headers.
func EncodeKeyValue(key, val string) (string, string) {
	key = strings.ToLower(key)
	if strings.HasSuffix(key, binHdrSuffix) {
		val = base64.StdEncoding.EncodeToString([]byte(val))
	}
	return key, val
}
