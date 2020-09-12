package testtools

import (
	"fmt"
	"reflect"
)
// FuncPair stores the:
// - Pointer to the value being stubbed out during the duration of the function
// - The original value which FuncPtr was pointing to
// - The type of the original value which FuncPtr was pointing to
//
// OrigFuncType is needed because nil original values will be represented as the nil empty interface (all
// indistinguishable from each other); there is no way to reconstruct the original nil with the right interface type.
type FuncPair struct {
	FuncPtr      interface{}
	OrigFunc     interface{}
	OrigFuncType reflect.Type
}

// funcLocks prevents the value being stubbed from modification by other tests running concurrently. The lock for the
// memory address containing the value being replaced is taken in Pairs. The lock is released after Restore runs.
var funcLocks = initFuncLocks()

//Restore is run on "defer", before stubbing so that the functions can have their original function stored before being replaced
func Restore(origFunctionPairs ...*FuncPair) {
	for _, p := range origFunctionPairs {
		if p.OrigFunc == nil {
			reflect.ValueOf(p.FuncPtr).Elem().Set(reflect.Zero(p.OrigFuncType))
		} else {
			reflect.ValueOf(p.FuncPtr).Elem().Set(reflect.ValueOf(p.OrigFunc))
		}

		address := fmt.Sprintf("%p", p.FuncPtr)
		funcLocks.unlock(address)
	}
	return
}

// Pairs encapsulates a nice input
func Pairs(origFuncPtrs ...interface{}) (funcPairs []*FuncPair) {
	var addresses []string
	for _, funcPtr := range origFuncPtrs {
		address := fmt.Sprintf("%p", funcPtr)
		funcLocks.createLock(address)
		addresses = append(addresses, address)
	}
	funcLocks.lockAddresses(addresses)
	for _, funcPtr := range origFuncPtrs {
		funcPair := &FuncPair{
			FuncPtr:      funcPtr,
			OrigFunc:     reflect.ValueOf(funcPtr).Elem().Interface(),
			OrigFuncType: reflect.ValueOf(funcPtr).Elem().Type(),
		}
		funcPairs = append(funcPairs, funcPair)
	}
	return
}