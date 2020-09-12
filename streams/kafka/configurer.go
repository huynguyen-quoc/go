package kafka

import (
	"github.com/huynguyen-quoc/go/streams/kafka/config"
)

// Configurer is the kafka config provider interface
//go:generate mockery --inpackage --case underscore --name Configurer
type Configurer interface {
	GetStreamID() string
	GetConfig() config.KafkaConfig
}

