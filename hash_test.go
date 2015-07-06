package mruby

import (
	"testing"
)

func TestHash(t *testing.T) {
	mrb := NewMrb()
	defer mrb.Close()

	value, err := mrb.LoadString(`{"foo" => "bar"}`)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	h := value.Hash()

	// Get
	value, err = h.Get(String("foo"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if value.String() != "bar" {
		t.Fatalf("bad: %s", value)
	}

	// Set
	err = h.Set(String("foo"), String("baz"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	value, err = h.Get(String("foo"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if value.String() != "baz" {
		t.Fatalf("bad: %s", value)
	}
}
