package endpoint

import (
	"context"

	"cloud.google.com/go/functions/metadata"
	"github.com/huynguyen-quoc/go/common/constant"
)

type Meta struct {
	Name        string
	Type        string
	Client      string
	Service     string
	CircuitName string
	Idempotency constant.IdempotencyLevel
}

// Endpoint ...
type Endpoint func(ctx context.Context, request interface{}) (metadata metadata.Metadata, response interface{}, err error)
