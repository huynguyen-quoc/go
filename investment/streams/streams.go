package streams

import "golang.org/x/net/context"

type Streams struct {
	IVStream StreamRunner
}

func (s Streams) Start(ctx context.Context) {
	s.IVStream.Start(ctx)
}

func (s Streams) Stop(ctx context.Context) {
	s.IVStream.Stop(ctx)
}

