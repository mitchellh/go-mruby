package mruby

import (
	"testing"
)

func TestArray(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	value, err := mrb.LoadString(`["foo", "bar", "baz"]`)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	v := value.Array()

	// Len
	if n := v.Len(); n != 3 {
		t.Fatalf("bad: %d", n)
	}

	// Get
	value, err = v.Get(1)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if value.String() != "bar" {
		t.Fatalf("bad: %s", value)
	}
}
