package metadata

import (
	"context"
)

type ctxKey string

// metaCtxKey is a custom type for the context key
type metaCtxKey int

const (
	mdIncKey metaCtxKey = iota
	mdOutKey
)

func newContext(ctx context.Context, ctxKey metaCtxKey, meta Metadata) context.Context {
	return context.WithValue(ctx, ctxKey, meta)
}

// NewIncomingContext ...
func NewIncomingContext(ctx context.Context, meta Metadata) context.Context {
	return newContext(ctx, mdIncKey, meta)
}

// NewOutgoingContext ...
func NewOutgoingContext(ctx context.Context, meta Metadata) context.Context {
	return newContext(ctx, mdOutKey, meta)
}

func fromContext(ctx context.Context, ctxKey interface{}) Metadata {
	if meta, ok := ctx.Value(ctxKey).(Metadata); ok {
		return meta
	}
	return Metadata{}
}

// FromIncomingContext ...
func FromIncomingContext(ctx context.Context) Metadata {
	return fromContext(ctx, mdIncKey)
}

// FromOutgoingContext ...
func FromOutgoingContext(ctx context.Context) Metadata {
	return fromContext(ctx, mdOutKey)
}
