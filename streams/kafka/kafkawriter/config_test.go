package kafkawriter

import (
	"context"
	"errors"
	"testing"

	"github.com/huynguyen-quoc/go/streams/core"
	"github.com/huynguyen-quoc/go/streams/kafka"
	"github.com/huynguyen-quoc/go/streams/kafka/config"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestStreamSetup_NewWriter(t *testing.T) {
	mockKafkaProducer := &kafka.MockProducerInitialization{}
	mockStreamProducer := &core.MockStreamProducer{}
	mockKafkaProducer.On("NewKafkaProducer", mock.Anything, mock.Anything, mock.Anything).Return(mockStreamProducer, nil).Once()
	mockKafkaProducer.On("NewKafkaProducer", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("some errors")).Once()
	mockConfigurer := &kafka.MockConfigurer{}
	mocKafkaConfig := config.KafkaConfig{}
	mockConfigurer.On("GetConfig", mock.Anything).Return(mocKafkaConfig).Twice()
	mockConfigurer.On("GetStreamID", mock.Anything).Return("test-stream").Twice()

	Convey("Given mock initialization for kafka producer with correct client", t, func() {

		Convey("When new writer is created", func() {
			setup := &WriterInit{Configurer: mockConfigurer, KafkaInit: mockKafkaProducer}
			client, err := setup.NewWriter(context.Background())
			Convey("Then return correct producer and error should be return nil", func() {
				So(err, ShouldBeNil)
				So(client, ShouldNotBeNil)
			})
		})

	})

	Convey("Given mock initialization for kafka producer with errors", t, func() {

		Convey("When new kafka producer return error", func() {
			setup := &WriterInit{Configurer: mockConfigurer, KafkaInit: mockKafkaProducer}
			client, err := setup.NewWriter(context.Background())
			Convey("Then return error and client is nil", func() {
				So(err.Error(), ShouldEqual, "some errors")
				So(client, ShouldBeNil)

			})
		})
	})
}
