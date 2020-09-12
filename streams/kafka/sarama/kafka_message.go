package sarama

import (
	"github.com/huynguyen-quoc/go/streams/sarama"
	"time"
)

type kafkaMessage struct {
	msg *sarama.Message
}

// Data is the implementation of the corresponding method of ConsumerMessage interface
func (m *kafkaMessage) Data() []byte {
	return m.msg.Value
}

// Partition is the implementation of the corresponding method of kafkaMessage interface
func (m *kafkaMessage) Partition() int32 {
	return m.msg.Partition
}

// Offset is the implementation of the corresponding method of kafkaMessage interface
func (m *kafkaMessage) Offset() int64 {
	return m.msg.Offset
}

// Timestamp is the implementation of the corresponding method of kafkaMessage interface
func (m *kafkaMessage) Timestamp() time.Time {
	return m.msg.Timestamp
}
