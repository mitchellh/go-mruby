package mruby

// #include "gomruby.h"
import "C"

// Hash represents an MrbValue that is a Hash in Ruby.
//
// A Hash can be obtained by calling the Hash function on MrbValue.
type Hash struct {
	*MrbValue
}

// Get reads a value from the hash.
func (h *Hash) Get(key Value) (*MrbValue, error) {
	keyVal := key.MrbValue(&Mrb{h.state}).value
	result := C.mrb_hash_get(h.state, h.value, keyVal)
	if h.state.exc != nil {
		return nil, newExceptionValue(h.state)
	}

	return newValue(h.state, result), nil
}

// Set sets a value on the hash
func (h *Hash) Set(key, val Value) error {
	keyVal := key.MrbValue(&Mrb{h.state}).value
	valVal := val.MrbValue(&Mrb{h.state}).value

	C.mrb_hash_set(h.state, h.value, keyVal, valVal)
	if h.state.exc != nil {
		return newExceptionValue(h.state)
	}

	return nil
}
