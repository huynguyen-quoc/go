package sarama

import (
	"context"
)

// Starter - defines common behaviors of Apache Kafka entities
type Starter interface {
	// Start initializes and starts
	Start(ctx context.Context) error
}

// Stopper - defines stop behaviors of Apache Kafka
type Stopper interface {

	// Stop stops the consumer
	Stop(ctx context.Context) error
}

// Producer - defines behaviors of Apache Kafka producer
type Producer interface {
	Starter
	Stopper

	Write(ctx context.Context, topic string, partitionKey []byte, value []byte) error
}

// Consumer - defines behaviors of Apache Kafka consumer
type Consumer interface {
	Starter
	Stopper

	// ReadMessage reads the message
	ReadMessage(ctx context.Context) (<-chan *Message, error)
}

type ConsumerInitialization interface {

	NewConsumer(ctx context.Context) (Consumer, error)
}

type ProducerInitialization interface {

	NewProducer(ctx context.Context) (Producer, error)
}

