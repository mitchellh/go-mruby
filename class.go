package mruby

import (
	"log"
	"unsafe"
)

// #include <stdlib.h>
// #include "gomruby.h"
import "C"

// Class is a class in mruby. To obtain a Class, use DefineClass or
// one of the variants on the Mrb structure.
type Class struct {
	class *C.struct_RClass
	mrb   *Mrb
}

// DefineClassMethod defines a class-level method on the given class.
func (c *Class) DefineClassMethod(name string, cb Func, as ArgSpec) {
	log.Printf("[TRACE] Class#DefineClassMethod(%q) start", name)
	defer log.Printf("[TRACE] Class#DefineClassMethod(%q) finish", name)

	insertMethod(c.mrb.state, c.class.c, name, cb)

	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))

	C.mrb_define_class_method(
		c.mrb.state,
		c.class,
		cs,
		C._go_mrb_func_t(),
		C.mrb_aspec(as))
}

// DefineConst defines a constant within this class.
func (c *Class) DefineConst(name string, value Value) {
	log.Printf("[TRACE] Class#DefineConst(%q, %q) start", name, value)
	defer log.Printf("[TRACE] Class#DefineConst(%q, %q) finish", name, value)

	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))

	C.mrb_define_const(
		c.mrb.state, c.class, cs, value.MrbValue(c.mrb).value)
}

// DefineMethod defines an instance method on the class.
func (c *Class) DefineMethod(name string, cb Func, as ArgSpec) {
	log.Printf("[TRACE] Class#DefineMethod(%q) start", name)
	defer log.Printf("[TRACE] Class#DefineMethod(%q) finish", name)

	insertMethod(c.mrb.state, c.class, name, cb)

	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))

	C.mrb_define_method(
		c.mrb.state,
		c.class,
		cs,
		C._go_mrb_func_t(),
		C.mrb_aspec(as))
}

// Value returns a *Value for this Class. *Values are sometimes required
// as arguments where classes should be valid.
func (c *Class) MrbValue(m *Mrb) *MrbValue {
	log.Printf("[TRACE] Class#MrbValue(%#v) start", m)
	defer log.Printf("[TRACE] Class#MrbValue(#%v) finish", m)

	return newValue(c.mrb.state, C.mrb_obj_value(unsafe.Pointer(c.class)))
}

// Instantiate the class with the given args.
func (c *Class) New(args ...Value) (*MrbValue, error) {
	log.Printf("[TRACE] Class#New(%q) start", args)
	defer log.Printf("[TRACE] Class#New(%q) finish", args)

	var argv []C.mrb_value = nil
	var argvPtr *C.mrb_value = nil
	if len(args) > 0 {
		// Make the raw byte slice to hold our arguments we'll pass to C
		argv = make([]C.mrb_value, len(args))
		for i, arg := range args {
			argv[i] = arg.MrbValue(c.mrb).value
		}

		argvPtr = &argv[0]
	}

	result := C.mrb_obj_new(c.mrb.state, c.class, C.mrb_int(len(argv)), argvPtr)
	if c.mrb.state.exc != nil {
		return nil, newExceptionValue(c.mrb.state)
	}

	return newValue(c.mrb.state, result), nil
}

func newClass(mrb *Mrb, c *C.struct_RClass) *Class {
	return &Class{
		class: c,
		mrb:   mrb,
	}
}
