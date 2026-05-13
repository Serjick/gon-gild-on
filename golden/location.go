package golden

import (
	"path/filepath"
)

type (
	// LocationVars is a variables for Locator to resolve path to golden file.
	LocationVars struct {
		TestName string
	}

	// Locator is a golden file path resolver.
	Locator func(LocationVars) string
)

// DefaultFilename is a default golden file basename.
const DefaultFilename = "golden.tmpl"

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
	return func(v LocationVars) string {
		return filepath.Join("testdata", "golden", dir, v.TestName, filename)
	}
}
