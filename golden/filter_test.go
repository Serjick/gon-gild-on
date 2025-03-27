package golden_test

import (
	"fmt"

	"github.com/Serjick/gon-gild-on/golden"
)

func ExampleNewDataFilterEmpty_nil() {
	f := golden.NewDataFilterEmpty()
	fmt.Println(f(nil))
	// Output: true
}

func ExampleNewDataFilterEmpty_zero() {
	f := golden.NewDataFilterEmpty()
	var s struct {
		F *struct{}
	}
	fmt.Println(f(s))
	// Output: true
}

func ExampleNewDataFilterEmpty_empty() {
	f := golden.NewDataFilterEmpty()
	fmt.Println(f(make(map[string]struct{})))
	// Output: true
}

func ExampleNewDataFilterEmpty_ok() {
	f := golden.NewDataFilterEmpty()
	s := struct {
		F *struct{}
	}{
		F: new(struct{}),
	}
	fmt.Println(f(s))
	// Output: false
}
