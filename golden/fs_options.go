package golden

import (
	"maps"
	"reflect"
	"slices"
)

// WithOptions is a multiple options immutable setter.
func (f *FS) WithOptions(opts ...FSOption) *FS {
	newF := *f

	newFPtr := &newF
	for i := range opts {
		newFPtr = opts[i](newFPtr)
	}

	return newFPtr
}

// WithSource is a golden files Source immutable setter.
func (f *FS) WithSource(src Source) *FS {
	newF := *f

	newF.src = src

	return &newF
}

// WithRoot is a root directory immutable setter.
func (f *FS) WithRoot(dir string) *FS {
	newF := *f

	newF.root = dir

	return &newF
}

// WithWriter is a dir and file writers implementations immutable setter.
func (f *FS) WithWriter(dir DirWriter, file FileWriter) *FS {
	newF := *f

	newF.writer.Dir = dir
	newF.writer.File = file

	return &newF
}

// WithLocator is a golden file location resolver immutable setter.
func (f *FS) WithLocator(l Locator) *FS {
	newF := *f

	newF.locator = l

	return &newF
}

// WithFormatter is a golden file content formatter immutable setter.
func (f *FS) WithFormatter(fmt Formatter) *FS {
	newF := *f

	newF.formatter = fmt

	return &newF
}

// WithDataFilter is a immutable setter for filter of actual data to prevent writes.
func (f *FS) WithDataFilter(df DataFilter) *FS {
	newF := *f

	newF.filter = df

	return &newF
}

// WithForceUpdate is a immutable setter to always overwrite golden file with actual data.
func (f *FS) WithForceUpdate() *FS {
	newF := *f

	newF.updallow = func() bool {
		return true
	}

	return &newF
}

// WithTmplFuncFactory is a immutable merger of text/template functions collection factories.
func (f *FS) WithTmplFuncFactory(tf TmplFuncFactory) *FS {
	newF := *f

	newF.tmplfuncs = append(slices.Clone(f.tmplfuncs), tf)

	return &newF
}

// WithPreSaveHook is a immutable setter for pre golden file save hook.
func (f *FS) WithPreSaveHook(h PreSaveHook) *FS {
	newF := *f

	newF.hooks.preSave = h

	return &newF
}

// WithDirEntryOverrides is a immutable setter for dir entry options override.
// Passing empty [FSOption] list equivalent to override deletion.
func (f *FS) WithDirEntryOverrides(path string, opts ...FSOption) *FS {
	newF := *f

	overrides := maps.Clone(newF.overrides)
	if overrides == nil {
		overrides = make(map[string][]FSOption)
	}

	overrides[path] = slices.Clone(opts)
	if len(overrides[path]) == 0 {
		delete(overrides, path)
	}

	newF.overrides = overrides

	return &newF
}

func (f *FS) withRootEnsure(caller string) *FS {
	newF := *f

	newF.caller = caller

	if newF.root == "" {
		src := f.src(SourceVars{RenderCallerDir: caller})

		root := caller
		if v := reflect.ValueOf(src); v.Kind() == reflect.String { // NOTE: it is likely os.dirFS
			root = v.String()
		}

		return newF.WithRoot(root)
	}

	return &newF
}
