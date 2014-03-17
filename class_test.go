package mruby

import (
	"testing"
)

func TestClassDefineClassMethod(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	class := mrb.DefineClass("Hello", mrb.ObjectClass())
	class.DefineClassMethod("foo", testCallback, ArgsNone())
	value, err := mrb.LoadString("Hello.foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testCallbackResult(t, value)
}

func TestClassDefineConst(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	class := mrb.DefineClass("Hello", mrb.ObjectClass())
	class.DefineConst("FOO", String("bar"))
	value, err := mrb.LoadString("Hello::FOO")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if value.String() != "bar" {
		t.Fatalf("bad: %s", value)
	}
}

func TestClassDefineMethod(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	class := mrb.DefineClass("Hello", mrb.ObjectClass())
	class.DefineMethod("foo", testCallback, ArgsNone())
	value, err := mrb.LoadString("Hello.new.foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testCallbackResult(t, value)
}

func TestClassNew(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	class := mrb.DefineClass("Hello", mrb.ObjectClass())
	class.DefineMethod("foo", testCallback, ArgsNone())

	instance, err := class.New()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	value, err := instance.Call("foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testCallbackResult(t, value)
}

func TestClassValue(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	class := mrb.DefineClass("Hello", mrb.ObjectClass())
	value := class.MrbValue(mrb)
	if value.Type() != TypeClass {
		t.Fatalf("bad: %d", value.Type())
	}
}
