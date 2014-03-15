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
func go_mrb_func_call(s *C.mrb_state, v C.mrb_value) C.mrb_value {
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
	value := f(&Mrb{s}, newValue(s, v))
	fmt.Printf("WHAT: %d\n", value.Type())
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
