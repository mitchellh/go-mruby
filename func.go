package mruby

import (
	"fmt"
	"log"
	"unsafe"
)

// #include <stdlib.h>
// #include "gomruby.h"
import "C"

// Func is the signature of a function in Go that you use to expose to Ruby
// code.
//
// The first return value is the actual return value for the code.
//
// The second return value is an exception, if any. This will be raised.
type Func func(m *Mrb, self *MrbValue) (Value, Value)

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
func go_mrb_func_call(s *C.mrb_state, v *C.mrb_value, c_exc *C.mrb_value) *C.mrb_value {
	log.Printf("[TRACE] go_mrb_func_call() start")
	defer log.Printf("[TRACE] go_mrb_func_call() finish")

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
	mrb := &Mrb{s}
	result, exc := f(mrb, newValue(s, *v))
	if exc != nil {
		*c_exc = exc.MrbValue(mrb).value
		return &mrb.NilValue().value
	}

	return &result.MrbValue(mrb).value
}

func insertMethod(s *C.mrb_state, c *C.struct_RClass, n string, f Func) {
	log.Printf("[TRACE] insertMethod() start")
	defer log.Printf("[TRACE] insertMethod() finish")

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

	cs := C.CString(n)
	defer C.free(unsafe.Pointer(cs))

	sym := C.mrb_intern_cstr(s, cs)
	methodLookup[sym] = f
}
