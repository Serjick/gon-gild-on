package golden

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
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
	src fs.FS

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
func NewFS(src fs.FS, opts ...FSOption) *FS {
	newF := &FS{
		src:       src,
		root:      "",
		writer:    Writer{Dir: os.MkdirAll, File: os.WriteFile},
		locator:   NewLocatorDefault(),
		formatter: NewJSONFormatter(),
		filter:    NewDataFilterEmpty(),
		updallow:  NewUpdateAllowerByFlag(),
		tmplfuncs: nil,
		hooks:     NewHooksDefault(),
	}
	if _, caller, _, ok := runtime.Caller(1); ok {
		newF.root = filepath.Dir(caller)
	}

	for i := range opts {
		newF = opts[i](newF)
	}

	return newF
}

// RenderFile locate and read golden file, render it as text/template with actual data,
// and write result back if `-update` flag defined. File will be auto created if doesn't exists.
// Value of `actual` is the result of operation which being tested, and could be of type
// golden.Data, json.RawMessage or any other type convertible to text.
func (f *FS) RenderFile(t TestingT, actual any) ([]byte, error) {
	data := f.ensureData(actual)
	file := f.locator(LocationVars{TestName: t.Name()})

	expected, err := f.ensureFile(t, file, data)
	if err != nil {
		return nil, fmt.Errorf("ensure %q failure: %w", file, err)
	}

	vars, err := data.TmplVars()
	if err != nil {
		return nil, fmt.Errorf("template payload failure: %w", err)
	}

	b, err := f.renderTmpl(expected, f.tmplFuncs(t), map[string]any{"Actual": vars})
	if err != nil {
		return nil, fmt.Errorf("template render failure: %w", err)
	}

	return b, nil
}

func (*FS) renderTmpl(tmpl []byte, funcs template.FuncMap, payload any) ([]byte, error) {
	t, err := template.New("").Funcs(funcs).Parse(string(tmpl))
	if err != nil {
		return nil, fmt.Errorf("template parse failure: %w", err)
	}

	var b bytes.Buffer
	if err := t.Execute(&b, payload); err != nil {
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
