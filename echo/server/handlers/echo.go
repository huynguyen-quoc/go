package handlers

import "context"

// EchoService implements echo.Service
type EchoService struct{}

//  New returns a new Service implementation
func New() *EchoService {
	// This function should take your service's dependencies as parameters and store them in the Service object for use by
	// the handlers. This can be used for dependency injection.
	return &EchoService{
		// TODO: Add dependencies here
	}
}

func (e EchoService) Echo(_ context.Context, value string) (string, error) {
	return value, nil
}
