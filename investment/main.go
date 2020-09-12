package main

import (
	"context"
	json2 "encoding/json"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/huynguyen-quoc/go/streams/kafka/config"
	"github.com/huynguyen-quoc/go/streams/kafka/kafkareader"
	"github.com/huynguyen-quoc/go/streams/kafka/kafkawriter"
	"github.com/huynguyen-quoc/go/streams/kafka/sarama"
	"github.com/huynguyen-quoc/go/streams/schema/ivstream"
	"os"
	"os/signal"
	"time"
)

var (
	brokers = []string{"localhost:9092"}
	topic   = "example-test"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	doneCh := make(chan struct{})
	go runProducer()
	go runConsumer()
	go func() {
		select {
		case <-signals:
			doneCh <- struct{}{}
		}
	}()
	<-doneCh


}

func runConsumer() {

	configData := &Config{
		Kafka: &config.KafkaConfig{
			Brokers:  brokers,
			ClientID: "12345",
			Stream:   topic,
			KafkaVersion: "2.6.0",

		},
	}

	var reader kafkareader.Client
	var err error
	kafkaReaderInit := &kafkareader.StreamSetup{
		Entity:     &ivstream.InvestmentOpenEntity{},
		Configurer: configData,
	}

	reader, err = kafkaReaderInit.NewReader(context.Background())
	if err != nil {
		fmt.Printf("error for init kafka [%v]\n", err)
		return
	}
	doneCh := reader.GetDataChan()
	go func() {
		for data := range doneCh {
			coreData := data.Event.(*ivstream.InvestmentOpenEntity)
			json, _ := json2.Marshal(coreData)
			fmt.Printf("result is [%v]\n", string(json))
		}
	}()


}

// Emit messages forever every second
func runProducer() {
	configData := &Config{
		Kafka: &config.KafkaConfig{
			Brokers: brokers,
			ClientID: "12345",
			Stream:  topic,
			KafkaVersion: "2.6.0",
			EnableSync: false,
		},
	}
	kafkaWriterInit := &kafkawriter.StreamSetup{
		Entity:     &ivstream.InvestmentOpenEntity{},
		Configurer: configData,
		KafkaInit: sarama.KafkaProducer{},
	}
	client, err := kafkaWriterInit.NewWriter(context.Background())
	if err != nil {
		fmt.Printf("error newWriter [%v]\n", err)
		return
	}
	for {
		time.Sleep(1 * time.Second)
		rand, _ := uuid.GenerateUUID()
		entity := &ivstream.InvestmentOpenEntity{
			Message: rand,
		}
		err = client.Save(entity)
		if err != nil {
			fmt.Printf("error send [%v]\n", err)
		}
		fmt.Printf("success send [%v]\n", rand)
	}
}
