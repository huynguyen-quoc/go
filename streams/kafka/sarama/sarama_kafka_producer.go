package sarama

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/huynguyen-quoc/go/common/util"
	"github.com/huynguyen-quoc/go/streams/core"
	"github.com/huynguyen-quoc/go/streams/kafka/config"
	"github.com/huynguyen-quoc/go/streams/sarama"
	"sync"
	"time"
)

const (
	shutdownTimeout = 5 * time.Second
)

type kafkaProducer struct {
	sarama.Producer
	topic string
	wg    sync.WaitGroup
}

type saramaKafkaProducer struct {
	kafkaConfig *config.KafkaConfig
	stream      string
}

func (s saramaKafkaProducer) newSaramaKafkaProducer(ctx context.Context) (*kafkaProducer, error) {
	setup := &sarama.ProducerSetup{
		Config: &sarama.ProducerConfig{
			Brokers:          s.kafkaConfig.Brokers,
			ClientID:         s.stream,
			Sync:             s.kafkaConfig.EnableSync,
			ClusterType:      s.kafkaConfig.ClusterType,
			KafkaVersion:     s.kafkaConfig.KafkaVersion,
			CompressionCodec: s.kafkaConfig.CompressionCodec,
			CompressionLevel: s.kafkaConfig.CompressionLevel,
			RequiredAcks:     s.kafkaConfig.RequiredAcks,
		},
	}
	saramaProducer, err := setup.NewProducer(ctx)

	if err != nil {
		fmt.Printf("Setup SKafka Producer failed, error=[%v]\n", err)
		return nil, err
	}
	err = saramaProducer.Start(ctx)
	if err != nil {
		fmt.Printf("Starting  kafka producer failed; Error: [%v]\n", err)
		return nil, err
	}
	p := &kafkaProducer{
		Producer: saramaProducer,
		topic: s.stream,
	}
	return p, nil
}

// Add gets the record and writes it to the appropriate topic
func (producer *kafkaProducer) add(msg core.Message) error {
	bytes, err := proto.Marshal(msg)
	if err != nil {
		fmt.Printf("Failed to marshal the message with error: [%v]\n", err)
		return err
	}
	partitionBytes := proto.EncodeVarint(uint64(msg.GetStreamInfo().StreamPartitionID))
	return producer.writeBytesToStream(context.Background(), partitionBytes, bytes)
}

func (producer *kafkaProducer) writeBytesToStream(ctx context.Context, partitionBytes, dataBytes []byte) error {
	err := producer.Write(ctx, producer.topic, partitionBytes, dataBytes)
	if err != nil {
		fmt.Printf("Failed to write to the topic=[%s] with error=[%v]\n", producer.topic, err)
		return err
	}
	return nil
}

func (producer *kafkaProducer) shutdown() {
	err := producer.Stop(context.Background())
	if err != nil {
		fmt.Printf("Failed to shutdown the kafka producer with error=[%v]\n", err)
	}

	err = util.Wait(&producer.wg, shutdownTimeout)
	if err != nil {
		fmt.Printf("Failed to shutdown with time out error=[%v]\n", err)
	}
}
