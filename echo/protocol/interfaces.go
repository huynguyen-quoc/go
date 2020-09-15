package protocol

import (
	"golang.org/x/net/context"
)

type Runner interface {

	Run(ctx context.Context) error
}
