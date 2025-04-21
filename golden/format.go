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
	_ Formatter = FnFormatter(nil)
	_ Formatter = (*JSONFormatter)(nil)
)

func NewFmtFormatter() FnFormatter {
	return fmt.Sprintln
}

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
