package gildedtestify_test

import (
	"fmt"
	"os"
	"testing"
	"text/template"

	"github.com/Serjick/gon-gild-on/golden"
	"github.com/Serjick/gon-gild-on/golden/gildedtestify"
)

func ExampleTmplFuncs_UUID_pass() {
	t := new(testing.T)
	f := gildedtestify.NewTmplFuncFactory()
	err := template.Must(template.New("").Funcs(f(t, golden.TmplFuncFactoryVars{})).
		Parse("{{ testifyUUID .v }}")).
		Execute(os.Stdout, map[string]any{"v": "c1b33595-a514-415c-bfa8-b824ada32c4a"})
	fmt.Print(" ", t.Failed(), err)
	// Output:
	// c1b33595-a514-415c-bfa8-b824ada32c4a false <nil>
}

func ExampleTmplFuncs_UUID_fail() {
	t := new(testing.T)
	f := gildedtestify.NewTmplFuncFactory()
	err := template.Must(template.New("").Funcs(f(t, golden.TmplFuncFactoryVars{})).
		Parse("{{ testifyUUID .v }}")).
		Execute(os.Stdout, map[string]any{"v": "z1b33595-a514-415c-bfa8-b824ada32c4a"})
	fmt.Print(" ", t.Failed(), err)
	// Output:
	// z1b33595-a514-415c-bfa8-b824ada32c4a true <nil>
}

func ExampleTmplFuncs_Time_pass() {
	t := new(testing.T)
	f := gildedtestify.NewTmplFuncFactory()
	err := template.Must(template.New("").Funcs(f(t, golden.TmplFuncFactoryVars{})).
		Parse("{{ testifyTimeRFC3339Nano .v }}")).
		Execute(os.Stdout, map[string]any{"v": "2021-01-17T14:28:55.987654Z"})
	fmt.Print(" ", t.Failed(), err)
	// Output:
	// 2021-01-17T14:28:55.987654Z false <nil>
}

func ExampleTmplFuncs_Time_fail() {
	t := new(testing.T)
	f := gildedtestify.NewTmplFuncFactory()
	err := template.Must(template.New("").Funcs(f(t, golden.TmplFuncFactoryVars{})).
		Parse("{{ testifyTimeRFC3339Nano .v }}")).
		Execute(os.Stdout, map[string]any{"v": "2021-01-17 14:28:55"})
	fmt.Print(" ", t.Failed(), err)
	// Output:
	// 2021-01-17 14:28:55 true <nil>
}

func ExampleTmplFuncs_TimeInRange_pass() {
	t := new(testing.T)
	f := gildedtestify.NewTmplFuncFactory()
	err := template.Must(template.New("").Funcs(f(t, golden.TmplFuncFactoryVars{})).
		Parse("{{ testifyTimeInTestcaseRange .v }}")).
		Execute(os.Stdout, map[string]any{"v": "2021-01-17T14:28:55.987654Z"})
	fmt.Print(" ", t.Failed(), err)
	// Output:
	// 2021-01-17T14:28:55.987654Z false <nil>
}

func ExampleTmplFuncs_TimeInRange_fail() {
	t := new(testing.T)
	f := gildedtestify.NewTmplFuncFactory()
	err := template.Must(template.New("").Funcs(f(t, golden.TmplFuncFactoryVars{})).
		Parse("{{ testifyTimeInTestcaseRange .v }}")).
		Execute(os.Stdout, map[string]any{"v": "2021-01-17T14:28:55.987653Z"})
	fmt.Print(" ", t.Failed(), err)
	// Output:
	// 2021-01-17T14:28:55.987653Z true <nil>
}
