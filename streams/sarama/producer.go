package sarama

import (
	"context"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"strings"
	"sync"
	"time"
)

const (
	defaultProducerFlushTimeout = 500 * time.Millisecond

	defaultProducerWriteTimeout = 3 * time.Second

	defaultKafkaVersion = "0.10.2.1"

	defaultCompressionCodec = sarama.CompressionSnappy
)

// syncProducer is the kafka specific sync producer which implements the
// Producer interface
type syncProducer struct {
	mtx      sync.RWMutex
	shutdown bool

	sdkProducer sarama.SyncProducer
}

// newSyncProducer instantiates the sync producer
func newSyncProducer(config *ProducerConfig) (*syncProducer, error) {
	sdkSyncProducer, err := sarama.NewSyncProducer(config.Brokers, getSaramaProducerConfig(config))
	if err != nil {
		fmt.Printf("Failed to instantiate sarama sync producer with error=[%v]\n", err)
		return nil, err
	}

	p := &syncProducer{
		sdkProducer: sdkSyncProducer,
	}

	return p, nil
}

func (s *syncProducer) Start(ctx context.Context) error {
	return nil
}

func (s *syncProducer) Stop(ctx context.Context) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.shutdown {
		return errors.New("producer already stop")
	}

	s.shutdown = true
	return nil
}

func (s *syncProducer) Write(ctx context.Context, topic string, partitionKey []byte, value []byte) error {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	// if we are already shutdown, return immediately
	if s.shutdown {
		return errors.New("producer already stop")
	}

	pMsg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(partitionKey),
		Value: sarama.ByteEncoder(value),
	}
	_, _, err := s.sdkProducer.SendMessage(pMsg)
	if err != nil {
		fmt.Printf("Failed to produce message on topic=[%s]. Error=[%s]\n", topic, err)
		return err
	}
	return nil
}

type asyncProducer struct {
	mtx                sync.RWMutex
	wg                 *sync.WaitGroup
	running            bool
	lastFlushTimeNanos int64
	writeTimeout       time.Duration
	sdkProducer        sarama.AsyncProducer
	closed             bool
}

// This is used to pass metadata to the sarama writer
type producerMessageMeta struct {
	writeTime time.Time
}

func newAsyncProducer(config *ProducerConfig) (*asyncProducer, error) {
	sdkProducer, err := getAsyncProducer(config)
	if err != nil {
		return nil, err
	}

	p := &asyncProducer{
		wg:           &sync.WaitGroup{},
		sdkProducer:  sdkProducer,
		writeTimeout: defaultProducerWriteTimeout,
	}
	if config.WriteTimeoutInMillis > int64(0) {
		p.writeTimeout = time.Duration(config.WriteTimeoutInMillis) * time.Millisecond
	}

	return p, nil
}

func (producer *asyncProducer) Start(ctx context.Context) error {
	producer.mtx.Lock()
	defer producer.mtx.Unlock()

	if producer.running {
		return errors.New("producer already started")
	}

	producer.running = true
	producer.startEventLoops()
	return nil
}

func (producer *asyncProducer) Stop(ctx context.Context) error {
	panic("implement me")
}

func (producer *asyncProducer) Write(ctx context.Context, topic string, partitionKey []byte, value []byte) error {
	producer.mtx.RLock()
	defer producer.mtx.RUnlock()

	if !producer.running {
		return errors.New("producer already stop")
	}
	now := time.Now()
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(partitionKey),
		Value: sarama.ByteEncoder(value),
		Metadata: &producerMessageMeta{
			writeTime: now,
		},
	}
	ctx, cancelFn := context.WithTimeout(ctx, producer.writeTimeout)
	defer cancelFn()
	select {
	case producer.sdkProducer.Input() <- msg:
		break
	case <-ctx.Done():
		fmt.Printf("timed out sending message to kafka for topic=[%s], partition key=[%s]\n", topic, string(partitionKey))
		return errors.New("producer already stop")
	}
	return nil
}

