package golden_test

import (
	"fmt"

	"github.com/Serjick/gon-gild-on/golden"
)

func ExampleNewLocatorDefault() {
	l := golden.NewLocatorDefault()
	fmt.Println(l(golden.LocationVars{
		TestName: "TestFoo/Bar",
	}))
	// Output: testdata/golden/TestFoo/Bar/golden.tmpl
}

func ExampleNewLocatorFilename() {
	l := golden.NewLocatorFilename("example.json")
	fmt.Println(l(golden.LocationVars{
		TestName: "TestFoo",
	}))
	// Output: testdata/golden/TestFoo/example.json
}

func ExampleNewLocatorSubDir() {
	l := golden.NewLocatorSubDir("api")
	fmt.Println(l(golden.LocationVars{
		TestName: "TestFoo",
	}))
	// Output: testdata/golden/api/TestFoo/golden.tmpl
}

func ExampleNewLocatorSubDirFilename() {
	l := golden.NewLocatorSubDirFilename("api", "example.json")
	fmt.Println(l(golden.LocationVars{
		TestName: "TestFoo",
	}))
	// Output: testdata/golden/api/TestFoo/example.json
}
