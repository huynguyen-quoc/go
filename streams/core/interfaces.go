package core

import (
	"github.com/golang/protobuf/proto"
	"github.com/huynguyen-quoc/go/streams/schema/common"
	"time"
)

// WriterDTO defines the base on all stream DTOs.
//go:generate mockery --inpackage --case underscore  --name WriterDTO
type WriterDTO interface {
	// ToProtoBuf returns a pb.Message filled with the supplied data from the DTO
	ToProtoBuf(isUpdate bool) Message
}

// Message is the base type for all messages ( use for proto buf extensions)
type Message interface {
	proto.Message
	Descriptor() ([]byte, []int)
	GetStreamInfo() common.StreamInfoEntity
}

// StreamConsumer defines a general Consumer
//go:generate mockery --inpackage --case underscore  --name StreamConsumer
type StreamConsumer interface {
	// initialize the consumer and start consuming the stream
	Start()

	// Stop consuming from the stream and shutdown
	Shutdown()

	// GetDataChan gets ConsumerMessage output channel
	GetDataChan() <-chan ConsumerMessage

	// Done will return a channel that returns when the StreamConsumer is done/shutdown
	Done() <-chan struct{}
}

// ConsumerMessage defines the general message interface for the stream
//go:generate mockery --inpackage --case underscore  --name ConsumerMessage
type ConsumerMessage interface {
	// Data gets the payload of the message
	Data() []byte

	// Partition gets Kafka partition of the message
	Partition() int32

	// Offset gets offset of the message
	Offset() int64

	// Timestamp gets timestamp of the message
	Timestamp() time.Time
}

// StreamProducer defines a general producer
//go:generate mockery --inpackage --case underscore  --name StreamProducer
type StreamProducer interface {
	// SaveToStream will enqueue the supplied DTO for sending to the stream as a new object.
	SaveToStream(dto WriterDTO) error

	// SaveBytesToStream will enqueue the supplied bytes for sending to the stream as a new object
	SaveBytesToStream(partitionBytes, dataBytes []byte) error

	// Shutdown shuts down the producer
	Shutdown()

	// Done will return a channel that returns when the StreamConsumer is done/shutdown
	Done() <-chan struct{}
}
