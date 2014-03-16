package mruby

import (
	"testing"
)

func testCallback(m *Mrb, self *MrbValue) *MrbValue {
	return m.FixnumValue(42)
}

func testCallbackResult(t *testing.T, v *MrbValue) {
	if v.Type() != TypeFixnum {
		t.Fatalf("bad type: %d", v.Type())
	}

	if v.Fixnum() != 42 {
		t.Fatalf("bad: %d", v.Fixnum())
	}
}
