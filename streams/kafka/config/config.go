package config

// IteratorType is a custom type related to consumer/iterator
type IteratorType string

// QueueType is a custom type to choose the appropriate queue
type QueueType string

// OffsetType is a custom type to choose the appropriate offset
type OffsetType string

// OffsetType is a custom type to choose the appropriate offset
type RequireAcksType int16

// Custom Iterator Types
const (
	Oldest   IteratorType = "oldest"
	Sequence              = "sequence"
	Latest                = "latest"
)

// Custom Offset Types for Consumer
const (
	// OffsetNewest to consume latest msgs
	OffsetNewest OffsetType = "newest"
	// OffsetOldest to consume oldest msgs
	OffsetOldest = "oldest"
	// NoResponse doesn't send any response, the TCP ACK is all you get.
	NoResponse RequireAcksType = 0
	// WaitForLocal waits for only the local commit to succeed before responding.
	WaitForLocal RequireAcksType = 1
	// WaitForAll waits for all in-sync replicas to commit before responding.
	WaitForAll RequireAcksType = -1
)



// KafkaConfig defines the basic config for sk2
type KafkaConfig struct {
	// Cluster and Stream Configurations used by both Producers & Consumers
	Brokers      []string `json:"brokers"`
	ClientID     string   `json:"clientID"`
	ClusterType  string   `json:"clusterType"`
	Enabled      bool     `json:"enabled"`
	KafkaVersion string   `json:"kafkaVersion"`
	Stream       string   `json:"stream"`

	// Consumer Configurations
	ConsumerGroupID string     `json:"consumerGroupID"`
	OffsetType      OffsetType `json:"offsetType"`

	// Producer Configurations
	CompressionCodec string `json:"compressionCodec"`
	CompressionLevel int    `json:"compressionLevel"`
	EnableSync       bool   `json:"enableSync"`
	RequiredAcks     int16  `json:"requiredAcks"` // number of acks required from replicas (passed in Produce Requests)

	AppName string `json:"appname"`
}

// Validate will the settings are complete and valid or panic
func (cfg KafkaConfig) Validate() {
	errorString := ""
	if cfg.Enabled {
		if len(cfg.Brokers) <= 0 {
			errorString += "Brokers must be set\n"
		}

		if len(cfg.ClientID) == 0 {
			errorString += "clientID must be set\n"
		}
		if len(errorString) > 0 {
			panic(errorString)
		}
	}
}
