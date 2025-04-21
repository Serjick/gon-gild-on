# gon-gild-on

Gild up your tests

[![CI](https://github.com/Serjick/gon-gild-on/actions/workflows/ci.yml/badge.svg)](https://github.com/Serjick/gon-gild-on/actions/workflows/ci.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/Serjick/gon-gild-on.svg)](https://pkg.go.dev/github.com/Serjick/gon-gild-on) [![Go Report Card](https://goreportcard.com/badge/github.com/Serjick/gon-gild-on)](https://goreportcard.com/report/github.com/Serjick/gon-gild-on) [![Coverage](https://gist.githubusercontent.com/Serjick/6b5b53429842ee281aebae5fa4473752/raw/coverage.svg)](https://github.com/Serjick/gon-gild-on/actions/workflows/test.yml) [![GolangCI-Lint](https://github.com/Serjick/gon-gild-on/actions/workflows/lint.yml/badge.svg)](https://github.com/Serjick/gon-gild-on/actions/workflows/lint.yml) [![License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://github.com/Serjick/gon-gild-on/blob/main/LICENSE)

Gon-gild-on is a lightweight, flexible library for golden file testing in Go. It provides a clean, composable API for creating, reading, rendering, and validating golden files, with powerful template and diffing capabilities.

## What are Golden Files?

Golden files (or "snapshot testing") is a testing technique where the expected output of a test is stored in a file. During tests, the actual output is compared against the stored "golden" output. This approach is particularly useful for:

- Testing complex data structures with large outputs
- Validating text output, JSON responses, or any structured data
- Reducing boilerplate in tests that verify output hasn't changed

When to Use Golden Files?

- Expected output is large or complex
- Output has a stable structure but variable content
- Visual comparison of differences is helpful
- Changes to output should be explicitly acknowledged by updating the golden files

## Features

- **Flexible Data Handling**: Support for arbitrary Go values, JSON, and extensible data types
- **Template Support**: Use Go's `text/template` in your golden files with powerful template functions
- **Smart Diffs & Updates**: Automatically update golden files while preserving templates
- **Multiple Formatters**: Format output as JSON, Go structs (via go-spew), and more
- **Embeddable**: Works with both regular files and `embed.FS`
- **Extensible**: Easily integrate with your existing testing workflow

## Installation

```sh
go get github.com/Serjick/gon-gild-on
```

## Quick Start

```go
package main

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/Serjick/gon-gild-on/golden"
)

const fixture = `{
    "name": "John",
    "count": 42
}
`

func TestExample(t *testing.T) {
	// Create a new FS for handling golden files
	fs := golden.NewFS()

	// Your test data
	got := fixture

	// Render golden file with actual data
	want, err := fs.RenderFile(t, json.RawMessage(got))
	if err != nil {
		t.FailNow()
	}

	// Assert golden file and test data are equal
	if !reflect.DeepEqual([]byte(got), want) {
		t.FailNow()
	}

	// _, self, _, _ := runtime.Caller(0)
	// fmt.Println(os.ReadFile(filepath.Dir(self) + "/testdata/golden/TestExample/golden.tmpl"))
	// Output:
	// {
	//     "name": "John",
	//     "count": 42
	// }
	// <nil>
}
```
[![Run](goplay.button.svg)](https://go.dev/play/p/Aw3K1wAgNhE)

## Updating Golden Files

Run your tests with the `-update` flag to automatically update golden files:

```sh
go test -update ./...
```

or use `WithForceUpdate` option.

This will update the golden files with the actual test data while template syntax may be preserved.

## Advanced Usage

### Change source directory

```go
// Use current working directory
source := golden.NewSourceCwd()
fs := golden.NewFS(golden.WithFSSource(source))
```

### Custom Formatters

```go
type Fmt func(string, ...any) string

func(f Fmt) Bytes(data any) ([]byte, error) {
	return []byte(f("%#v", data)), nil
}

// Use a stdlib fmt formatter
formatter := Fmt(fmt.Sprintf)
fs := golden.NewFS(golden.WithFSFormatter(formatter))
```

### Custom File Locations

```go
// Put golden files in a specific subdirectory
locator := golden.NewLocatorSubDir("api")
fs := golden.NewFS(golden.WithFSLocator(locator))
```

### Working with JSON

```go
// Create JSON data directly
jsonData := json.RawMessage(`{"key": "value"}`)
want, err := fs.RenderFile(t, jsonData)
```

## Extensions

Golden comes with several extension packages:

- `gildedk8sapimachinery`: Integration with Kubernetes apimachinery for JSON patches
- `gildedsergigodiff`: Smart template diffing using sergi/go-diff
- `gildedspew`: Pretty-printing Go values using davecgh/go-spew
- `gildedtestify`: Template functions for assertions using stretchr/testify

### Using Extensions

#### davecgh/go-spew Formatter for Go Values

```go
import "github.com/Serjick/gon-gild-on/golden/gildedspew"

fs := golden.NewFS(golden.WithFSFormatter(gildedspew.NewFormatter()))
```

#### JSON Merge Patch as actual Data

```go
import "github.com/Serjick/gon-gild-on/golden/gildedk8sapimachinery"

before := []byte(`{"key": "old"}`)
after := []byte(`{"key": "new", "added": true}`)
data := gildedk8sapimachinery.NewDataJSONMergePatch(before, after)
```

#### Testify Assertions in Templates

```go
import "github.com/Serjick/gon-gild-on/golden/gildedtestify"

fs := golden.NewFS(ogolden.WithFSTmplFuncFactory(gildedtestify.NewTmplFuncFactory()))
```

With this you can use assertions in your templates:

```
"uuid": "{{ testifyUUID .Actual.uuid }}",
"timestamp": "{{ testifyTimeRFC3339Nano .Actual.timestamp }}"
"created_at": "{{ testifyTimeInTestcaseRange .Actual.created_at }}"
```

#### Smart Template Diffing

```go
import "github.com/Serjick/gon-gild-on/golden/gildedsergigodiff"

fs := golden.NewFS(
	golden.WithFSPreSaveHook(gildedsergigodiff.NewTextTemplateDiffMatchPatchPreSaveHook()),
)
```

This preserves template expressions when updating golden files.

## Full API Documentation

See the [Go Reference](https://pkg.go.dev/github.com/Serjick/gon-gild-on) for complete API documentation.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the [MIT License](../LICENSE).
