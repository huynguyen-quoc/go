package config

import "github.com/huynguyen-quoc/go/streams/kafka/config"

type AppConfig struct {
	Grpc  GrpcConfig          `json:"grpc"`
	Http  HttpConfig          `json:"http"`
	Kafka *config.KafkaConfig `json:"kafka"`
}

type GrpcConfig struct {
	Port int
}

type HttpConfig struct {
	Port int
}

// GetConfig returns the Kafka KafkaConfig for the app
func (app *AppConfig) GetConfig() config.KafkaConfig {
	return *app.Kafka
}

// GetStreamID returns the topic name of the stream
func (app *AppConfig) GetStreamID() string {
	return app.Kafka.Stream
}
