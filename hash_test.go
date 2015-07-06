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
	value, err = h.Get(String("foo"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if value.String() != "bar" {
		t.Fatalf("bad: %s", value)
	}
}
