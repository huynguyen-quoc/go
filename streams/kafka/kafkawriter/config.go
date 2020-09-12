package kafkawriter

import (
	"context"
	"fmt"
	"github.com/huynguyen-quoc/go/streams/kafka"
)

// StreamConfig defines the interface required to init the connection to StreamConfig server
type StreamConfig struct {
	Configurer kafka.Configurer
}

type StreamSetup struct {
	Entity     kafka.Entity
	Configurer kafka.Configurer
	KafkaInit  kafka.ProducerInitialization
}

func (s StreamSetup) NewWriter(ctx context.Context) (Client, error) {
	cfg := &StreamConfig{
		Configurer: s.Configurer,
	}

	kafkaConfig := cfg.Configurer.GetConfig()

	producer, err := s.KafkaInit.NewKafkaProducer(ctx, kafkaConfig, cfg.Configurer.GetStreamID())
	if err != nil {
		fmt.Printf("failed to init the producer for streamID=[%s] err=[%+v]\n", cfg.Configurer.GetStreamID(), err)
		return nil, err
	}

	return &writerImpl{
		producer: producer,
	}, nil
}
