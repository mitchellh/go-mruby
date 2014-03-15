package mruby

// #cgo CFLAGS: -Ivendor/mruby/include
// #cgo LDFLAGS: -lm libmruby.a
// #include "gomruby.h"
import "C"

import (
)

// Mrb represents a single instance of mruby.
type Mrb struct {
	state *C.mrb_state
}

type ArenaIndex int
type Class *C.struct_RClass

func NewMrb() *Mrb {
	state := C.mrb_open()

	return &Mrb{
		state: state,
	}
}

// Restores the arena index so the objects between the save and this point
// can be garbage collected in the future.
//
// See ArenaSave for more documentation.
func (m *Mrb) ArenaRestore(idx ArenaIndex) {
	C.mrb_gc_arena_restore(m.state, C.int(idx))
}

// This saves the index into the arena.
//
// Restore the arena index later by calling ArenaRestore.
//
// The arena is where objects returned by functions such as LoadString
// are stored. By saving the index and then later restoring it with
// ArenaRestore, these objects can be garbage collected. Otherwise, the
// objects will never be garbage collected.
//
// Of course, when Close() is called, all objects in the arena are
// garbage collected anyways, so if you're only calling mruby for a short
// period of time, you might not have to worry about saving/restoring the
// arena.
func (m *Mrb) ArenaSave() ArenaIndex {
	return ArenaIndex(C.mrb_gc_arena_save(m.state))
}

func (m *Mrb) ObjectClass() Class {
	return Class(m.state.object_class)
}

func (m *Mrb) DefineClass(name string, super Class) Class {
	if super == nil {
		panic("WHAT")
	}

	return Class(C.mrb_define_class(m.state, C.CString(name), super))
}

// LoadString loads the given code, executes it, and returns its final
// value that it might return.
func (m *Mrb) LoadString(code string) (*Value, error) {
	value := C.mrb_load_string(m.state, C.CString(code))
	if m.state.exc != nil {
		return nil, newExceptionValue(m.state)
	}

	return newValue(m.state, value), nil
}

// Close a Mrb, this must be called to properly free resources, and
// should only be called once.
func (m *Mrb) Close() {
	C.mrb_close(m.state)
}
