package gildedspew_test

import (
	"fmt"

	"github.com/Serjick/gon-gild-on/golden/gildedspew"
)

func ExampleFormatter_Bytes() {
	var s struct {
		F *struct{}
	}
	b, err := gildedspew.NewFormatter().Bytes(&s)
	fmt.Printf("%s", b)
	fmt.Print(err)
	// Output:
	// (*struct { F *struct {} })({
	//  F: (*struct {})(<nil>)
	// })
	// <nil>
}
