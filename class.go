package mruby

import "unsafe"

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
	insertMethod(c.mrb.state, c.class.c, name, cb)

	C.mrb_define_class_method(
		c.mrb.state,
		c.class,
		C.CString(name),
		C._go_mrb_func_t(),
		C.mrb_aspec(as))
}

// DefineMethod defines an instance method on the class.
func (c *Class) DefineMethod(name string, cb Func, as ArgSpec) {
	insertMethod(c.mrb.state, c.class, name, cb)

	C.mrb_define_method(
		c.mrb.state,
		c.class,
		C.CString(name),
		C._go_mrb_func_t(),
		C.mrb_aspec(as))
}

// Value returns a *Value for this Class. *Values are sometimes required
// as arguments where classes should be valid.
func (c *Class) MrbValue() *MrbValue {
	return newValue(c.mrb.state, C.mrb_obj_value(unsafe.Pointer(c.class)))
}

func newClass(mrb *Mrb, c *C.struct_RClass) *Class {
	return &Class{
		class: c,
		mrb:   mrb,
	}
}
