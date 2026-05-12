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

// WithFSDirEntryOverrides is a immutable setter for FS dir entry options override.
func WithFSDirEntryOverrides(path string, opts ...FSOption) FSOption {
	return FSOption(func(fs *FS) *FS {
		return fs.WithDirEntryOverrides(path, opts...)
	})
}
