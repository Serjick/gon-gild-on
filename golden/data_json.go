package golden

import (
	"encoding/json"
	"fmt"
)

var _ Data = DataJSON{}

// DataJSON is a JSON datum representation.
type DataJSON json.RawMessage

// TmplVars converts data into format compartible with text/template.
func (d DataJSON) TmplVars() (any, error) {
	return d.decode()
}

// Format render data as text to save into golden file.
func (d DataJSON) Format(f Formatter) ([]byte, error) {
	b, err := f.Bytes(json.RawMessage(d))
	if err != nil {
		return nil, fmt.Errorf("json data format failure: %w", err)
	}

	return b, nil
}

// Valid checks whether or not data should be saved as golden file.
func (d DataJSON) Valid(f DataFilter) bool {
	v, err := d.decode()

	return err == nil && !f(v)
}

func (d DataJSON) decode() (any, error) {
	var v any
	if err := json.Unmarshal(json.RawMessage(d), &v); err != nil {
		return nil, fmt.Errorf("unmarshal failure: %w", err)
	}

	return v, nil
}
