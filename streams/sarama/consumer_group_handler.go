package sarama

import (
	"github.com/Shopify/sarama"
	"sync"
	"time"
)

// ConsumerGroupReadHandler represents a Sarama consumer group consumer
type ConsumerGroupReadHandler struct {
	readFunc         func(m *sarama.ConsumerMessage, session sarama.ConsumerGroupSession)
	saramaClient     sarama.Client
	topic            string
	commitInterval   time.Duration
	clientID         string
	readWg           *sync.WaitGroup
}


func (h *ConsumerGroupReadHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h *ConsumerGroupReadHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (h *ConsumerGroupReadHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29

	h.readWg.Add(1)

	for message := range claim.Messages() {
		h.readFunc(message, session)
	}

	h.readWg.Done()
	return nil
}
