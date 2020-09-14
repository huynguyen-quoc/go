package server

import (
	"context"
	appConfig "github.com/huynguyen-quoc/go/investment/config"
	"github.com/huynguyen-quoc/go/investment/protocol"
	"github.com/huynguyen-quoc/go/investment/streams"
	"github.com/huynguyen-quoc/go/streams/kafka/config"
)

type Server struct {
	GrpcRunner protocol.Runner
	HttpRunner protocol.Runner
}

var (
	brokers = []string{"localhost:9092"}
	topic   = "example-test"
)

func (s Server) RunServer() error {
	ctx := context.Background()
	kafkaConfig := &config.KafkaConfig{
		Brokers:      brokers,
		ClientID:     "12345",
		Stream:       topic,
		KafkaVersion: "2.6.0",
		EnableSync:   false,
	}
	grpc := appConfig.GrpcConfig{Port: 9000}
	http := appConfig.HttpConfig{Port: 8080}

	applicationConfig := &appConfig.AppConfig{Kafka: kafkaConfig, Grpc: grpc, Http: http}

	ivStream := streams.InvestmentStream{Config: *applicationConfig}

	streamServer := streams.Streams{
		IVStream: ivStream,
	}

	go streamServer.Start(context.Background())

	// run HTTP gateway
	go func() {
		_ = s.HttpRunner.Run(ctx, applicationConfig)
	}()

	return s.GrpcRunner.Run(ctx, applicationConfig)
}
