package grpc

import (
	"context"

	serverMetadata "github.com/huynguyen-quoc/go/server/metadata"
	"google.golang.org/grpc/metadata"
)

var metadataToGRPC = map[string]string{
	serverMetadata.KeyRequestAuthorization: "authorization",
	serverMetadata.KeyCacheTTL:             "cache-ttl",
	serverMetadata.KeyCacheControl:         "cache-control",
	serverMetadata.KeyClientID:             "client-id",
	serverMetadata.KeyClientIP:             "ip",
}

var headerToMetadataKey = map[string]string{}

func convertMeta(md metadata.MD) serverMetadata.Metadata {
	kitMeta := serverMetadata.Metadata(md)
	for k, v := range md {
		if newKey, ok := headerToMetadataKey[k]; ok {
			kitMeta[newKey] = v
		}
	}
	return kitMeta
}

func convertMetaCtx(ctx context.Context, md metadata.MD) context.Context {
	if len(md) == 0 {
		return ctx
	}
	meta := serverMetadata.Join(serverMetadata.FromIncomingContext(ctx), convertMeta(md))
	return serverMetadata.NewIncomingContext(ctx, meta)
}

// map grab-kit metadata to gRPC header names
func headersToGRPC(md serverMetadata.Metadata) metadata.MD {
	out := metadata.MD{}
	for k, v := range md {
		if newKey, ok := metadataToGRPC[k]; ok {
			out[newKey] = v
			continue
		}
		out[k] = v
	}
	return out
}

func init() {
	// generate the inverse map
	for k, v := range metadataToGRPC {
		headerToMetadataKey[v] = k
	}
}
