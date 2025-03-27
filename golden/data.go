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
	Data interface {
		// TmplVars converts data into format compartible with text/template.
		TmplVars() (any, error)
		// Format render data as text to save into golden file.
		Format(Formatter) ([]byte, error)
		// Valid checks whether or not data should be saved as golden file.
		Valid(DataFilter) bool
	}

	DataAny struct {
		any
	}

	DataJSON json.RawMessage
)

func (d DataAny) TmplVars() (any, error) {
	return d.any, nil
}

func (d DataAny) Format(f Formatter) ([]byte, error) {
	b, err := f.Bytes(d.any)
	if err != nil {
		return nil, fmt.Errorf("data format failure: %w", err)
	}

	return b, nil
}

func (d DataAny) Valid(f DataFilter) bool {
	return !f(d.any)
}

func (d DataJSON) TmplVars() (any, error) {
	return d.decode()
}

func (d DataJSON) Format(f Formatter) ([]byte, error) {
	b, err := f.Bytes(json.RawMessage(d))
	if err != nil {
		return nil, fmt.Errorf("json data format failure: %w", err)
	}

	return b, nil
}

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
