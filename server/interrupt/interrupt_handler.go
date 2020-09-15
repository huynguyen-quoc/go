package interrupt

import (
	"os"
	"os/signal"
)

// interruptHandler wraps systems signal receiving channel and provides a better interface for user
type interruptHandler struct {
	signals       []os.Signal
	interruptChan chan os.Signal
	stopChan      chan struct{}
}

// Start starts capturing os signals
func (ih *interruptHandler) Start(signals ...os.Signal) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, signals...)
	ih.signals = signals
	ih.interruptChan = interruptChan
	ih.stopChan = make(chan struct{})
}

// Wait blocks and waits for predefined system signal
func (ih *interruptHandler) Wait() {
	select {
	case <-ih.interruptChan:
		//stopped by interrupt

	case <-ih.stopChan:
		// stopped by user
	}
	signal.Stop(ih.interruptChan)
}

// Stop blocks and waits for predefined system signal
func (ih *interruptHandler) Stop() {
	close(ih.stopChan)
}

// NewInterruptHandler captures system signals
func NewInterruptHandler() Handler {
	return &interruptHandler{}
}
