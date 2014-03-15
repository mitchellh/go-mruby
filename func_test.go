package mruby

import (
	"fmt"
	"testing"
)

func testCallback(m *Mrb, self *Value) *Value {
	v := m.FixnumValue(42)
	fmt.Printf("TYPE: %d\n", v.Type())
	return v
}

func testCallbackResult(t *testing.T, v *Value) {
	if v.Type() != TypeFixnum {
		t.Fatalf("bad type: %d", v.Type())
	}
}
