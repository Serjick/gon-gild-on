package golden

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
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
	hooks     Hooks
}

// NewFS instantiate FS.
func NewFS(src fs.FS, opts ...FSOption) *FS {
	var root string
	if _, caller, _, ok := runtime.Caller(1); ok {
		root = filepath.Dir(caller)
	}

	newF := &FS{
		src:       src,
		root:      root,
		writer:    Writer{Dir: os.MkdirAll, File: os.WriteFile},
		locator:   NewLocatorDefault(),
		formatter: NewJSONFormatter(),
		filter:    NewDataFilterEmpty(),
		updallow:  NewUpdateAllowerByFlag(),
		hooks:     NewHooksDefault(),
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

	b, err := f.renderTmpl(expected, map[string]any{"Actual": vars})
	if err != nil {
		return nil, fmt.Errorf("template render failure: %w", err)
	}

	return b, nil
}

func (*FS) renderTmpl(tmpl []byte, payload any) ([]byte, error) {
	t, err := template.New("").Funcs(nil).Parse(string(tmpl))
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

func (f *FS) ensureFile(t TestingT, path string, actual Data) ([]byte, error) {
	file, err := f.src.Open(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return f.writeFile(t, filepath.Join(f.root, path), nil, actual)
		}

		return nil, fmt.Errorf("source fille open failure: %w", err)
	}
	defer file.Close() //nolint:errcheck // file opened only for read, close error is not important

	current, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("%q read failure: %w", path, err)
	}

	if f.updallow() {
		return f.writeFile(t, filepath.Join(f.root, path), current, actual)
	}

	return current, nil
}

func (f *FS) writeFile(t TestingT, path string, current []byte, actual Data) ([]byte, error) {
	buf, err := actual.Format(f.formatter)
	if err != nil {
		return nil, fmt.Errorf("%T formating failure: %w", f.formatter, err)
	}

	if !actual.Valid(f.filter) { // skip auto creation of empty file
		return buf, nil
	}

	dir := filepath.Dir(path)
	if err := f.writer.Dir(dir, DefaultDirPerm); err != nil {
		return nil, fmt.Errorf("%q create failure: %w", dir, err)
	}

	hookVars := PreSaveHookVars{Current: current}
	if err := f.writer.File(path, f.hooks.preSave(t, buf, hookVars), DefaultFilePerm); err != nil {
		return nil, fmt.Errorf("%q write failure: %w", path, err)
	}

	return buf, nil
}
