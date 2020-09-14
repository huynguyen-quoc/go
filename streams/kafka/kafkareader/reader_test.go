package kafkareader

import (
	"github.com/golang/protobuf/proto"
	"github.com/huynguyen-quoc/go/streams/core"
	"github.com/huynguyen-quoc/go/streams/schema/ivstream"
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"testing"
	"time"
)

func Test_initialize(t *testing.T) {
	Convey("Given stream consumer and kafka Entity", t, func() {
		mockConsumer := &core.MockStreamConsumer{}
		kafkaEntity := &ivstream.InvestmentOpenEntity{}

		Convey("When initialize", func() {

			reader := initialize(mockConsumer, kafkaEntity)
			Convey("Then return correct producer and error should be return nil", func() {
				So(reader, ShouldNotBeNil)
				So(reader.consumer, ShouldEqual, mockConsumer)
			})
		})

	})

}

func Test_readerImpl_Done(t *testing.T) {
	Convey("Given shutdown channel", t, func() {
		channel := make(chan struct{})
		data := struct{}{}
		go func() {
			defer close(channel)

			channel <- data
		}()

		reader := readerImpl{
			shutdownChan: channel,
		}

		Convey("When call done", func() {
			done := reader.Done()

			Convey("Then result should returns correct struct result", func() {
				var got []struct{}
				want := []struct{}{
					data,
				}
				for i := range done {
					got = append(got, i)
				}
				So(len(got), ShouldEqual, 1)
				So(reflect.DeepEqual(got, want), ShouldBeTrue)
			})
		})

	})
}

func Test_readerImpl_GetDataChan(t *testing.T) {
	entity := &ivstream.InvestmentOpenEntity{}

	consumerMessage := &core.MockConsumerMessage{}
	value, _ := proto.Marshal(&ivstream.InvestmentOpen{})
	consumerMessage.On("Data").Return(value).Twice()
	invalidConsumerMessage := &core.MockConsumerMessage{}
	invalidMessage := []byte("invalid")
	invalidConsumerMessage.On("Data").Return(invalidMessage).Once()
	entityChannel := make(chan *Entity, 1)
	dataChannel := make(chan core.ConsumerMessage, 1)
	shutdownChan := make(chan struct{})
	mockConsumer := &core.MockStreamConsumer{}
	mockConsumer.On("Shutdown").Once()
	mockConsumer.On("Start").Twice()

	go func() {
		ptReceiveCh := getUnidirectionalChannel(dataChannel, consumerMessage)
		mockConsumer.On("GetDataChan").Return(ptReceiveCh)

		ptReceiveChInvalid := getUnidirectionalChannel(dataChannel, invalidConsumerMessage)
		mockConsumer.On("GetDataChan").Return(ptReceiveChInvalid)

		time.Sleep(100 * time.Millisecond)
		ptReceiveChShutDown := getUnidirectionalChannel(dataChannel, consumerMessage)
		mockConsumer.On("GetDataChan").Return(ptReceiveChShutDown)
	}()

	reader := readerImpl{
		streamEntity: entity,
		consumer:     mockConsumer,
		shutdownChan: shutdownChan,
		dataCh:       entityChannel,
	}

	Convey("GetDataChan success", t, func() {

		ch := reader.GetDataChan()

		Convey("Then result should returns correct kafka entity", func() {
			var got []*Entity
			go func() {
				time.Sleep(200 * time.Millisecond)
				close(entityChannel)
			}()

			for i := range ch {
				got = append(got, i)
			}

			So(len(got), ShouldEqual, 2)
			So(reflect.DeepEqual(got[0].Event, entity), ShouldBeTrue)
		})
	})

	Convey("Message is not correct", t, func() {

		ch := reader.GetDataChan()
		Convey("Then result should returns correct kafka entity", func() {
			var got *Entity
			go func() {
				time.Sleep(200 * time.Millisecond)
				close(entityChannel)
			}()

			for i := range ch {
				got = i
			}

			So(got, ShouldBeNil)
		})
	})

}

func Test_readerImpl_Shutdown(t *testing.T) {

	Convey("When can shutdown the reader", t, func() {
		entity := &ivstream.InvestmentOpenEntity{}
		entityChannel := make(chan *Entity, 1)
		mockConsumer := &core.MockStreamConsumer{}
		mockConsumer.On("Shutdown").Once()

		reader := readerImpl{
			streamEntity: entity,
			consumer:     mockConsumer,
			shutdownChan: make(chan struct{}),
			dataCh:       entityChannel,
		}
		err := reader.Shutdown()

		Convey("Then result should returns nil errors", func() {
			So(err, ShouldBeNil)

		})
	})

}
func Test_readerImpl_processMessage(t *testing.T) {
	Convey("Given reader implement object", t, func() {

		entity := &ivstream.InvestmentOpenEntity{}
		entityChannel := make(chan *Entity, 1)
		mockConsumer := &core.MockStreamConsumer{}
		mockConsumer.On("Shutdown").Once()

		reader := readerImpl{
			streamEntity: entity,
			consumer:     mockConsumer,
			shutdownChan: make(chan struct{}),
			dataCh:       entityChannel,
		}

		Convey("When process message ", func() {

			err := reader.Shutdown()

			Convey("Then result should returns nil errors", func() {
				So(err, ShouldBeNil)

			})
		})

	})

}

func getUnidirectionalChannel(ch chan core.ConsumerMessage, message core.ConsumerMessage) <-chan core.ConsumerMessage {
	ch <- message
	return ch
}

