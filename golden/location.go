package golden

import (
	"fmt"
	"path/filepath"
)

// DefaultFilename is a default golden file basename.
const DefaultFilename = "golden.tmpl"

type (
	// Location is a directory and filename tuple.
	Location struct {
		// Dir is a [Location] directory.
		Dir string
		// File is a [Location] path relative to [Location.Dir].
		// It may or may not contain directory.
		File string
	}

	// LocationVars is a variables for Locator to resolve path to golden file.
	LocationVars struct {
		TestName string
	}

	// Locator is a golden file path resolver.
	Locator func(LocationVars) fmt.Stringer
)

// NewLocatorDefault instantiates default [Locator].
func NewLocatorDefault() Locator {
	return NewLocatorSubDir("")
}

// NewLocatorFilename instantiates [Locator] with custom golden file basename.
func NewLocatorFilename(name string) Locator {
	return NewLocatorSubDirFilename("", name)
}

// NewLocatorSubDir instantiates [Locator] with custom directories.
func NewLocatorSubDir(dir string) Locator {
	return NewLocatorSubDirFilename(dir, DefaultFilename)
}

// NewLocatorSubDirFilename instantiates [Locator] with custom directories and golden file basename.
func NewLocatorSubDirFilename(dir, filename string) Locator {
	return func(v LocationVars) fmt.Stringer {
		return Location{
			Dir:  filepath.Join("testdata", "golden", dir, filepath.FromSlash(v.TestName)),
			File: filename,
		}
	}
}

// WithFile is a [Location.File] immutable setter.
func (l Location) WithFile(f ...string) Location {
	l.File = filepath.Join(f...)

	return l
}

// Rel returns directory relative to [Location.Dir].
func (l Location) Rel() string {
	return filepath.Dir(l.File)
}

// String returns path.
func (l Location) String() string {
	return filepath.Join(l.Dir, l.File)
}
