package sarama

import (
	"context"
	"fmt"
	"github.com/huynguyen-quoc/go/streams/core"
	"github.com/huynguyen-quoc/go/streams/sarama"
	"sync"
)

type kafkaConsumer struct {
	sarama.Consumer
	dataChan     chan core.ConsumerMessage
	wg           sync.WaitGroup
	shutdownOnce sync.Once
	shutdownChan chan struct{}
}

func (k *kafkaConsumer) Start() {
	err := k.Consumer.Start(context.Background())
	if err != nil {
		fmt.Printf("Failed to start kafka consumer. Error: [%v]\n", err)
	}
}

func (k *kafkaConsumer) Shutdown() {
	k.shutdownOnce.Do(func() {
		defer close(k.shutdownChan)
		err := k.Consumer.Stop(context.Background())
		if err != nil {
			fmt.Printf("Failed to shutdown the kafka consumer. Error: [%v]\n", err)
		}
		if k.dataChan != nil {
			k.wg.Wait()
			close(k.dataChan)
		}
	})
}

func (k *kafkaConsumer) GetDataChan() <-chan core.ConsumerMessage {
	readMsgChan, err := k.Consumer.ReadMessage(context.Background())
	if err != nil {
		fmt.Printf("Error when read channel [%v]\n", err)
	}
	k.wg.Add(1)
	go func() {
		defer k.wg.Done()
		for msg := range readMsgChan {
			k.dataChan <- k.convert(msg)
		}
	}()
	return k.dataChan
}

func (k *kafkaConsumer) convert(msg *sarama.Message) core.ConsumerMessage {
	return &kafkaMessage{
		msg: msg,
	}
}

func (k *kafkaConsumer) Done() <-chan struct{} {
	return k.shutdownChan
}

func getConsumerGroupID(clientID string, stream string, consumerGroupID string) string {
	id := fmt.Sprintf("%s-%s", clientID, stream)

	if len(consumerGroupID) > 0 {
		id = consumerGroupID
	}

	return id
}
