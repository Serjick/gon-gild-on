package golden

// FSOption is a single setting setter for FS in a immutable way.
type FSOption func(*FS) *FS

// WithFSSource is a golden files Source immutable setter for FS.
func WithFSSource(src Source) FSOption {
	return FSOption(func(fs *FS) *FS {
		return fs.WithSource(src)
	})
}

// WithFSRoot is a root directory immutable setter for FS writes.
func WithFSRoot(dir string) FSOption {
	return FSOption(func(fs *FS) *FS {
		return fs.WithRoot(dir)
	})
}

// WithFSWriter is a dir and file writers implementations immutable setter for FS.
func WithFSWriter(d DirWriter, f FileWriter) FSOption {
	return FSOption(func(fs *FS) *FS {
		return fs.WithWriter(d, f)
	})
}

// WithFSLocator is a golden file location resolver immutable setter for FS.
func WithFSLocator(l Locator) FSOption {
	return FSOption(func(fs *FS) *FS {
		return fs.WithLocator(l)
	})
}

// WithFSFormatter is a golden file content formatter immutable setter for FS.
func WithFSFormatter(f Formatter) FSOption {
	return FSOption(func(fs *FS) *FS {
		return fs.WithFormatter(f)
	})
}

// WithFSDataFilter is a immutable setter for filter of actual data to prevent FS writes.
func WithFSDataFilter(f DataFilter) FSOption {
	return FSOption(func(fs *FS) *FS {
		return fs.WithDataFilter(f)
	})
}

// WithFSForceUpdate is to force golden file overwrite with actual data by FS.
func WithFSForceUpdate() FSOption {
	return FSOption(func(fs *FS) *FS {
		return fs.WithForceUpdate()
	})
}

// WithFSPreSaveHook is a immutable setter for pre golden file FS save hook.
func WithFSPreSaveHook(h PreSaveHook) FSOption {
	return FSOption(func(fs *FS) *FS {
		return fs.WithPreSaveHook(h)
	})
}

// WithFSTmplFuncFactory is a golden files text/template functions factories immutable merger for FS.
func WithFSTmplFuncFactory(tf TmplFuncFactory) FSOption {
	return FSOption(func(fs *FS) *FS {
		return fs.WithTmplFuncFactory(tf)
	})
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
	newF.tmplfuncs = make([]TmplFuncFactory, len(f.tmplfuncs))
	copy(newF.tmplfuncs, f.tmplfuncs)
	newF.tmplfuncs = append(newF.tmplfuncs, tf)

	return &newF
}

// WithPreSaveHook is a immutable setter for pre golden file save hook.
func (f *FS) WithPreSaveHook(h PreSaveHook) *FS {
	newF := *f
	newF.hooks.preSave = h

	return &newF
}
