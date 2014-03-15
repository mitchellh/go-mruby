package mruby

import (
	"unsafe"
)

// #include "gomruby.h"
import "C"

// Value represents an mrb_value.
type Value struct {
	isExc bool
	s     *mrbState
	value C.mrb_value
}

// Close garbage collects a value. This must be called when a Value
// is done being used.
func (v *Value) Close() {
	// TODO: GC value?
	// C.mrb_free(v.s.state, C._go_mrb_ptr(v.value))
	v.s.Close()
}

// IsExc returns true if the value was raised as an exception.
func (v *Value) IsExc() bool {
	return v.isExc
}

// String returns the "to_s" result of this value.
func (v *Value) String() string {
	value := C.mrb_any_to_s(v.s.state, v.value)
	result := C.GoString(C.mrb_string_value_ptr(v.s.state, value))
	// TODO: GC value?
	return result
}

func newExceptionValue(s *mrbState) *Value {
	if s.state.exc == nil {
		panic("exception value init without exception")
	}

	// Convert the RObject* to an mrb_value
	value := C.mrb_obj_value(unsafe.Pointer(s.state.exc))

	result := newValue(s, value)
	result.isExc = true
	return result
}

func newValue(s *mrbState, v C.mrb_value) *Value {
	// Increase the ref count on our state
	s.Open()

	return &Value{
		isExc: false,
		s:     s,
		value: v,
	}
}
