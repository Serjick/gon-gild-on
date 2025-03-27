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

const DefaultFilename = "golden.tmpl"

func NewLocatorDefault() Locator {
	return NewLocatorSubDir("")
}

func NewLocatorFilename(name string) Locator {
	return NewLocatorSubDirFilename("", name)
}

func NewLocatorSubDir(dir string) Locator {
	return NewLocatorSubDirFilename(dir, DefaultFilename)
}

func NewLocatorSubDirFilename(dir, filename string) Locator {
	return func(v LocationVars) string {
		return filepath.Join("testdata", "golden", dir, v.TestName, filename)
	}
}
