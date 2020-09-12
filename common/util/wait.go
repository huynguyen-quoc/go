package util

import (
	"errors"
	"sync"
	"time"
)

// Wait blocks until WaitGroup is done or timeout
func Wait(wg *sync.WaitGroup, timeout time.Duration) error {
	waitForDone := make(chan struct{})
	go func() {
		defer close(waitForDone)
		wg.Wait()
	}()

	select {
	case <-waitForDone:
		return nil

	case <-time.After(timeout):
		return errors.New("timeout")
	}
}

