package kafka

import (
	"context"
	"github.com/huynguyen-quoc/go/streams/core"
	"github.com/huynguyen-quoc/go/streams/kafka/config"
)

//go:generate mockery --inpackage --case underscore  --name ConsumerInitialization
type ConsumerInitialization interface {
	NewKafkaConsumer(ctx context.Context) (core.StreamConsumer, error)
}

//go:generate mockery --inpackage --case underscore  --name ProducerInitialization
type ProducerInitialization interface {
	NewKafkaProducer(ctx context.Context, kafkaConfig config.KafkaConfig, streamID string) (core.StreamProducer, error)
}

