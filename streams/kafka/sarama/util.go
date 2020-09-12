package sarama

import (
	"github.com/huynguyen-quoc/go/streams/kafka/config"
	"github.com/huynguyen-quoc/go/streams/sarama"
)

// getKafkaOffsetPosition returns the kafka.OffsetPosition given the config
func getKafkaOffsetPosition(offsetType config.OffsetType) sarama.OffsetPosition {
	switch offsetType {
	case config.OffsetNewest:
		return sarama.OffsetNewest
	case config.OffsetOldest:
		return sarama.OffsetOldest
	default: // in case of unsupported Kafka OffsetType, default will be OffsetNewest
		return sarama.OffsetNewest
	}
}
