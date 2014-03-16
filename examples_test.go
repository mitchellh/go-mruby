package mruby

import (
	"fmt"
)

func ExampleCustomFunction() {
	mrb := NewMrb()
	defer mrb.Close()

	// Our custom function we'll expose to Ruby
	addFunc := func(m *Mrb, self *MrbValue) Value {
		args := m.GetArgs()
		return Int(args[0].Fixnum() + args[1].Fixnum())
	}

	// Lets define a custom class and a class method we can call.
	class := mrb.DefineClass("Example", nil)
	class.DefineClassMethod("add", addFunc, ArgsReq(2))

	// Let's call it and inspect the result
	result, err := mrb.LoadString(`Example.add(12, 30)`)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Result: %s\n", result.String())
	// Output:
	// Result: 42
}
