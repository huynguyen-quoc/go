package sarama

import (
	"context"
	"fmt"
	"github.com/huynguyen-quoc/go/streams/core"
	"github.com/huynguyen-quoc/go/streams/kafka/config"
	"github.com/huynguyen-quoc/go/streams/sarama"
)

type KafkaConsumer struct {
	KafkaConfig config.KafkaConfig
	StreamID string
}

func (k KafkaConsumer) NewKafkaConsumer(ctx context.Context) (core.StreamConsumer, error) {
	consumerGroupID := getConsumerGroupID(k.KafkaConfig.ClientID, k.StreamID, k.KafkaConfig.ConsumerGroupID)

	consumerConfig := &sarama.ConsumerConfig{
		Brokers:         k.KafkaConfig.Brokers,
		Topic:           k.StreamID,
		ClientID:        k.KafkaConfig.ClientID,
		ConsumerGroupID: consumerGroupID,
		InitOffset:      getKafkaOffsetPosition(k.KafkaConfig.OffsetType),
		ClusterType:     k.KafkaConfig.ClusterType,
	}

	consumerSetup := &sarama.ConsumerSetup{
		Config: consumerConfig,
	}

	c, err := consumerSetup.NewConsumer(ctx)
	if err != nil {
		fmt.Printf("instantiating kafka consumer failed; Error: [%v]\n", err)
		return nil, err
	}

	consumer := &kafkaConsumer{
		Consumer:     c,
		shutdownChan: make(chan struct{}),
	}

	return consumer, nil
}
