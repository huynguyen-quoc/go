package server

import (
	"context"

	"log"

	"github.com/huynguyen-quoc/go/echo/config"
	"github.com/huynguyen-quoc/go/echo/streams"
	"github.com/huynguyen-quoc/go/service"
	streamConfig "github.com/huynguyen-quoc/go/streams/kafka/config"
)

var (
	brokers = []string{"localhost:9092"}
	topic   = "example-test"
)

type Server struct {
	StreamRunner streams.StreamRunner
}

func (s Server) RunServer() {
	kafkaConfig := &streamConfig.KafkaConfig{
		Brokers:      brokers,
		ClientID:     "12345",
		Stream:       topic,
		KafkaVersion: "2.6.0",
		EnableSync:   false,
	}
	serviceConfig := service.Configs{
		Name:        "Test Service",
		EnabledGRPC: true,
		GRPC:        nil,
		OnShutdown:  []func(){},
	}

	appConfig := &config.AppConfig{
		ServerConfig: serviceConfig,
		Kafka:        kafkaConfig,
	}

	ivStream := &streams.InvestmentStream{Config: appConfig}

	streamServer := streams.Streams{
		IVStream: ivStream,
	}

	go streamServer.Start(context.Background())

	//applicationConfig := &appConfig.AppConfig{Kafka: kafkaConfig, Grpc: grpcConfig, Http: httpConfig}
	serve(appConfig)
}

// serve creates the service and starts the Grab-Kit server
func serve(appConfig *config.AppConfig, configs ...service.Config) {
	o := appConfig.ServerConfig

	for _, opt := range configs {
		opt(&o)
	}

	runner := service.NewService(configs...)
	//echoHandler := &echo.EndpointEchoHandler{Server: runner}
	//echoHandler.RegisterEchoHandlers(handlers.New())

	// Start the service
	if err := runner.Run(context.Background()); err != nil {
		log.Printf("Unable to start service %s due to %#v", "Echo", err)
	}
}
