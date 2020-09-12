package kafkareader

import (
	"context"
	"github.com/huynguyen-quoc/go/streams/kafka"
)

//go:generate mockery --inpackage --case underscore --name Client
type Client interface {

	// GetDataChan will return a channel of entity (cancel the context to stop the flow of data)
	GetDataChan() <-chan *Entity

	// Shutdown will stop reading from the stream
	Shutdown() error

	// Done will return a channel that returns when the last msg from the stream is read during shutdown
	Done() <-chan struct{}

}

//go:generate mockery --inpackage --case underscore --name ReaderInitialization
type ReaderInitialization interface {

	NewReader(ctx context.Context, entity kafka.Entity, configurer kafka.Configurer) (Client, error)
}