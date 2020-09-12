package kafkareader

import (
	"context"
	"fmt"
	"github.com/huynguyen-quoc/go/streams/kafka"
	"github.com/huynguyen-quoc/go/streams/kafka/sarama"
)

// StreamConfig defines the interface required to init the connection to stream server
type StreamConfig struct {
	Configurer kafka.Configurer
}

type StreamSetup struct {
	Entity     kafka.Entity
	Configurer kafka.Configurer
}

func (s StreamSetup) NewReader(ctx context.Context) (Client, error) {
	cfg := &StreamConfig{
		Configurer: s.Configurer,
	}

	kafkaConfig := cfg.Configurer.GetConfig()

	kafkaConsumer := sarama.KafkaConsumer{
		KafkaConfig: kafkaConfig,
		StreamID: cfg.Configurer.GetStreamID(),
	}
	consumer, err := kafkaConsumer.NewKafkaConsumer(ctx)
	if err != nil {
		fmt.Printf("failed to init the consumer for streamID=[%s] err=[%+v]\n", cfg.Configurer.GetStreamID(), err)
		return nil, err
	}

	client, _ := initialize(consumer, s.Entity)

	return client, nil
}
