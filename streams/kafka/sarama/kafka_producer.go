package sarama

import (
	"context"
	"log"

	"github.com/huynguyen-quoc/go/streams/core"
	"github.com/huynguyen-quoc/go/streams/kafka/config"
)

type KafkaProducer struct {
}

func (k KafkaProducer) NewKafkaProducer(ctx context.Context, kafkaConfig config.KafkaConfig, streamID string) (core.StreamProducer, error) {

	p := &producer{
		runningCh: make(chan struct{}),
	}

	// if we have a valid kafkaConfig, use kafka
	// Note: all the feature flags checks are performed by the caller
	if len(kafkaConfig.Brokers) > 0 {

		log.Printf("Initialize kafka producer for stream=[%s] topicName=[%s]\n", streamID, kafkaConfig.Stream)
		log.Printf("Using Kafka to write streamName=[%s] brokers=[%+v]\n", kafkaConfig.Stream, kafkaConfig.Brokers)

		saramaKafkaProducer := &saramaKafkaProducer{
			kafkaConfig: &kafkaConfig,
			stream:      kafkaConfig.Stream,
		}
		kafkaProducer, err := saramaKafkaProducer.newSaramaKafkaProducer(ctx)
		if err != nil {
			return nil, err
		}

		p.kafkaProducer = kafkaProducer
	}

	return p, nil
}
