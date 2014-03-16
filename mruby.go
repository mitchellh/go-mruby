package mruby

import "unsafe"

// #cgo CFLAGS: -Ivendor/mruby/include
// #cgo LDFLAGS: -lm libmruby.a
// #include "gomruby.h"
import "C"

// Mrb represents a single instance of mruby.
type Mrb struct {
	state *C.mrb_state
}

// ArenaIndex represents the index into the arena portion of the GC.
type ArenaIndex int

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

// Define a new top-level class.
func (m *Mrb) DefineClass(name string, super *Class) *Class {
	if super == nil {
		panic("WHAT")
	}

	return newClass(
		m, C.mrb_define_class(m.state, C.CString(name), super.class))
}

// Class returns the class with the given name and superclass. Note that
// if you call this with a class that doesn't exist, mruby will abort the
// application (like a panic, but not a Go panic).
func (m *Mrb) Class(name string, super *Class) *Class {
	var class *C.struct_RClass
	if super == nil {
		class = C.mrb_class_get(m.state, C.CString(name))
	} else {
		class = C.mrb_class_get_under(m.state, super.class, C.CString(name))
	}

	return newClass(m, class)
}

// ConstDefined checks if the given constant is defined in the scope.
//
// This should be used, for example, before a call to Class, because a
// failure in Class will crash your program (by design). You can retrieve
// the Value of a Class by calling Value().
func (m *Mrb) ConstDefined(name string, scope *Value) bool {
	b := C.mrb_const_defined(
		m.state, scope.value, C.mrb_intern_cstr(m.state, C.CString(name)))
	return C.ushort(b) != 0
}

// GetArgs returns all the arguments that were given to the currnetly
// called function (currently on the stack).
func (m *Mrb) GetArgs() []*Value {
	getArgLock.Lock()
	defer getArgLock.Unlock()

	// If we haven't initialized the accumulator yet, do it. We then
	// keep this slice cached around forever.
	if getArgAccumulator == nil {
		getArgAccumulator = make([]*C.mrb_value, 0, 5)
	}

	// Get all the arguments and put it into our accumulator
	C._go_mrb_get_args_all(m.state)

	// Convert those all to values
	values := make([]*Value, len(getArgAccumulator))
	for i, v := range getArgAccumulator {
		values[i] = newValue(m.state, *v)

		// Unset the accumulator value for GC
		getArgAccumulator[i] = nil
	}

	// Clear reset the accumulator to zero length
	getArgAccumulator = getArgAccumulator[:0]

	return values
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
	// Delete all the methods from the state
	delete(stateMethodTable, m.state)

	// Close the state
	C.mrb_close(m.state)
}

//-------------------------------------------------------------------
// Functions below return Values or constant Classes
//-------------------------------------------------------------------

// Returns the Object top-level class.
func (m *Mrb) ObjectClass() *Class {
	return newClass(m, m.state.object_class)
}

// Returns the Object top-level class.
func (m *Mrb) KernelModule() *Class {
	return newClass(m, m.state.kernel_module)
}

// Returns the top-level `self` value.
func (m *Mrb) TopSelf() *Value {
	return newValue(m.state, C.mrb_obj_value(unsafe.Pointer(m.state.top_self)))
}

// Returns a Value for "false"
func (m *Mrb) FalseValue() *Value {
	return newValue(m.state, C.mrb_false_value())
}

// Returns a Value for "true"
func (m *Mrb) TrueValue() *Value {
	return newValue(m.state, C.mrb_true_value())
}

// Returns a Value for a fixed number.
func (m *Mrb) FixnumValue(v int) *Value {
	return newValue(m.state, C.mrb_fixnum_value(C.mrb_int(v)))
}


