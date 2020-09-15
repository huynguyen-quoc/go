package echo

import "golang.org/x/net/context"

// Service ...
type Service interface {

	// Echo
	Echo(ctx context.Context, value string) (string, error)
}
