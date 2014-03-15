package mruby

import (
)

// #include "gomruby.h"
import "C"

// Class is a class in mruby
type Class struct {
	class   *C.struct_RClass
	mrb     *Mrb
}

func (c *Class) DefineClassMethod(name string, cb Func) {
	insertMethod(c.mrb.state, c.class.c, name, cb)

	C.mrb_define_class_method(
		c.mrb.state,
		c.class,
		C.CString(name),
		C._go_mrb_func_t(),
		C._go_MRB_ARGS_ANY())
}

func (c *Class) DefineMethod(name string, cb Func) {
	insertMethod(c.mrb.state, c.class, name, cb)

	C.mrb_define_method(
		c.mrb.state,
		c.class,
		C.CString(name),
		C._go_mrb_func_t(),
		C._go_MRB_ARGS_ANY())
}

func newClass(mrb *Mrb, c *C.struct_RClass) *Class {
	return &Class{
		class:   c,
		mrb:     mrb,
	}
}
