package interrupt

import (
	"os"
)

// Handler defines methods for handling interrupt
type Handler interface {
	// Start starts capturing os signals
	Start(signals ...os.Signal)

	// Wait blocks and waits for predefined system signal
	Wait()

	// Stop stops signal waiting process explicitly
	Stop()
}
