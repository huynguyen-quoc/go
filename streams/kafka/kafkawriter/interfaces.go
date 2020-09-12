package kafkawriter

import (
	"context"
	"github.com/huynguyen-quoc/go/streams/kafka"
)

// Client defines the methods for interacting with streams writer
//go:generate mockery --inpackage --case underscore --name Client
type Client interface {
	// Shutdown will complete all pending writes and then stop the writer.
	// Calls to shutdown an already closed write will be a NOP
	Shutdown()

	// Done will return a channel that returns when the shutdown is complete
	Done() <-chan struct{}

	// Save will enqueue the supplied entity for sending to the stream as a new object.
	Save(entity kafka.Entity) error

	// SaveBytesToStream will enqueue the supplied bytes for sending to the stream as a new object
	SaveBytesToStream(partitionBytes, dataBytes []byte) error

}

//go:generate mockery --inpackage --case underscore --name WriterInitialization
type WriterInitialization interface {

	NewReader(ctx context.Context, configurer kafka.Configurer) (Client, error)
}