func (producer *asyncProducer) stopInternal(parent context.Context, fullShutdown bool) error {
	if producer.closed {
		fmt.Println("producer already closed")
		return nil
	}
	producer.closed = true
	producer.sdkProducer.AsyncClose()

	shutdownChan := make(chan struct{})
	go func() {
		producer.wg.Wait()
		close(shutdownChan)
	}()

	ctx, cancelFunc := context.WithTimeout(parent, defaultProducerFlushTimeout)
	defer cancelFunc()

	select {
	case <-shutdownChan:
		fmt.Println( "Successfully stopped kafka producer")
	case <-ctx.Done():
		fmt.Printf("Close kafka producer error=[%v]\n", ctx.Err())
		return ctx.Err()
	}

	return nil
}

func (producer *asyncProducer) startEventLoops() {
	sdkProducer := producer.sdkProducer

	producer.wg.Add(1)

	go func() {
		defer producer.wg.Done()

		for err := range sdkProducer.Errors() {
			// even errors should be calculated as attempted add
			fmt.Printf("Failed to produce message. error=[%v]\n", err)
		}
	}()
}

func getAsyncProducer(config *ProducerConfig) (sarama.AsyncProducer, error) {
	sdkProducer, err := sarama.NewAsyncProducer(config.Brokers, getSaramaProducerConfig(config))
	if err != nil {
		fmt.Printf("Failed to instantiate sarama producer with error=[%v]\n", err)
	}

	return sdkProducer, err
}

func getSaramaProducerConfig(config *ProducerConfig) *sarama.Config {
	// Configs common for all producers
	producerConfig := sarama.NewConfig()
	producerConfig.ClientID = config.ClientID
	producerConfig.Version = getKafkaVersion(config.KafkaVersion)

	if config.Sync {
		producerConfig.Producer.Return.Successes = true
	}

	producerConfig.Producer.RequiredAcks = getRequiredACKs(config.RequiredAcks)
	producerConfig.Producer.Compression = getCompressionCodec(config.CompressionCodec)
	producerConfig.Producer.CompressionLevel = getCompressionLevel(config.CompressionLevel)

	fmt.Printf("ClientID=[%s], Kafka version=[%v]\n", producerConfig.ClientID, producerConfig.Version)
	fmt.Printf("Compression codec=[%s], Compression level=[%d]\n", producerConfig.Producer.Compression, producerConfig.Producer.CompressionLevel)

	if config.Sync {
		// Sync producer specific configs should go below
		producerConfig.Producer.Return.Errors = true
		if config.RequiredAcks != 0 {
			producerConfig.Producer.RequiredAcks = sarama.RequiredAcks(config.RequiredAcks)
		}
		return producerConfig
	}

	// Async producer specific configs should go below

	return producerConfig
}

func getKafkaVersion(kafkaVersion string) sarama.KafkaVersion {
	version, err := sarama.ParseKafkaVersion(kafkaVersion)
	if err != nil {
		fmt.Printf("Failed to parse sarama version with error=[%v], use default one instead\n", err)
		version, _ = sarama.ParseKafkaVersion(defaultKafkaVersion)
	}
	return version
}

func getRequiredACKs(requireAcks int16) sarama.RequiredAcks {

	switch requireAcks {
	case -1:
		return sarama.WaitForAll
	case 0:
		return sarama.NoResponse
	case 1:
		return sarama.WaitForLocal
	default:
		fmt.Printf("Unknown require acks type=[%d], use default type=[%v]", requireAcks, sarama.WaitForAll)
		return sarama.WaitForAll
	}
}

func getCompressionCodec(compressionCodec string) sarama.CompressionCodec {
	compressionCodecLower := strings.ToLower(compressionCodec)

	switch compressionCodecLower {
	case "gzip":
		return sarama.CompressionGZIP
	case "lz4":
		return sarama.CompressionLZ4
	case "snappy":
		return sarama.CompressionSnappy
	case "none":
		return sarama.CompressionNone
	case "":
		return defaultCompressionCodec
	default:
		fmt.Printf("Unknown compression codec used: [%s]. Using default codec: [%s]\n", compressionCodec, defaultCompressionCodec)
		return defaultCompressionCodec
	}
}

func getCompressionLevel(compressionLevel int) int {
	if compressionLevel == 0 {
		return sarama.CompressionLevelDefault
	}

	return compressionLevel
}
