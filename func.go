package mruby

import (
	"fmt"
)

// #include "gomruby.h"
import "C"

// Func is the signature of a function in Go that you use to expose to Ruby
// code.
type Func func(m *Mrb, self *Value) *Value

type classMethodMap map[*C.struct_RClass]methodMap
type methodMap map[C.mrb_sym]Func
type stateMethodMap map[*C.mrb_state]classMethodMap


// stateMethodTable is the lookup table for methods that we define in Go and
// expose in Ruby. This is cleaned up by Mrb.Close.
var stateMethodTable stateMethodMap

func init() {
	stateMethodTable = make(stateMethodMap)
}

//export go_mrb_func_call
func go_mrb_func_call(s *C.mrb_state, v *C.mrb_value) C.mrb_value {
	// Lookup the classes that we've registered methods for in this state
	classTable := stateMethodTable[s]
	if classTable == nil {
		panic(fmt.Sprintf("func call from unknown state: %p", s))
	}

	// Get the call info, which we use to lookup the proc
	ci := s.c.ci

	// Lookup the class itself
	methodTable := classTable[ci.proc.target_class]
	if methodTable == nil {
		panic(fmt.Sprintf("func call on unknown class"))
	}

	// Lookup the method
	f := methodTable[ci.mid]
	if f == nil {
		panic(fmt.Sprintf("func call on unknown method"))
	}

	// Call the method to get our *Value
	// TODO(mitchellh): reuse the Mrb instead of allocating every time
	value := f(&Mrb{s}, newValue(s, *v))
	return value.value
}

func insertMethod(s *C.mrb_state, c *C.struct_RClass, n string, f Func) {
	classLookup := stateMethodTable[s]
	if classLookup == nil {
		classLookup = make(classMethodMap)
		stateMethodTable[s] = classLookup
	}

	methodLookup := classLookup[c]
	if methodLookup == nil {
		methodLookup = make(methodMap)
		classLookup[c] = methodLookup
	}

	sym := C.mrb_intern_cstr(s, C.CString(n))
	methodLookup[sym] = f
}

// ArgSpec defines how many arguments a function should take and
// what kind. Multiple ArgSpecs can be combined using the "|"
// operator.
type ArgSpec C.mrb_aspec

// ArgsAny allows any number of arguments.
func ArgsAny() ArgSpec {
	return ArgSpec(C._go_MRB_ARGS_ANY())
}

// ArgsArg says the given number of arguments are required and
// the second number is optional.
func ArgsArg(r, o int) ArgSpec {
	return ArgSpec(C._go_MRB_ARGS_ARG(C.int(r), C.int(o)))
}

// ArgsBlock says it takes a block argument.
func ArgsBlock() ArgSpec {
	return ArgSpec(C._go_MRB_ARGS_BLOCK())
}

// ArgsNone says it takes no arguments.
func ArgsNone() ArgSpec {
	return ArgSpec(C._go_MRB_ARGS_NONE())
}

// ArgsReq says that the given number of arguments are required.
func ArgsReq(n int) ArgSpec {
	return ArgSpec(C._go_MRB_ARGS_REQ(C.int(n)))
}

// ArgsOpt says that the given number of arguments are optional.
func ArgsOpt(n int) ArgSpec {
	return ArgSpec(C._go_MRB_ARGS_OPT(C.int(n)))
}
