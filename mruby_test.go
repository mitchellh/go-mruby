package mruby

import (
	"testing"
)

func TestNewMrb(t *testing.T) {
	mrb := NewMrb()
	mrb.Close()
}

func TestMrbLoadString(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	value := mrb.LoadString("p 'HELLO'")
	if value.IsExc() {
		t.Fatalf("err: %s", value)
	}
}
