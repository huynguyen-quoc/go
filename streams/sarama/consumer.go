package sarama

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

const (
	sleepDurationWhenErr = 500 * time.Millisecond
)

// Message encompasses the sarama message object
type Message struct {
	*sarama.ConsumerMessage
}

type consumerState uint32

const (
	notStarted consumerState = iota
	started
	// paused state is only intermediate, not supposed to be seen from outside
	paused
	// stopped state is final
	stopped
)

type readType uint32

const (
	readNotStarted readType = iota
	readMsg
)

type partitionOffsetInfo struct {
	sync.Mutex
	ackOffsets      []int64
	committedOffset int64
	skippedOffset   int64
}

type consumer struct {
	// STATE FLAG RELATED
	mtx sync.RWMutex
	//running  bool
	//closed   bool
	state consumerState

	readType readType

	// CONFIG RELATED
	topic  string
	config *ConsumerConfig

	// UNDERLYING SDK RELATED
	sdkConsumerGroup sarama.ConsumerGroup
	client           sarama.Client

	// OUTPUT RELATED
	msgChan chan *Message

	// LIFE CYCLE RELATED
	// used for breaking infinite Consume loop. Should be closed before calling ConsumerGroup.Close()
	stopReadChan chan struct{}
	shutdownChan chan struct{}

	// job routines controller for all reading logic
	readWg *sync.WaitGroup
}

func (consumer *consumer) Start(parent context.Context) error {
	consumer.mtx.Lock()
	defer consumer.mtx.Unlock()

	if consumer.state != notStarted {
		return errors.New("consumer is already started")
	}

	consumer.state = started
	consumer.startEventLoops()
	return nil
}

func (consumer *consumer) Stop(ctx context.Context) error {
	panic("implement me")
}

func (consumer *consumer) ReadMessage(ctx context.Context) (<-chan *Message, error) {
	var messageChan chan *Message
	err := errors.New("msgChan channel not initialized")
	consumer.mtx.Lock()
	defer consumer.mtx.Unlock()
	if consumer.readType == readNotStarted {
		messageChan = consumer.msgChan
		consumer.doRead(readMsg)
		err = nil
	}

	return messageChan, err
}

func (consumer *consumer) doRead(readT readType) {
	if readT == readNotStarted {
		return
	}
	consumer.readType = readT

	readFunc := func(m *sarama.ConsumerMessage, session sarama.ConsumerGroupSession) {
		consumer.readMsg(m, session)
	}

	consumer.readWg.Add(1)
	go func() {
		defer consumer.readWg.Done()

		for {
			select {
			case <-consumer.stopReadChan:
				return
			default:
				groupHandler := &ConsumerGroupReadHandler{
					readFunc:     readFunc,
					saramaClient: consumer.client,
					topic:        consumer.topic,

					// default sarama config commit interval
					commitInterval: time.Second,
					clientID:       consumer.config.ClientID,
					readWg:         consumer.readWg,
				}
				err := consumer.sdkConsumerGroup.Consume(context.Background(), []string{consumer.topic}, groupHandler)
				if err != nil {
					log.Printf("failed to consume message with error [%v]\n", err)
				}
			}
		}
	}()
}

// startEventLoops will start event listeners
func (consumer *consumer) startEventLoops() {

	consumer.readWg.Add(1)
	go func() {
		defer consumer.readWg.Done()
		for err := range consumer.sdkConsumerGroup.Errors() {
			log.Printf("consumer group failed to consume message. Error: [%v]\n", err)
			time.Sleep(sleepDurationWhenErr)
		}
	}()

}

func (consumer *consumer) readMsg(message *sarama.ConsumerMessage, session sarama.ConsumerGroupSession) {
	select {
	case consumer.msgChan <- consumer.toMessage(message):
	case <-consumer.stopReadChan:
		return
	}
	session.MarkMessage(message, "")
}

func (consumer *consumer) toMessage(message *sarama.ConsumerMessage) *Message {
	return &Message{message}
}

func (consumer *consumer) createSaramaComponents(config *ConsumerConfig) error {
	consumer.topic = config.Topic
	saramaClient, err := getSaramaClient(config)
	if err != nil {
		return err
	}
	consumer.client = saramaClient

	consumerGroup, err := getSDKConsumerGroup(config, saramaClient)
	if err != nil {
		return err
	}
	consumer.sdkConsumerGroup = consumerGroup

	return nil
}

func getSDKConsumerGroup(config *ConsumerConfig, client sarama.Client) (sarama.ConsumerGroup, error) {
	return sarama.NewConsumerGroupFromClient(config.ConsumerGroupID, client)
}

func getSaramaClient(config *ConsumerConfig) (sarama.Client, error) {
	return sarama.NewClient(config.Brokers, getSaramaConfig(config))
}

func getSaramaConfig(config *ConsumerConfig) *sarama.Config {
	consumerGroupConfig := sarama.NewConfig()
	consumerGroupConfig.Consumer.Return.Errors = true
	consumerGroupConfig.ClientID = config.ClientID
	consumerGroupConfig.Version = sarama.V2_6_0_0

	consumerGroupRebalanceStrategy := config.ConsumerGroupRebalanceStrategy.Strategy()
	// sarama.NewConfig() uses BalanceStrategyRange by default, we will not change the default if it is not set in ConsumerConfig
	if consumerGroupRebalanceStrategy != nil {
		consumerGroupConfig.Consumer.Group.Rebalance.Strategy = consumerGroupRebalanceStrategy
	}

	if config.InitOffset == OffsetOldest {
		consumerGroupConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
		return consumerGroupConfig
	}
	// by default start from latest
	consumerGroupConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	return consumerGroupConfig
}
