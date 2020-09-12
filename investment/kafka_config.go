package main

import (
	"github.com/huynguyen-quoc/go/streams/kafka/config"
)

type Config struct {
	// service specific configs
	Kafka *config.KafkaConfig `json:"kafka"`
}

// GetConfig returns the Kafka Config for the app
func (app *Config) GetConfig() config.KafkaConfig {
	return *app.Kafka
}

// GetStreamID returns the topic name of the stream
func (app *Config) GetStreamID() string {
	return app.Kafka.Stream
}
