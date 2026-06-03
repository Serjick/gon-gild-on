package golden

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// DirKeepMarkerFilename is a filename for unofficial
// convention used to force Git to track an empty directory.
const DirKeepMarkerFilename = ".gitkeep"

// dirEntryFsSave is a [fs.DirEntry] handler to save one as golden.
type dirEntryFsSave struct {
	t      TestingT
	root   string
	source fs.FS

	visited map[string]struct{}
}

// newDirEntryFsSave instantiate [dirEntryFsSave].
func newDirEntryFsSave(t TestingT, root string, source fs.FS) dirEntryFsSave {
	return dirEntryFsSave{
		t:       t,
		root:    root,
		source:  source,
		visited: make(map[string]struct{}),
	}
}

// Handle is for create or replace dir/file as golden.
func (s dirEntryFsSave) Handle(path string, d fs.DirEntry, renderer *FS) error {
	if !d.IsDir() {
		data, err := fs.ReadFile(s.source, path)
		if err != nil {
			return fmt.Errorf("read %q file failure: %w", path, err)
		}

		return s.handleFile(path, data, renderer)
	}

	if err := s.ensureGoldenPathExists(renderer); err != nil {
		return fmt.Errorf("ensure %q failure: %w", renderer.getLocation(s.t), err)
	}

	s.visited[filepath.Join(s.getGoldenDir(renderer), renderer.getLocation(s.t).File)] = struct{}{}

	return s.handleEntry(path, renderer)
}

// handleDir is for create dir/file as golden.
func (s dirEntryFsSave) handleEntry(path string, renderer *FS) error {
	info, err := s.getGoldenInfo(renderer)
	if err != nil {
		return err
	}

	if info.IsDir() {
		err := os.MkdirAll(filepath.Join(s.root, path), DefaultDirPerm)
		if err != nil && !os.IsExist(err) {
			return fmt.Errorf("make %q dir failure: %w", path, err)
		}

		return nil
	}

	return s.handleFile(path, DataAny{any: nil}, renderer)
}

// getGoldenInfo returns golden dir/file [fs.FileInfo].
func (s dirEntryFsSave) getGoldenInfo(renderer *FS) (fs.FileInfo, error) {
	dir, filename := filepath.Split(filepath.Join(renderer.root, renderer.getLocation(s.t).String()))

	dirOrFile, err := os.DirFS(dir).Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open %q failure: %w", renderer.getLocation(s.t), err)
	}
	defer dirOrFile.Close() //nolint:errcheck // file opened only for read, close error is not important

	meta, err := dirOrFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat %q failure: %w", renderer.getLocation(s.t), err)
	}

	return meta, nil
}

// ensureGoldenPathExists guarantees golden path existent.
func (s dirEntryFsSave) ensureGoldenPathExists(renderer *FS) error {
	dir, filename := filepath.Split(filepath.Join(s.getGoldenDir(renderer), renderer.getLocation(s.t).File))

	file, err := os.DirFS(dir).Open(filename)
	switch {
	case err == nil:
		defer file.Close() //nolint:errcheck // file opened only for read, close error is not important

		return s.ensureGoldenDir(renderer, file, filepath.Join(dir, filename))
	case os.IsNotExist(err):
		if er := os.MkdirAll(filepath.Join(dir, filename), DefaultDirPerm); er != nil {
			return fmt.Errorf("make %q dir in %q failure: %w", filename, dir, er)
		}
	default:
		return fmt.Errorf("open %q failure: %w", filepath.Join(dir, filename), err)
	}

	return nil
}

// ensureGoldenDir guarantees golden dir existent if allowed by [FS.updallow].
func (s dirEntryFsSave) ensureGoldenDir(renderer *FS, existent fs.File, path string) error {
	if !renderer.updallow() {
		return nil
	}

	meta, err := existent.Stat()
	if err != nil {
		return fmt.Errorf("stat %q failure: %w", path, err)
	}

	if err := existent.Close(); err != nil {
		return fmt.Errorf("close %q failure: %w", path, err)
	}

	if meta.IsDir() {
		return nil
	}

	return s.replaceGoldenDir(path)
}

// replaceGoldenDir deletes existent file and replaces it with directory.
func (dirEntryFsSave) replaceGoldenDir(path string) error {
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("remove %q file failure: %w", path, err)
	}

	if err := os.Mkdir(path, DefaultDirPerm); err != nil {
		return fmt.Errorf("make %q dir failure: %w", path, err)
	}

	return nil
}

// handleFile is for create or replace file as golden.
func (s dirEntryFsSave) handleFile(path string, actual any, renderer *FS) error {
	data, err := renderer.handleFile(s.t, actual)
	if err != nil {
		return fmt.Errorf("render %q file failure: %w", renderer.getLocation(s.t), err)
	}

	s.visited[filepath.Join(renderer.root, renderer.getLocation(s.t).String())] = struct{}{}

	if err = os.WriteFile(filepath.Join(s.root, path), data, DefaultFilePerm); err != nil {
		return fmt.Errorf("write %q file failure: %w", path, err)
	}

	return nil
}

// getGoldenDir return golden dir for spicified [FS].
func (s dirEntryFsSave) getGoldenDir(f *FS) string {
	return filepath.Join(f.root, f.getLocation(s.t).Dir)
}
