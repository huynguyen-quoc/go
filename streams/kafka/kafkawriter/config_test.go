package kafkawriter

import (
	"context"
	"errors"
	"github.com/huynguyen-quoc/go/common/testtools"
	"github.com/huynguyen-quoc/go/streams/core"
	"github.com/huynguyen-quoc/go/streams/kafka"
	"github.com/huynguyen-quoc/go/streams/kafka/config"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestStreamSetup_NewWriter(t *testing.T) {
	mockKafkaProducer := &kafka.MockProducerInitialization{}
	mockStreamProducer := &core.MockStreamProducer{}
	mockConfigurer := &kafka.MockConfigurer{}
	mocKafkaConfig := config.KafkaConfig{}
	mockConfigurer.On("GetConfig", mock.Anything).Return(mocKafkaConfig)
	mockConfigurer.On("GetStreamID", mock.Anything).Return("test-stream")

	Convey("Given mock initialization for kafka producer with correct client", t, func() {
		defer testtools.Restore(testtools.Pairs(&mockConfigurer, mockKafkaProducer)...)

		mockKafkaProducer.On("NewKafkaProducer", mock.Anything, mock.Anything, mock.Anything).Return(mockStreamProducer, nil)

		Convey("When new writer", func() {
			setup := &StreamSetup{Configurer: mockConfigurer, KafkaInit: mockKafkaProducer}
			client, err := setup.NewWriter(context.Background())
			Convey("Then return correct producer and error should be return nil", func() {
				So(err, ShouldBeNil)
				So(client, ShouldNotBeNil)
				mockKafkaProducer.AssertNumberOfCalls(t, "NewKafkaProducer", 1)
				mockConfigurer.AssertNumberOfCalls(t, "GetConfig", 1)
				mockConfigurer.AssertNumberOfCalls(t, "GetStreamID", 1)
			})
		})

	})

	Convey("Given mock initialization for kafka producer with errors", t, func() {
		defer testtools.Restore(testtools.Pairs(&mockConfigurer, mockKafkaProducer)...)

		mockKafkaProducer.On("NewKafkaProducer", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("some errors"))
		Convey("When new kafka producer return error", func() {
			setup := &StreamSetup{Configurer: mockConfigurer, KafkaInit: mockKafkaProducer}
			client, _ := setup.NewWriter(context.Background())
			Convey("Then return error and client is nil", func() {
				//So(err.Error(), ShouldEqual, "some errors")
				So(client, ShouldBeNil)
				//mockKafkaProducer.AssertNumberOfCalls(t, "NewKafkaProducer", 1)
				//mockConfigurer.AssertNumberOfCalls(t, "GetConfig", 1)
				//mockConfigurer.AssertNumberOfCalls(t, "GetStreamID", 1)
			})
		})
	})
}
