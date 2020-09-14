package protocol

import (
	"github.com/huynguyen-quoc/go/investment/config"
	"golang.org/x/net/context"
)

type Runner interface {

	Run(ctx context.Context, config *config.AppConfig) error
}
