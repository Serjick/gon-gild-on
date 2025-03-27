package golden

import (
	"os"
)

type (
	DirWriter  func(string, os.FileMode) error
	FileWriter func(string, []byte, os.FileMode) error
	Writer     struct {
		Dir  DirWriter
		File FileWriter
	}
)

const (
	DefaultDirPerm  os.FileMode = os.ModeDir | 0o755
	DefaultFilePerm os.FileMode = 0o644
)
