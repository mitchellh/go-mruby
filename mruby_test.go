package mruby

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewMrb(t *testing.T) {
	mrb := NewMrb()
	mrb.Close()
}

func TestMrbArena(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	idx := mrb.ArenaSave()
	mrb.ArenaRestore(idx)
}

func TestMrbClass(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	var class *Class
	class = mrb.Class("Object", nil)
	if class == nil {
		t.Fatal("class should not be nil")
	}

	mrb.DefineClass("Hello", mrb.ObjectClass())
	class = mrb.Class("Hello", mrb.ObjectClass())
	if class == nil {
		t.Fatal("class should not be nil")
	}
}

func TestMrbConstDefined(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	if !mrb.ConstDefined("Object", mrb.ObjectClass()) {
		t.Fatal("Object should be defined")
	}

	mrb.DefineClass("Hello", mrb.ObjectClass())
	if !mrb.ConstDefined("Hello", mrb.ObjectClass()) {
		t.Fatal("Hello should be defined")
	}
}

func TestMrbDefineClass(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	mrb.DefineClass("Hello", mrb.ObjectClass())
	_, err := mrb.LoadString("Hello")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	mrb.DefineClass("World", nil)
	_, err = mrb.LoadString("World")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestMrbDefineClass_methodException(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	cb := func(m *Mrb, self *MrbValue) (Value, Value) {
		v, err := m.LoadString(`raise "exception"`)
		if err != nil {
			exc := err.(*Exception)
			return nil, exc.MrbValue
		}

		return v, nil
	}

	class := mrb.DefineClass("Hello", mrb.ObjectClass())
	class.DefineClassMethod("foo", cb, ArgsNone())
	_, err := mrb.LoadString(`Hello.foo`)
	if err == nil {
		t.Fatal("should error")
	}
}

func TestMrbDefineClassUnder(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	// Define an outer
	hello := mrb.DefineClass("Hello", mrb.ObjectClass())
	_, err := mrb.LoadString("Hello")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Inner
	mrb.DefineClassUnder("World", nil, hello)
	_, err = mrb.LoadString("Hello::World")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Inner defaults
	mrb.DefineClassUnder("Another", nil, nil)
	_, err = mrb.LoadString("Another")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestMrbDefineModule(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	mrb.DefineModule("Hello")
	_, err := mrb.LoadString("Hello")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestMrbDefineModuleUnder(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	// Define an outer
	hello := mrb.DefineModule("Hello")
	_, err := mrb.LoadString("Hello")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Inner
	mrb.DefineModuleUnder("World", hello)
	_, err = mrb.LoadString("Hello::World")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Inner defaults
	mrb.DefineModuleUnder("Another", nil)
	_, err = mrb.LoadString("Another")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestMrbFixnumValue(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	value := mrb.FixnumValue(42)
	if value.Type() != TypeFixnum {
		t.Fatalf("should be fixnum")
	}
}

func TestMrbFullGC(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	ai := mrb.ArenaSave()
	value := mrb.StringValue("foo")
	if value.IsDead() {
		t.Fatal("should not be dead")
	}

	mrb.ArenaRestore(ai)
	mrb.FullGC()
	if !value.IsDead() {
		t.Fatal("should be dead")
	}
}

func TestMrbGetArgs(t *testing.T) {
	cases := []struct {
		args   string
		types  []ValueType
		result []string
	}{
		{
			`("foo")`,
			[]ValueType{TypeString},
			[]string{`"foo"`},
		},

		{
			`(true)`,
			[]ValueType{TypeTrue},
			[]string{`true`},
		},

		{
			`(Hello)`,
			[]ValueType{TypeClass},
			[]string{`Hello`},
		},

		{
			`() { }`,
			[]ValueType{TypeProc},
			nil,
		},

		{
			`(Hello, "bar", true)`,
			[]ValueType{TypeClass, TypeString, TypeTrue},
			[]string{`Hello`, `"bar"`, "true"},
		},

		{
			`("bar", true) {}`,
			[]ValueType{TypeString, TypeTrue, TypeProc},
			nil,
		},
	}

	for _, tc := range cases {
		var actual []*MrbValue
		testFunc := func(m *Mrb, self *MrbValue) (Value, Value) {
			actual = m.GetArgs()
			return self, nil
		}

		mrb := NewMrb()
		defer mrb.Close()
		class := mrb.DefineClass("Hello", mrb.ObjectClass())
		class.DefineClassMethod("test", testFunc, ArgsAny())
		_, err := mrb.LoadString(fmt.Sprintf("Hello.test%s", tc.args))
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		if tc.result != nil {
			if len(actual) != len(tc.result) {
				t.Fatalf("%s: expected %d, got %d",
					tc.args, len(tc.result), len(actual))
			}
		}

		actualStrings := make([]string, len(actual))
		actualTypes := make([]ValueType, len(actual))
		for i, v := range actual {
			str, err := v.Call("inspect")
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			actualStrings[i] = str.String()
			actualTypes[i] = v.Type()
		}

		if !reflect.DeepEqual(actualTypes, tc.types) {
			t.Fatalf("code: %s\nexpected: %#v\nactual: %#v",
				tc.args, tc.types, actualTypes)
		}

		if tc.result != nil {
			if !reflect.DeepEqual(actualStrings, tc.result) {
				t.Fatalf("expected: %#v\nactual: %#v",
					tc.result, actualStrings)
			}
		}
	}
}

func TestMrbLoadString(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	value, err := mrb.LoadString(`"HELLO"`)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if value == nil {
		t.Fatalf("should have value")
	}
}

func TestMrbLoadString_twice(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	value, err := mrb.LoadString(`"HELLO"`)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if value == nil {
		t.Fatalf("should have value")
	}

	value, err = mrb.LoadString(`"WORLD"`)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if value.String() != "WORLD" {
		t.Fatalf("bad: %s", value)
	}
}

func TestMrbLoadStringException(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	_, err := mrb.LoadString(`raise "An exception"`)

	if err == nil {
		t.Fatal("exception expected")
	}

	value, err := mrb.LoadString(`"test"`)
	if err != nil {
		t.Fatal("exception should have been cleared")
	}

	if value.String() != "test" {
		t.Fatal("bad test value returned")
	}
}

func TestMrbRaise(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	cb := func(m *Mrb, self *MrbValue) (Value, Value) {
		return nil, m.GetArgs()[0]
	}

	class := mrb.DefineClass("Hello", mrb.ObjectClass())
	class.DefineClassMethod("foo", cb, ArgsReq(1))
	_, err := mrb.LoadString(`Hello.foo(ArgumentError.new("ouch"))`)
	if err == nil {
		t.Fatal("should have error")
	}
	if err.Error() != "ouch" {
		t.Fatalf("bad: %s", err)
	}
}

func TestMrbYield(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	cb := func(m *Mrb, self *MrbValue) (Value, Value) {
		result, err := m.Yield(m.GetArgs()[0], Int(12), Int(30))
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		return result, nil
	}

	class := mrb.DefineClass("Hello", mrb.ObjectClass())
	class.DefineClassMethod("foo", cb, ArgsBlock())
	value, err := mrb.LoadString(`Hello.foo { |a, b| a + b }`)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if value.Fixnum() != 42 {
		t.Fatalf("bad: %s", value)
	}
}

func TestMrbYieldException(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	cb := func(m *Mrb, self *MrbValue) (Value, Value) {
		result, err := m.Yield(m.GetArgs()[0])
		if err != nil {
			exc := err.(*Exception)
			return nil, exc.MrbValue
		}

		return result, nil
	}

	class := mrb.DefineClass("Hello", mrb.ObjectClass())
	class.DefineClassMethod("foo", cb, ArgsBlock())
	_, err := mrb.LoadString(`Hello.foo { raise "exception" }`)
	if err == nil {
		t.Fatal("should error")
	}

	_, err = mrb.LoadString(`Hello.foo { 1 }`)
	if err != nil {
		t.Fatal("exception should have been cleared")
	}
}

func TestMrbRun(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	parser := NewParser(mrb)
	defer parser.Close()
	context := NewCompileContext(mrb)
	defer context.Close()

	parser.Parse(`
		if $do_raise
			raise "exception"
		else
			"rval"
		end`,
		context,
	)

	proc := parser.GenerateCode()

	// Enable proc exception raising & verify
	mrb.LoadString(`$do_raise = true`)
	_, err := mrb.Run(proc, nil)

	if err == nil {
		t.Fatalf("expected exception, %#v", err)
	}

	// Disable proc exception raising
	// If we still have an exception, it wasn't cleared from the previous invocation.
	mrb.LoadString(`$do_raise = false`)
	rval, err := mrb.Run(proc, nil)
	if err != nil {
		t.Fatalf("unexpected exception, %#v", err)
	}

	if rval.String() != "rval" {
		t.Fatalf("expected return value 'rval', got %#v", rval)
	}
}
