package mruby

import (
	"log"
	"sync"
)

// #include "gomruby.h"
import "C"

// ArgSpec defines how many arguments a function should take and
// what kind. Multiple ArgSpecs can be combined using the "|"
// operator.
type ArgSpec C.mrb_aspec

// ArgsAny allows any number of arguments.
func ArgsAny() ArgSpec {
	log.Printf("[TRACE] ArgsAny() start")
	defer log.Printf("[TRACE] ArgsAny() finish")
	return ArgSpec(C._go_MRB_ARGS_ANY())
}

// ArgsArg says the given number of arguments are required and
// the second number is optional.
func ArgsArg(r, o int) ArgSpec {
	log.Printf("[TRACE] ArgsAny(%d, %d) start", r, o)
	defer log.Printf("[TRACE] ArgsAny(%d, %d) finish", r, o)
	return ArgSpec(C._go_MRB_ARGS_ARG(C.int(r), C.int(o)))
}

// ArgsBlock says it takes a block argument.
func ArgsBlock() ArgSpec {
	log.Printf("[TRACE] ArgsBlock() start")
	defer log.Printf("[TRACE] ArgsBlock() finish")
	return ArgSpec(C._go_MRB_ARGS_BLOCK())
}

// ArgsNone says it takes no arguments.
func ArgsNone() ArgSpec {
	log.Printf("[TRACE] ArgsNone() start")
	defer log.Printf("[TRACE] ArgsNone() finish")
	return ArgSpec(C._go_MRB_ARGS_NONE())
}

// ArgsReq says that the given number of arguments are required.
func ArgsReq(n int) ArgSpec {
	log.Printf("[TRACE] ArgsReq() start")
	defer log.Printf("[TRACE] ArgsReq() finish")
	return ArgSpec(C._go_MRB_ARGS_REQ(C.int(n)))
}

// ArgsOpt says that the given number of arguments are optional.
func ArgsOpt(n int) ArgSpec {
	log.Printf("[TRACE] ArgsOpt(%d) start", n)
	defer log.Printf("[TRACE] ArgsOpt(%d) finish", n)
	return ArgSpec(C._go_MRB_ARGS_OPT(C.int(n)))
}

// The global accumulator when Mrb.GetArgs is called. There is a
// global lock around this so that the access to it is safe.
var getArgAccumulator []*C.mrb_value
var getArgLock sync.Mutex

//export go_get_arg_append
func go_get_arg_append(v *C.mrb_value) {
	getArgAccumulator = append(getArgAccumulator, v)
}
