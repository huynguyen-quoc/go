package sarama

import (
	"context"
	"log"

	"github.com/huynguyen-quoc/go/streams/core"
	"github.com/huynguyen-quoc/go/streams/kafka/config"
	"github.com/huynguyen-quoc/go/streams/sarama"
)

type KafkaConsumer struct{}

func (k KafkaConsumer) NewKafkaConsumer(ctx context.Context, kafkaConfig config.KafkaConfig, streamID string) (core.StreamConsumer, error) {
	consumerGroupID := getConsumerGroupID(kafkaConfig.ClientID, streamID, kafkaConfig.ConsumerGroupID)

	consumerConfig := &sarama.ConsumerConfig{
		Brokers:         kafkaConfig.Brokers,
		Topic:           streamID,
		ClientID:        kafkaConfig.ClientID,
		ConsumerGroupID: consumerGroupID,
		InitOffset:      getKafkaOffsetPosition(kafkaConfig.OffsetType),
		ClusterType:     kafkaConfig.ClusterType,
	}

	consumerSetup := &sarama.ConsumerSetup{
		Config: consumerConfig,
	}

	c, err := consumerSetup.NewConsumer(ctx)
	if err != nil {
		log.Printf("instantiating kafka consumer failed; Error: [%v]\n", err)
		return nil, err
	}

	consumer := &kafkaConsumer{
		Consumer:     c,
		shutdownChan: make(chan struct{}),
		dataChan:     make(chan core.ConsumerMessage),
	}

	return consumer, nil
}
