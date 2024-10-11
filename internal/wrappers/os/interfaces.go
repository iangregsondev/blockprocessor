package os

import (
	"os"
)

type OS interface {
	FileExists(filename string) bool
	MkdirAll(name string, perm os.FileMode) error
}
