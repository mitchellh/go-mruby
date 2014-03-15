package mruby

// #cgo CFLAGS: -Ivendor/mruby/include
// #cgo LDFLAGS: -lm libmruby.a
// #include "gomruby.h"
import "C"

import (
	"sync/atomic"
)

// Mrb represents a single instance of mruby.
type Mrb struct {
	s *mrbState
}

func NewMrb() *Mrb {
	return &Mrb{
		s: &mrbState{
			ref:   1,
			state: C.mrb_open(),
		},
	}
}

// LoadString loads the given code, executes it, and returns its final
// value.
func (m *Mrb) LoadString(code string) *Value {
	value := C.mrb_load_string(m.s.state, C.CString(code))
	if m.s.state.exc != nil {
		return newExceptionValue(m.s)
	}

	return newValue(m.s, value)
}

// Close a Mrb, this must be called to properly free resources, and
// should only be called once.
func (m *Mrb) Close() {
	m.s.Close()
}

// mrbState wraps a C.mrb_state but keeps track of a reference count
// so that we can clean up the state properly when we're done.
type mrbState struct {
	ref   int32
	state *C.mrb_state
}

func (s *mrbState) Open() {
	atomic.AddInt32(&s.ref, 1)
}

func (s *mrbState) Close() {
	if atomic.AddInt32(&s.ref, -1) == 0 {
		C.mrb_close(s.state)
	}
}
