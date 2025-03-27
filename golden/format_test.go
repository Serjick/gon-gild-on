package golden_test

import (
	"encoding/json"
	"fmt"

	"github.com/Serjick/gon-gild-on/golden"
)

func ExampleJSONFormatter_Bytes_nil() {
	b, err := golden.NewJSONFormatter().Bytes(nil)
	fmt.Printf("%s", b)
	fmt.Print(err)
	// Output:
	// null
	// <nil>
}

func ExampleJSONFormatter_Bytes_struct() {
	var s struct {
		F struct{}
	}
	b, err := golden.NewJSONFormatter().Bytes(&s)
	fmt.Printf("%s", b)
	fmt.Print(err)
	// Output:
	// {
	//     "F": {}
	// }
	// <nil>
}

func ExampleJSONFormatter_Bytes_raw() {
	b, err := golden.NewJSONFormatter().Bytes(json.RawMessage(`{"F": {}}`))
	fmt.Printf("%s", b)
	fmt.Print(err)
	// Output:
	// {
	//     "F": {}
	// }
	// <nil>
}
