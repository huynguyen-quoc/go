package streams

import "golang.org/x/net/context"

type StreamStartRunner interface {
	Start(ctx context.Context)
}

type StreamStopRunner interface {
	Stop(ctx context.Context)
}

type StreamRunner interface {
	StreamStartRunner
	StreamStopRunner
}

type ProducerRunner interface {
	StreamStartRunner
	SendData(ctx context.Context, data interface{}) error
}

type ConsumerRunner interface {
	StreamStartRunner
	StreamStopRunner
}
