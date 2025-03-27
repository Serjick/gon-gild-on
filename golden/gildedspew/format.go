package gildedspew

import (
	"bytes"

	"github.com/davecgh/go-spew/spew"

	"github.com/Serjick/gon-gild-on/golden"
)

var _ golden.Formatter = (*Formatter)(nil)

// Formatter is implementation of golden.Formatter based on github.com/davecgh/go-spew/spew.
type Formatter struct {
	s *spew.ConfigState
}

func NewFormatter() *Formatter {
	cfg := spew.NewDefaultConfig()
	cfg.DisablePointerAddresses = true
	cfg.DisableCapacities = true
	cfg.SortKeys = true

	return &Formatter{
		s: cfg,
	}
}

// Bytes dump any data with spew.Fdump().
func (f *Formatter) Bytes(data any) ([]byte, error) {
	var b bytes.Buffer

	f.s.Fdump(&b, data)

	return b.Bytes(), nil
}
