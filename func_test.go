package mruby

import (
	"testing"
)

func testCallback(m *Mrb, self *Value) *Value {
	return m.FixnumValue(42)
}

func testCallbackResult(t *testing.T, v *Value) {
	if v.Type() != TypeFixnum {
		t.Fatalf("bad type: %d", v.Type())
	}

	if v.Fixnum() != 42 {
		t.Fatalf("bad: %d", v.Fixnum())
	}
}
