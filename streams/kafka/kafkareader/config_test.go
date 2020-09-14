package kafkareader

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

func TestReaderInitial_NewReader(t *testing.T) {
	mockKafkaConsumer := &kafka.MockConsumerInitialization{}
	mockStreamProducer := &core.MockStreamConsumer{}
	mockKafkaConsumer.On("NewKafkaConsumer", mock.Anything, mock.Anything, mock.Anything).Return(mockStreamProducer, nil).Once()
	mockKafkaConsumer.On("NewKafkaConsumer", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("some errors")).Once()
	mockConfigurer := &kafka.MockConfigurer{}
	mocKafkaConfig := config.KafkaConfig{}
	mockConfigurer.On("GetConfig", mock.Anything).Return(mocKafkaConfig).Twice()
	mockConfigurer.On("GetStreamID", mock.Anything).Return("test-stream").Twice()

	Convey("Given mock initialization for kafka consumer with correct client", t, func() {

		Convey("When new reader is created", func() {
			setup := &ReaderInit{Configurer: mockConfigurer, KafkaInit: mockKafkaConsumer}
			client, err := setup.NewReader(context.Background())
			Convey("Then return correct producer and error should be return nil", func() {
				So(err, ShouldBeNil)
				So(client, ShouldNotBeNil)
			})
		})

	})

	Convey("Given mock initialization for kafka consumer with errors", t, func() {

		Convey("When new kafka consumer return error", func() {
			setup := &ReaderInit{Configurer: mockConfigurer, KafkaInit: mockKafkaConsumer}
			client, err := setup.NewReader(context.Background())
			Convey("Then return error and client is nil", func() {
				So(err.Error(), ShouldEqual, "some errors")
				So(client, ShouldBeNil)

			})
		})
	})
}
