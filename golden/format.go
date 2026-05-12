package golden

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Formatter is golden file content converter.
type Formatter interface {
	// Bytes converts arbitrary data into text to use as golden file content.
	Bytes(any) ([]byte, error)
}

type (
	// FnFormatter is implementation of Formatter compartible with stdlib fmt package.
	FnFormatter func(...any) string

	// JSONFormatter is implementation of Formatter based on encoding/json.
	JSONFormatter struct {
		prefix string
		indent string
	}
)

var (
	_ Formatter   = FnFormatter(nil)
	_ Formatter   = (*JSONFormatter)(nil)
	_ DataAdapter = (*JSONFormatter)(nil)
)

// NewFmtFormatter instantiates [fmt.Sprintln] as [FnFormatter].
func NewFmtFormatter() FnFormatter {
	return fmt.Sprintln
}

// NewStrFormatter instantiates [fmt.Sprintf] with simple string pattern as [FnFormatter].
func NewStrFormatter() FnFormatter {
	return NewPatternFormatter("%s")
}

// NewPatternFormatter instantiates [fmt.Sprintf] with specified pattern as [FnFormatter].
func NewPatternFormatter(pattern string) FnFormatter {
	return func(args ...any) string {
		return fmt.Sprintf(pattern, args...)
	}
}

// NewJSONFormatter instantiates [JSONFormatter].
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{
		prefix: "",
		indent: "    ",
	}
}

// Bytes dump any data as is.
func (f FnFormatter) Bytes(data any) ([]byte, error) {
	return []byte(f(data)), nil
}

// Bytes dump any data as json.
func (f *JSONFormatter) Bytes(data any) ([]byte, error) {
	var buf bytes.Buffer

	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent(f.prefix, f.indent)

	if err := enc.Encode(data); err != nil {
		return nil, fmt.Errorf("json encode failure: %w", err)
	}

	return buf.Bytes(), nil
}

// AdaptRaw adapts raw bytes as JSON Data.
func (*JSONFormatter) AdaptRaw(b []byte) Data { //nolint:ireturn // [DataAdapter] implementation
	return DataJSON(b)
}
