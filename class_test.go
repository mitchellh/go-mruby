package mruby

import (
	"testing"
)

func TestClassDefineClassMethod(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	class := mrb.DefineClass("Hello", mrb.ObjectClass())
	class.DefineClassMethod("foo", testCallback)
	value, err := mrb.LoadString("Hello.foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testCallbackResult(t, value)
}

func TestClassDefineMethod(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	class := mrb.DefineClass("Hello", mrb.ObjectClass())
	class.DefineMethod("foo", testCallback)
	value, err := mrb.LoadString("Hello.new.foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testCallbackResult(t, value)
}
