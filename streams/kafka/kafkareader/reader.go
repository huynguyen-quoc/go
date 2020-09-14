package kafkareader

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/huynguyen-quoc/go/streams/core"
	"github.com/huynguyen-quoc/go/streams/kafka"
	"sync"
)

// Entity ...
type Entity struct {
	Event kafka.Entity
}

func initialize(consumer core.StreamConsumer, entity kafka.Entity) *readerImpl {
	return &readerImpl{
		consumer:     consumer,
		streamEntity: entity,
		shutdownChan: make(chan struct{}),
		dataCh:       make(chan *Entity),
	}
}

type readerImpl struct {
	streamEntity kafka.Entity
	consumer     core.StreamConsumer
	wg           sync.WaitGroup
	initOnce     sync.Once
	shutdownOnce sync.Once
	dataCh       chan *Entity
	shutdownChan chan struct{}
}

func (r *readerImpl) GetDataChan() <-chan *Entity {
	r.initOnce.Do(func() {
		r.consumer.Start()

		r.wg.Add(1)
		go r.process(r.consumer.GetDataChan(), r.processMessage)
	})

	return r.dataCh
}

func (r *readerImpl) Shutdown() error {
	r.shutdownOnce.Do(func() {
		r.consumer.Shutdown()

		waitForShutdown := make(chan struct{})
		go func() {
			r.wg.Wait()

			close(r.dataCh)
			close(waitForShutdown)
		}()

		message := r.streamEntity.GetMessage().String()

		select {
		case <-waitForShutdown:
			fmt.Printf("[stream:%s] Consumer shutdown gracefully and drained packets in channel", message)
		}
		close(r.shutdownChan)
	})

	return nil
}

func (r *readerImpl) Done() <-chan struct{} {
	return r.shutdownChan
}

func (r *readerImpl) process(messageChan <-chan core.ConsumerMessage, processEvent func(message core.ConsumerMessage)) {
	defer r.wg.Done()

	for message := range messageChan {
		processEvent(message)
	}
}

func (r *readerImpl) processMessage(message core.ConsumerMessage) {
	msg := r.streamEntity.GetMessage()
	err := proto.Unmarshal(message.Data(), msg)
	if err != nil {
		fmt.Printf("Error while unmarshalling data to message [%v]\n", err)
		return
	}
	data := r.streamEntity.FromPB(msg)
	select {
	case <-r.shutdownChan:
		// consumer shutdown
		return

	case r.dataCh <- &Entity{
		Event: data,
	}:
		//success
	}
}
