package mruby

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

// Returns the Object top-level class.
func (m *Mrb) ObjectClass() *Class {
	return newClass(m, m.state.object_class)
}

// Define a new top-level class.
func (m *Mrb) DefineClass(name string, super *Class) *Class {
	if super == nil {
		panic("WHAT")
	}

	return newClass(
		m, C.mrb_define_class(m.state, C.CString(name), super.class))
}

// GetArgs returns the arguments to a called function and should
// be called from a Func callback.
//
// The format string represents what arguments and how many you want
// to retrieve. Each character in the format string represents an
// argument, which will result in in another element in the resulting
// slice. The contents of the slice elements for each format type are
// shown in the table below.
//
//  format specifiers:
//
//    o:      Object         [mrb_value]
//    C:      class/module   [mrb_value]
//    S:      String         [mrb_value]
//    A:      Array          [mrb_value]
//    H:      Hash           [mrb_value]
//    s:      String         [char*,int]            Receive two arguments.
//    z:      String         [char*]                NUL terminated string.
//    a:      Array          [mrb_value*,mrb_int]   Receive two arguments.
//    f:      Float          [mrb_float]
//    i:      Integer        [mrb_int]
//    b:      Boolean        [mrb_bool]
//    n:      Symbol         [mrb_sym]
//    d:      Data           [void*,mrb_data_type const]
//      2nd argument will be used to check data type so it won't
//      be modified
//    &:      Block          [mrb_value]
//    *:      rest argument  [mrb_value*,int]
//      Receive the rest of the arguments as an array.
//    |:      optional
//      Next argument of '|' and later are optional.
//    ?:      optional given [mrb_bool]
//      true if preceding argument (optional) is given.
//
func (m *Mrb) GetArgs(format string) []interface{} {
	result := make([]interface{}, 0, len(format))

	// Iterate over each character, which must return an mrb_value
	for i := 0; i < len(format); i++ {
		f := format[i:i+1]
		switch f {
		case "o":
			fallthrough
		case "C":
			fallthrough
		case "S":
			fallthrough
		case "A":
			fallthrough
		case "H":
			fallthrough
		case "&":
			value := newValue(
				m.state,
				C._go_mrb_get_arg_value(m.state, C.CString(format[i:i+1])))
			result = append(result, value)
		default:
			panic("unknown format type: " + f)
		}
	}

	return result
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
// Functions below return Values
//-------------------------------------------------------------------

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


