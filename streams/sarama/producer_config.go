package sarama

import (
	"context"
	"sync"
)

// ProducerConfig defines a struct for Kafka producer configuration
type ProducerConfig struct {
	// Cluster Level Config
	Brokers      []string   `json:"brokers"`
	Topic        string     `json:"topic"`
	ClientID     string     `json:"client_id"`
	KafkaVersion string     `json:"kafkaVersion"`
	StatsdTag    string     `json:"-"`

	// Producer Level Configurations
	EnableRetry      bool   `json:"enableRetry"`
	ClusterType      string `json:"clusterType"`
	CompressionCodec string `json:"compressionCodec"`
	CompressionLevel int    `json:"compressionLevel"`
	Sync             bool   `json:"sync"`
	RequiredAcks     int16  `json:"requiredAcks"` // number of acks required from replicas (passed in Produce Requests)

	// Defaults to 3 seconds
	WriteTimeoutInMillis int64 `json:"writeTimeoutInMillis"`
}

type ProducerSetup struct {
	Config *ProducerConfig
}

func (p ProducerSetup) NewProducer(ctx context.Context) (Producer, error) {
	if p.Config.Sync {
		return newSyncProducer(p.Config)
	}
	return newAsyncProducer(p.Config)
}

func (s *ConsumerSetup) I(ctx context.Context) (Consumer, error) {
	consumer := &consumer{
		shutdownChan:  make(chan struct{}),
		stopReadChan:  make(chan struct{}),
		msgChan:       make(chan *Message),

		config:      s.Config,
		readWg:      &sync.WaitGroup{},
	}

	err := consumer.createSaramaComponents(s.Config)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}


