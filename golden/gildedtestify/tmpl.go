package gildedtestify

import (
	"html/template"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Serjick/gon-gild-on/golden"
)

// TmplFuncs is a text/template functions collection.
type TmplFuncs struct {
	require.TestingT
	t time.Time
}

var timeNow = time.Now //nolint:gochecknoglobals // private alias

// UUID asserts string is a proper UUID.
func (f TmplFuncs) UUID(s string) string {
	_, err := uuid.Parse(s)
	assert.NoError(f, err)

	return s
}

// Time asserts string is a proper time in specified format.
func (f TmplFuncs) Time(format, s string) string {
	_, err := time.Parse(format, s)
	assert.NoError(f, err)

	return s
}

// TimeInRange asserts string is a proper time and in specific range.
func (f TmplFuncs) TimeInRange(from, to time.Time, s string) string {
	t, err := time.Parse(time.RFC3339Nano, s)
	require.NoError(f, err)
	assert.WithinRange(f, t, from, to)

	return s
}

func (f TmplFuncs) FuncMap() template.FuncMap {
	return template.FuncMap{
		"testifyUUID": f.UUID,
		"testifyTimeRFC3339Nano": func(s string) string {
			return f.Time(time.RFC3339Nano, s)
		},
		"testifyTimeInTestcaseRange": func(s string) string {
			return f.TimeInRange(f.t, time.Now(), s)
		},
	}
}

// NewTmplFuncFactory creates a factory of text/templates functions collection
// to check arbitrary values with testify/assert tests.
func NewTmplFuncFactory() golden.TmplFuncFactory {
	createdAt := timeNow()

	return func(t golden.TestingT, _ golden.TmplFuncFactoryVars) template.FuncMap {
		return TmplFuncs{
			TestingT: t,
			t:        createdAt,
		}.FuncMap()
	}
}
