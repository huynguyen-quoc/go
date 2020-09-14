package kafkareader

import (
	"context"
	"fmt"
	"github.com/huynguyen-quoc/go/streams/kafka"
)


type ReaderInit struct {
	Entity     kafka.Entity
	Configurer kafka.Configurer
	KafkaInit  kafka.ConsumerInitialization
}

func (s ReaderInit) NewReader(ctx context.Context) (Client, error) {
	cfg := &kafka.StreamConfig{
		Configurer: s.Configurer,
	}

	kafkaConfig := cfg.Configurer.GetConfig()
	streamID := cfg.Configurer.GetStreamID()

	consumer, err := s.KafkaInit.NewKafkaConsumer(ctx, kafkaConfig, streamID)
	if err != nil {
		fmt.Printf("failed to init the consumer for streamID=[%s] err=[%+v]\n", streamID, err)
		return nil, err
	}

	client, _ := initialize(consumer, s.Entity)

	return client, nil
}

