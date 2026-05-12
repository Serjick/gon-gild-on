package golden

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// gcRedundantFiles remove if allowed all golden files not presented in whitelist.
func (f *FS) gcRedundantFiles(t TestingT, whitelist map[string]struct{}) error {
	dir := filepath.Join(f.root, f.getLocation(t).Dir)

	err := fs.WalkDir(os.DirFS(dir), ".", f.newGCRedundantFilesWalkDirFunc(t, whitelist))
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("%q fs walk failure: %w", dir, err)
	}

	for name, opts := range f.overrides {
		err := f.WithDirEntryOverrides(name).WithOptions(opts...).gcRedundantFiles(t, whitelist)
		if err != nil {
			return err
		}
	}

	return nil
}

// newGCRedundantFilesWalkDirFunc is a [fs.WalkDirFunc] constructor for [gcRedundantFiles].
func (f *FS) newGCRedundantFilesWalkDirFunc(t TestingT, whitelist map[string]struct{}) fs.WalkDirFunc {
	return func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Base(path) == DirKeepMarkerFilename {
			return nil
		}

		fsys, err := f.getInstanceForDirEntry(t, path)
		if err != nil {
			return err
		}

		if !fsys.updallow() {
			return nil
		}

		return fsys.ensureNoRedundant(t, entry, whitelist)
	}
}

// ensureNoRedundant ensures either path present in whitelist
// or no path and it's empty parents present in golden dir.
func (f *FS) ensureNoRedundant(t TestingT, entry fs.DirEntry, whitelist map[string]struct{}) error {
	path := filepath.Join(f.root, f.getLocation(t).String())

	if _, ok := whitelist[path]; ok {
		return nil
	}

	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("remove %q failure: %w", path, err)
	}

	if err := f.gcEmptyParents(path); err != nil {
		return fmt.Errorf("gc %q empty parents failure: %w", path, err)
	}

	if entry.IsDir() {
		return fs.SkipDir
	}

	return nil
}

// gcEmptyParents move upward by directory structure deleting
// empty entries until non empty directory appeared or [FS.root] reached.
func (f *FS) gcEmptyParents(path string) error {
	parent := filepath.Dir(path)

	for parent != "" && parent != "." && parent != f.root {
		entries, err := os.ReadDir(parent)
		if err != nil {
			return fmt.Errorf("read dir %q failure: %w", parent, err)
		}

		if len(entries) > 0 {
			break
		}

		if err := os.RemoveAll(parent); err != nil {
			return fmt.Errorf("remove %q failure: %w", parent, err)
		}

		parent = filepath.Dir(parent)
	}

	return nil
}
