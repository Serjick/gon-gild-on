package golden

import (
	"fmt"
)

var _ Data = DataAny{}

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

	// DataAdapter is to adapt arbitrary data by concrete [Formatter].
	DataAdapter interface {
		// AdaptRaw adapts raw bytes as concrete Data implementation.
		// Error free adapt function should be used, since validity will
		// be checked by any of Data methods.
		AdaptRaw([]byte) Data
	}

	// DataAny is a generic datum representation.
	DataAny struct {
		any
	}
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
