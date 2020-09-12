package kafkawriter

import (
	"github.com/huynguyen-quoc/go/streams/core"
	"github.com/huynguyen-quoc/go/streams/kafka"
)

type writerImpl struct {
	producer core.StreamProducer
}

func (w writerImpl) Shutdown() {
	w.producer.Shutdown()
}

func (w writerImpl) Done() <-chan struct{} {
	return w.producer.Done()
}

func (w writerImpl) Save(entity kafka.Entity) error {
	return w.producer.SaveToStream(entity)
}

func (w writerImpl) SaveBytesToStream(partitionBytes, dataBytes []byte) error {
	return w.producer.SaveBytesToStream(partitionBytes, dataBytes)
}

