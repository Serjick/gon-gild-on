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

// JSONFormatter is implementation of Formatter based on encoding/json.
type JSONFormatter struct {
	prefix string
	indent string
}

func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{
		prefix: "",
		indent: "    ",
	}
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
