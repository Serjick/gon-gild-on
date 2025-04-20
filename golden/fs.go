package golden

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"text/template"
)

// FS is golden files read and write store wrapper. Any implementation of
// io/fs.FS could be used as source (e.g. os.DirFS or embed.FS). Writes are
// performed over host filesystem if Writer option is not used.
type FS struct {
	src       Source
	root      string
	writer    Writer
	locator   Locator
	formatter Formatter
	filter    DataFilter
	updallow  UpdateAllower
	tmplfuncs []TmplFuncFactory
	hooks     Hooks
}

// NewFS instantiate FS.
func NewFS(opts ...FSOption) *FS {
	newF := newFSDefault()
	for i := range opts {
		newF = opts[i](newF)
	}

	return newF
}

func newFSDefault() *FS {
	return &FS{
		src:       NewSourceCaller(),
		root:      "",
		writer:    Writer{Dir: os.MkdirAll, File: os.WriteFile},
		locator:   NewLocatorDefault(),
		formatter: NewJSONFormatter(),
		filter:    NewDataFilterEmpty(),
		updallow:  NewUpdateAllowerByFlag(),
		tmplfuncs: nil,
		hooks:     NewHooksDefault(),
	}
}

// RenderFile locate and read golden file, render it as text/template with actual data,
// and write result back if `-update` flag defined. File will be auto created if doesn't exists.
// Value of `actual` is the result of operation which being tested, and could be of type
// golden.Data, json.RawMessage or any other type convertible to text.
func (f *FS) RenderFile(t TestingT, actual any) ([]byte, error) {
	var caller string
	if _, c, _, ok := runtime.Caller(1); ok {
		caller = filepath.Dir(c)
	}

	data := f.ensureData(actual)
	file := f.locator(LocationVars{TestName: t.Name()})

	expected, err := f.ensureFile(t, caller, file, data)
	if err != nil {
		return nil, fmt.Errorf("ensure %q failure: %w", file, err)
	}

	b, err := f.renderTmpl(expected, f.tmplFuncs(t), data)
	if err != nil {
		return nil, fmt.Errorf("template render failure: %w", err)
	}

	return b, nil
}

func (*FS) renderTmpl(tmpl []byte, funcs template.FuncMap, actual Data) ([]byte, error) {
	vars, err := actual.TmplVars()
	if err != nil {
		return nil, fmt.Errorf("template payload failure: %w", err)
	}

	t, err := template.New("").Funcs(funcs).Parse(string(tmpl))
	if err != nil {
		return nil, fmt.Errorf("template parse failure: %w", err)
	}

	var b bytes.Buffer
	if err := t.Execute(&b, map[string]any{"Actual": vars}); err != nil {
		return nil, fmt.Errorf("template execute failure: %w", err)
	}

	return b.Bytes(), nil
}

func (*FS) ensureData(actual any) Data { //nolint:ireturn // arbitrary implementations could be returned
	switch t := actual.(type) {
	case Data:
		return t
	case json.RawMessage:
		return DataJSON(t)
	default:
		return DataAny{any: t}
	}
}

func (f *FS) tmplFuncs(t TestingT) template.FuncMap {
	m := make(template.FuncMap)
	for _, tf := range f.tmplfuncs {
		maps.Copy(m, tf(t, TmplFuncFactoryVars{}))
	}

	return m
}
