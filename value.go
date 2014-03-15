package mruby

import (
	"unsafe"
)

// #include "gomruby.h"
import "C"

// Value represents an mrb_value.
type Value struct {
	value C.mrb_value
	state *C.mrb_state
}

// String returns the "to_s" result of this value.
func (v *Value) String() string {
	value := C.mrb_any_to_s(v.state, v.value)
	result := C.GoString(C.mrb_string_value_ptr(v.state, value))
	return result
}

func newExceptionValue(s *C.mrb_state) *Exception {
	if s.exc == nil {
		panic("exception value init without exception")
	}

	// Convert the RObject* to an mrb_value
	value := C.mrb_obj_value(unsafe.Pointer(s.exc))

	result := newValue(s, value)
	return &Exception{Value: result}
}

func newValue(s *C.mrb_state, v C.mrb_value) *Value {
	return &Value{
		state: s,
		value: v,
	}
}

// Exception is a special type of value that represents an error
// and implements the Error interface.
type Exception struct {
	*Value
}

func (e *Exception) Error() string {
	return e.String()
}
