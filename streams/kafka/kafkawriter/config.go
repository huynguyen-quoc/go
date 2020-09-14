package kafkawriter

import (
	"context"
	"fmt"
	"github.com/huynguyen-quoc/go/streams/kafka"
	"github.com/huynguyen-quoc/go/streams/kafka/sarama"
)


var (
	SaramaProducer kafka.ProducerInitialization = sarama.KafkaProducer{}
)

type WriterInit struct {
	Entity     kafka.Entity
	Configurer kafka.Configurer
	KafkaInit  kafka.ProducerInitialization
}

func (s WriterInit) NewWriter(ctx context.Context) (Client, error) {
	cfg := &kafka.StreamConfig{
		Configurer: s.Configurer,
	}

	kafkaConfig := cfg.Configurer.GetConfig()
	streamID := cfg.Configurer.GetStreamID()

	producer, err := s.KafkaInit.NewKafkaProducer(ctx, kafkaConfig, streamID)
	if err != nil {
		fmt.Printf("failed to init the producer for streamID=[%s] err=[%+v]\n", streamID, err)
		return nil, err
	}

	return &writerImpl{
		producer: producer,
	}, nil
}
