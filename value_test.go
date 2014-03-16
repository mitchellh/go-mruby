package mruby

import (
	"testing"
)

func TestMrbValueCall(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	value, err := mrb.LoadString(`"foo"`)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	result, err := value.Call("==", String("foo"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if result.Type() != TypeTrue {
		t.Fatalf("bad type")
	}
}

func TestMrbValueValue(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	falseV := mrb.FalseValue()
	if falseV.MrbValue(mrb) != falseV {
		t.Fatal("should be the same")
	}
}

func TestMrbValueValue_impl(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	var _ Value = mrb.FalseValue()
}

func TestMrbValueFixnum(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	value, err := mrb.LoadString("42")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if value.Fixnum() != 42 {
		t.Fatalf("bad fixnum")
	}
}

func TestMrbValueString(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	value, err := mrb.LoadString(`"foo"`)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if value.String() != "foo" {
		t.Fatalf("bad string")
	}
}

func TestIntMrbValue(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	var value Value = Int(42)
	v := value.MrbValue(mrb)
	if v.Fixnum() != 42 {
		t.Fatalf("bad value")
	}
}

func TestStringMrbValue(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	var value Value = String("foo")
	v := value.MrbValue(mrb)
	if v.String() != "foo" {
		t.Fatalf("bad value")
	}
}
