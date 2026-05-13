package golden

import (
	"encoding/json"
	"fmt"
)

var (
	_ Data = DataAny{}
	_ Data = DataJSON{}
)

type (
	// Data is a target representation of arbitrary datum.
	Data interface {
		// TmplVars converts data into format compartible with text/template.
		TmplVars() (any, error)
		// Format render data as text to save into golden file.
		Format(Formatter) ([]byte, error)
		// Valid checks whether or not data should be saved as golden file.
		Valid(DataFilter) bool
	}

	// DataAny is a generic datum representation.
	DataAny struct {
		any
	}

	// DataJSON is a JSON datum representation.
	DataJSON json.RawMessage
)

// TmplVars converts data into format compartible with text/template.
func (d DataAny) TmplVars() (any, error) {
	return d.any, nil
}

// Format render data as text to save into golden file.
func (d DataAny) Format(f Formatter) ([]byte, error) {
	b, err := f.Bytes(d.any)
	if err != nil {
		return nil, fmt.Errorf("data format failure: %w", err)
	}

	return b, nil
}

// Valid checks whether or not data should be saved as golden file.
func (d DataAny) Valid(f DataFilter) bool {
	return !f(d.any)
}

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
