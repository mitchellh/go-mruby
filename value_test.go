package mruby

import (
	"testing"
)

func TestExceptionString_afterClose(t *testing.T) {
	mrb := NewMrb()
	_, err := mrb.LoadString(`clearly a syntax error`)
	mrb.Close()

	// This panics before the bug fix that this test tests
	err.Error()
}

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

func TestMrbValueCallBlock(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	value, err := mrb.LoadString(`"foo"`)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	block, err := mrb.LoadString(`Proc.new { |_| "bar" }`)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	result, err := value.CallBlock("gsub", String("foo"), block)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if result.Type() != TypeString {
		t.Fatalf("bad type")
	}
	if result.String() != "bar" {
		t.Fatalf("bad: %s", result)
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

func TestMrbValueType(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	cases := []struct {
		Input    string
		Expected ValueType
	}{
		{
			`false`,
			TypeFalse,
		},
		// TypeFree - Type of value after GC collection
		{
			`true`,
			TypeTrue,
		},
		{
			`1`,
			TypeFixnum,
		},
		{
			`:test`,
			TypeSymbol,
		},
		// TypeUndef - Internal value used by mruby for undefined things (instance vars etc)
		// These all seem to get converted to exceptions before hitting userland
		{
			`1.1`,
			TypeFloat,
		},
		// TypeCptr
		{
			`Object.new`,
			TypeObject,
		},
		{
			`Object`,
			TypeClass,
		},
		{
			`module T; end; T`,
			TypeModule,
		},
		// TypeIClass
		// TypeSClass
		{
			`Proc.new { 1 }`,
			TypeProc,
		},
		{
			`[]`,
			TypeArray,
		},
		{
			`{}`,
			TypeHash,
		},
		{
			`"string"`,
			TypeString,
		},
		{
			`1..2`,
			TypeRange,
		},
		{
			`Exception.new`,
			TypeException,
		},
		// TypeFile
		// TypeEnv
		// TypeData
		// TypeFiber
		// TypeMaxDefine
		{
			`nil`,
			TypeNil,
		},
	}

	for _, c := range cases {
		r, err := mrb.LoadString(c.Input)
		if err != nil {
			t.Fatalf("loadstring failed for case %#v: %s", c, err)
		}
		if cType := r.Type(); cType != c.Expected {
			t.Fatalf("bad type: got %s, expected %s", cType, c.Expected)
		}
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
