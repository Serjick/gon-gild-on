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
// [io/fs.FS] could be used as source (e.g. [os.DirFS] or [embed.FS]). Writes are
// performed over host filesystem if [Writer] option is not used.
type FS struct {
	src       Source
	root      string
	caller    string
	writer    Writer
	locator   Locator
	formatter Formatter
	filter    DataFilter
	updallow  UpdateAllower
	tmplfuncs []TmplFuncFactory
	hooks     Hooks
	overrides map[string][]FSOption
}

// NewFS instantiate [FS].
func NewFS(opts ...FSOption) *FS {
	return newFSDefault().WithOptions(opts...)
}

func newFSDefault() *FS {
	return &FS{
		src:       NewSourceCaller(),
		root:      "",
		caller:    "",
		writer:    Writer{Dir: os.MkdirAll, File: os.WriteFile},
		locator:   NewLocatorDefault(),
		formatter: NewJSONFormatter(),
		filter:    NewDataFilterEmpty(),
		updallow:  NewUpdateAllowerByFlag(),
		tmplfuncs: nil,
		hooks:     NewHooksDefault(),
		overrides: make(map[string][]FSOption),
	}
}

// RenderFile locate and read golden file, render it as text/template with actual data,
// and write result back if `-update` flag defined. File will be auto created if doesn't exists.
// Value of `actual` is the result of operation which being tested, and could be of type
// [golden.Data], [json.RawMessage] or any other type convertible to text.
func (f *FS) RenderFile(t TestingT, actual any) ([]byte, error) {
	var caller string
	if _, c, _, ok := runtime.Caller(1); ok {
		caller = filepath.Dir(c)
	}

	return f.withRootEnsure(caller).handleFile(t, actual)
}

// RenderDir handle same as [FS.RenderFile] every [fs.FS] entry.
// Note that filename from [Locator] is always ignored.
// Symlinks are handled as regular directories/files, and
// like so they are will be saved as golden.
func (f *FS) RenderDir(t TestingT, actual fs.FS) (fs.FS, error) {
	var caller string
	if _, c, _, ok := runtime.Caller(1); ok {
		caller = filepath.Dir(c)
	}

	rootedF := f.withRootEnsure(caller)

	dirFS, visited, err := rootedF.handleDir(t, actual)
	if err != nil {
		return nil, err
	}

	if err := rootedF.gcRedundantFiles(t, visited); err != nil {
		return nil, fmt.Errorf("gc redundant files failure: %w", err)
	}

	return dirFS, nil
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

func (f *FS) getLocation(t tNamer) Location {
	l := f.locator(LocationVars{TestName: t.Name()})

	if l, ok := l.(Location); ok {
		return l
	}

	dir, file := filepath.Split(l.String())

	return Location{
		Dir:  filepath.Clean(dir),
		File: file,
	}
}

func (f *FS) ensureData(actual any) Data { //nolint:ireturn // arbitrary implementations could be returned
	switch cast := actual.(type) {
	case Data:
		return cast
	case json.RawMessage:
		return DataJSON(cast)
	case []byte:
		if c, ok := f.formatter.(DataAdapter); ok {
			return c.AdaptRaw(cast)
		}
	}

	return DataAny{any: actual}
}

func (f *FS) tmplFuncs(t TestingT) template.FuncMap {
	m := make(template.FuncMap)
	for _, tf := range f.tmplfuncs {
		maps.Copy(m, tf(t, TmplFuncFactoryVars{}))
	}

	return m
}
