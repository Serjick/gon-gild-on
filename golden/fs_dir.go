package golden

import (
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// handleDir is a handler for [FS.RenderDir].
func (f *FS) handleDir(t TestingT, actual fs.FS) (fs.FS, map[string]struct{}, error) {
	tmpDir := t.TempDir()
	saver := newDirEntryFsSave(t, tmpDir, actual)

	err := fs.WalkDir(actual, ".", f.newHandleDirWalkDirFunc(saver))
	if err != nil {
		return nil, nil, fmt.Errorf("actual fs walk failure: %w", err)
	}

	return os.DirFS(tmpDir), saver.visited, nil
}

// newHandleDirWalkDirFunc is a [fs.WalkDirFunc] constructor for [handleDir].
func (f *FS) newHandleDirWalkDirFunc(saver dirEntryFsSave) fs.WalkDirFunc {
	return func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		renderer, err := f.getInstanceForDirEntry(saver.t, path)
		if err != nil {
			return err
		}

		if err := saver.Handle(path, entry, renderer); err != nil {
			return fmt.Errorf("create %q dir entry failure: %w", renderer.getLocation(saver.t).Rel(), err)
		}

		return nil
	}
}

// getInstanceForDirEntry return relative path and overridden [FS]
// if override defined for path, or current [FS] instance otherwise.
func (f *FS) getInstanceForDirEntry(t TestingT, path string) (*FS, error) {
	overrides := slices.Collect(maps.Keys(f.overrides))
	slices.Sort(overrides)

	instance := f.WithLocator(func(LocationVars) fmt.Stringer {
		return f.getLocation(t).WithFile(f.getLocation(t).Rel(), path)
	})
	for _, override := range overrides {
		i, err := f.getOverrideInstance(t, override, path)
		if err == nil {
			instance = i
		}

		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
	}

	return instance, nil
}

// getOverrideInstance return relative path and overridden [FS], or
// [fs.ErrNotExist] if path is not applicable for override specified by name.
func (f *FS) getOverrideInstance(t TestingT, name, path string) (*FS, error) {
	overridden, err := f.getOverride(name, path)
	if err != nil {
		return nil, err
	}

	i, err := overridden.getInstanceForDirEntry(t, path)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (f *FS) getOverride(name, path string) (*FS, error) {
	if len(f.overrides[name]) == 0 {
		return nil, fs.ErrNotExist
	}

	overridden := f.WithDirEntryOverrides(name).WithOptions(f.overrides[name]...)
	target := filepath.Join(overridden.root, strings.TrimPrefix(path, overridden.root+string(filepath.Separator)))
	base := filepath.Join(overridden.root, name)

	rel, err := filepath.Rel(base, target)
	if err != nil {
		return nil, fmt.Errorf("rel %q to %q failure: %w", base, path, err)
	}

	if !filepath.IsLocal(rel) {
		return nil, fs.ErrNotExist
	}

	return overridden, nil
}
