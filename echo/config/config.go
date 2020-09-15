package config

import (
	"github.com/huynguyen-quoc/go/service"
	streamConfig "github.com/huynguyen-quoc/go/streams/kafka/config"
)

type AppConfig struct {
	ServerConfig service.Configs           `json:"serverConfig"`
	Kafka        *streamConfig.KafkaConfig `json:"kafka"`
}

// GetConfig returns the Kafka KafkaConfig for the app
func (app *AppConfig) GetConfig() streamConfig.KafkaConfig {
	return *app.Kafka
}

// GetStreamID returns the topic name of the stream
func (app *AppConfig) GetStreamID() string {
	return app.Kafka.Stream
}
