package golden

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
)

// ensureFile reads golden file if it is exists, creates othewise,
// and overwrites if allowed by UpdateAllower.
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

// writeFile dumps actual data formatted by Formatter content as golden file using Writer.
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

	hookVars := PreSaveHookVars{Current: current, TmplFuncs: f.tmplFuncs(t)}
	if err := f.writer.File(path, f.hooks.preSave(t, buf, hookVars), DefaultFilePerm); err != nil {
		return nil, fmt.Errorf("%q write failure: %w", path, err)
	}

	return buf, nil
}
