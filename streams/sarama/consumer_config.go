package sarama

import (
	"context"
	"github.com/Shopify/sarama"
	"sync"
)

// OffsetPosition defines a type for Kafka topic offset
type OffsetPosition string

const (
	// OffsetOldest marks the oldest offset in Kafka topic
	OffsetOldest OffsetPosition = "oldest"

	// OffsetNewest marks the newest offset in Kafka topic
	OffsetNewest OffsetPosition = "newest"

	// RangeBalanceStrategyName identifies strategies that use the range partition assignment strategy
	RangeBalanceStrategyName ConsumerGroupRebalanceStrategy = sarama.RangeBalanceStrategyName
	// RoundRobinBalanceStrategyName identifies strategies that use the round-robin partition assignment strategy
	RoundRobinBalanceStrategyName = sarama.RoundRobinBalanceStrategyName
	// StickyBalanceStrategyName identifies strategies that use the sticky-partition assignment strategy
	StickyBalanceStrategyName = sarama.StickyBalanceStrategyName
)

// ConsumerGroupRebalanceStrategy is the config for sarama.Config.Consumer.Group.Rebalance.Strategy
type ConsumerGroupRebalanceStrategy string

// Strategy is to get the sarama.BalanceStrategy corresponding the the string value
func (s ConsumerGroupRebalanceStrategy) Strategy() sarama.BalanceStrategy {
	switch s {
	case RangeBalanceStrategyName:
		return sarama.BalanceStrategyRange
	case RoundRobinBalanceStrategyName:
		return sarama.BalanceStrategyRoundRobin
	case StickyBalanceStrategyName:
		return sarama.BalanceStrategySticky
	default:
		return nil
	}
}

// ConsumerConfig defines a struct for Kafka consumer configurations
type ConsumerConfig struct {
	Brokers                        []string                       `json:"brokers"`
	ConsumerGroupID                string                         `json:"consumerGroupID"` // Unique identifier for a consumer group
	Topic                          string                         `json:"topic"`           // We are able to consume multiple topics with one consumer, define 1 topic for simplicity here
	ClientID                       string                         `json:"clientID"`
	InitOffset                     OffsetPosition                 `json:"initOffset"`
	ClusterType                    string                         `json:"clusterType"`
	ConsumerGroupRebalanceStrategy ConsumerGroupRebalanceStrategy `json:"consumerGroupRebalanceStrategy"` // default to sarama.BalanceStrategySticky when empty
}

type ConsumerSetup struct {
	Config *ConsumerConfig
}

func (s *ConsumerSetup) NewConsumer(ctx context.Context) (Consumer, error) {
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


