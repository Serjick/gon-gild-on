package golden

import (
	"os"
)

type (
	// DirWriter is a creator of directories.
	DirWriter func(string, os.FileMode) error
	// FileWriter is a creator and writer of files.
	FileWriter func(string, []byte, os.FileMode) error
	// Writer is composion of all writers.
	Writer struct {
		Dir  DirWriter
		File FileWriter
	}
)

const (
	// DefaultDirPerm is a default permissions for [DirWriter].
	DefaultDirPerm os.FileMode = os.ModeDir | 0o755
	// DefaultFilePerm is a default permissions for [FileWriter].
	DefaultFilePerm os.FileMode = 0o644
)
