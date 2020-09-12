package sarama

import (
	"context"
	"github.com/huynguyen-quoc/go/streams/core"
	"sync"
)


type producer struct {
	kafkaProducer  *kafkaProducer
	runningCh      chan struct{}
	shutdownOnce   sync.Once
}

func (p *producer) SaveToStream(dto core.WriterDTO) error {
	return p.addToStream(dto, true)
}

func (p *producer) UpdateToStream(dto core.WriterDTO) error {
	return p.addToStream(dto, false)
}

func (p *producer) SaveBytesToStream(partitionBytes, dataBytes []byte) error {
	return p.kafkaProducer.writeBytesToStream(context.Background(), partitionBytes, dataBytes)
}

func (p *producer) Shutdown() {
	p.shutdownOnce.Do(func() {
		defer close(p.runningCh)
		if p.kafkaProducer != nil {
			p.kafkaProducer.shutdown()
		}
	})
}

func (p *producer) Done() <-chan struct{} {
	return p.runningCh
}


func (p *producer) addToStream(dto core.WriterDTO, isUpdate bool) error {
	var err error
	msg := dto.ToProtoBuf(isUpdate)
	if p.kafkaProducer != nil {
		err = p.kafkaProducer.add(msg)
	}
	return err
}



