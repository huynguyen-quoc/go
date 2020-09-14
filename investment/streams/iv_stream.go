package streams

import (
	json2 "encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/huynguyen-quoc/go/investment/config"
	"github.com/huynguyen-quoc/go/streams/kafka/kafkareader"
	"github.com/huynguyen-quoc/go/streams/kafka/kafkawriter"
	"github.com/huynguyen-quoc/go/streams/schema/ivstream"
	"golang.org/x/net/context"
	"time"
)

type InvestmentStream struct {
	Config config.AppConfig
	writer kafkawriter.Client
	reader kafkareader.Client
}

func (i InvestmentStream) Start(ctx context.Context) {
	go i.newProducer()
	i.newConsumer(ctx)
}

func (i InvestmentStream) SendData(ctx context.Context, data interface{}) error {
	if i.writer == nil {
		fmt.Println("Investment Streams producer was not init")
		return errors.New("producer was not init")
	}

	entity, ok := data.(*ivstream.InvestmentOpenEntity)
	if !ok {
		fmt.Println("Data is not InvestmentOpen Entity")
		return errors.New("data is not investmentOpen entity")
	}

	result := i.writer.Save(entity)

	return result
}

func (i InvestmentStream) Stop(ctx context.Context) {
	panic("implement me")
}

func (i InvestmentStream) newConsumer(ctx context.Context) {

	var reader kafkareader.Client
	var err error
	kafkaReaderInit := &kafkareader.ReaderInit{
		Entity:     &ivstream.InvestmentOpenEntity{},
		Configurer: &i.Config,
		KafkaInit:  kafkareader.SaramaConsumer,
	}

	reader, err = kafkaReaderInit.NewReader(context.Background())
	if err != nil {
		fmt.Printf("error for init kafka consumer [%v]\n", err)
		return
	}
	i.reader = reader
	if i.reader == nil {
		fmt.Println("Investment Streams Consumer was not init")
		return
	}
	ch := i.reader.GetDataChan()
	i.consumeIvStream(ctx, ch)
}

func (i InvestmentStream) consumeIvStream(ctx context.Context, ch <-chan *kafkareader.Entity) {
	for data := range ch {
		entity, err := data.Event.(*ivstream.InvestmentOpenEntity)
		if !err {
			fmt.Printf("Wrong entity in ptReader, event=[%#v]\n", data.Event)
			return
		}
		errHandle := handlePTStream(ctx, *entity)
		if errHandle != nil {
			fmt.Printf("handling ptstream event failed, investmentId=[%d] err=[%v]", entity.InvestmentId, err)
		}
	}
	fmt.Printf("iv streams read channel drained, send signal to stop")
}

// Emit messages forever every second
func (i InvestmentStream) newProducer() {

	kafkaWriterInit := &kafkawriter.WriterInit{
		Entity:     &ivstream.InvestmentOpenEntity{},
		Configurer: &i.Config,
		KafkaInit:  kafkawriter.SaramaProducer,
	}
	client, err := kafkaWriterInit.NewWriter(context.Background())
	if err != nil {
		fmt.Printf("error newWriter [%v]\n", err)
		return
	}
	i.writer = client
	//// TODO: remove when finish structure project

	go func() {
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
	}()
}

var handlePTStream = func(ctx context.Context, data ivstream.InvestmentOpenEntity) error {
	json, _ := json2.Marshal(data)
	fmt.Printf("result is [%v]\n", string(json))
	return nil
}
