package testtools

import (
	"fmt"
	"sync"
)

// functionLocks ensure no race condition, but tries not to compromise on speed by locking individual functions, not entire tests
type functionLocks struct {
	structLocker *sync.Mutex
	locks        *sync.Map
}

func initFuncLocks() *functionLocks {
	f := &functionLocks{
		locks:        &sync.Map{},
		structLocker: &sync.Mutex{},
	}
	return f
}

func (f *functionLocks) createLock(address string) {
	f.structLocker.Lock()
	defer f.structLocker.Unlock()

	_, _ = f.locks.LoadOrStore(address, &sync.Mutex{})
}

func (f *functionLocks) lock(address string) {
	lock, found := f.locks.Load(address)
	if found {
		lock.(*sync.Mutex).Lock()
	} else {
		fmt.Printf("failed to locate & lock functionLock for address %s", address)
	}
}

func (f *functionLocks) unlock(address string) {
	lock, found := f.locks.Load(address)
	if found {
		lock.(*sync.Mutex).Unlock()
	} else {
		fmt.Printf("failed to locate & lock functionLock for address %s", address)
	}
}

func (f *functionLocks) lockAddresses(addresses []string) {
	f.structLocker.Lock()
	defer f.structLocker.Unlock()

	for _, address := range addresses {
		f.lock(address)
	}
}
