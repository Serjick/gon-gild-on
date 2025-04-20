package golden

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
)

type (
	// SourceVars is a for Source to instantiate fs.Fs.
	SourceVars struct {
		// RenderCallerDir is a directory of `FS.RenderFile` caller.
		RenderCallerDir string
	}
	// Source is a factory of filesystem where golden files are located.
	Source func(SourceVars) fs.FS
)

// NewSourceCwd creates Source which uses directory from
// which `go test` has been run as a fs.FS root.
func NewSourceCwd() Source {
	src := os.DirFS(".")

	return func(SourceVars) fs.FS {
		return src
	}
}

// NewSourceDir creates Source which uses specified dir as a fs.FS root.
func NewSourceDir(dir string) Source {
	return func(SourceVars) fs.FS {
		return os.DirFS(dir)
	}
}

// NewSourceFS creates Source which utilizes specified fs.FS.
func NewSourceFS(f fs.FS) Source {
	return func(SourceVars) fs.FS {
		return f
	}
}

// MustNewSourceRel creates Source which uses directory
// relative to caller as a fs.FS root.
func MustNewSourceRel() Source {
	_, caller, _, ok := runtime.Caller(1)
	if !ok {
		panic("no reative dir found")
	}
	src := os.DirFS(filepath.Dir(caller))

	return func(SourceVars) fs.FS {
		return src
	}
}

// NewSourceCaller creates Source which uses
// `SourceVars.RenderCallerDir` as a fs.FS root.
func NewSourceCaller() Source {
	return func(v SourceVars) fs.FS {
		return os.DirFS(v.RenderCallerDir)
	}
}
