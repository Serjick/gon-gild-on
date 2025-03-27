package golden

import (
	"reflect"
)

// DataFilter is for check whether or not data should be filtered.
type DataFilter func(any) bool

// NewDataFilterEmpty instantiate DataFilter to check whether or not
// data is empty, or equal to zero value, or zero length.
func NewDataFilterEmpty() DataFilter {
	return func(data any) bool {
		if data == nil {
			return true
		}

		if reflect.ValueOf(data).IsZero() {
			return true
		}

		defer func() {
			recover() //nolint:errcheck // just prevent crash when data type is not compartible with Len()
		}()

		return reflect.ValueOf(data).Len() == 0
	}
}
