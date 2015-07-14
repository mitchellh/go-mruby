package mruby

import (
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	type structString struct {
		Foo string
	}

	var outBool bool
	var outFloat64 float64
	var outInt int
	var outMap, outMap2 map[string]string
	var outPtrInt *int
	var outSlice []string
	var outString string
	var outStructString structString

	cases := []struct {
		Input    string
		Output   interface{}
		Expected interface{}
	}{
		// Booleans
		{
			"true",
			&outBool,
			true,
		},

		{
			"false",
			&outBool,
			false,
		},

		// Float
		{
			"1.2",
			&outFloat64,
			float64(1.2000000476837158),
		},

		// Int
		{
			"32",
			&outInt,
			int(32),
		},

		{
			`"32"`,
			&outInt,
			int(32),
		},

		// Map
		{
			`{"foo" => "bar"}`,
			&outMap,
			map[string]string{"foo": "bar"},
		},

		{
			`{32 => "bar"}`,
			&outMap2,
			map[string]string{"32": "bar"},
		},

		// Slice
		{
			`["foo", "bar"]`,
			&outSlice,
			[]string{"foo", "bar"},
		},

		// Ptr
		{
			`32`,
			&outPtrInt,
			32,
		},

		// String
		{
			`32`,
			&outString,
			"32",
		},

		{
			`"32"`,
			&outString,
			"32",
		},

		// Struct
		{
			`{"foo" => "bar"}`,
			&outStructString,
			structString{Foo: "bar"},
		},
	}

	for _, tc := range cases {
		mrb := NewMrb()
		value, err := mrb.LoadString(tc.Input)
		if err != nil {
			mrb.Close()
			t.Fatalf("err: %s\n\n%s", err, tc.Input)
		}

		err = Decode(tc.Output, value)
		mrb.Close()
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		val := reflect.ValueOf(tc.Output)
		for val.Kind() == reflect.Ptr {
			val = reflect.Indirect(val)
		}
		actual := val.Interface()
		if !reflect.DeepEqual(actual, tc.Expected) {
			t.Fatalf("bad: %#v\n\n%#v", actual, tc.Expected)
		}
	}
}
