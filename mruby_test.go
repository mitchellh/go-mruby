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

func TestMrbDefineClass(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	mrb.DefineClass("Hello", mrb.ObjectClass())
	_, err := mrb.LoadString("Hello")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestMrbGetArgs(t *testing.T) {
	cases := []struct {
		args   string
		format string
		types  []ValueType
		result []string
	}{
		{
			`("foo")`,
			"o",
			[]ValueType{TypeString},
			[]string{`"foo"`},
		},

		{
			`(true)`,
			"o",
			[]ValueType{TypeTrue},
			[]string{`true`},
		},

		{
			`(Hello)`,
			"C",
			[]ValueType{TypeClass},
			[]string{`Hello`},
		},

		{
			`() { }`,
			"&",
			[]ValueType{TypeProc},
			nil,
		},
	}

	for _, tc := range cases {
		var actual []interface{}
		testFunc := func(m *Mrb, self *Value) *Value {
			actual = m.GetArgs(tc.format)
			return self
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
		for i, raw := range actual {
			v := raw.(*Value)
			str, err := v.Call("inspect")
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			actualStrings[i] = str.String()
			actualTypes[i] = v.Type()
		}

		if !reflect.DeepEqual(actualTypes, tc.types) {
			t.Fatalf("expected: %#v\nactual: %#v",
				tc.types, actualTypes)
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
