package os

import (
	"os"
)

type Wrapper struct{}

func NewOSWrapper() OS {
	return Wrapper{}
}

func (w Wrapper) FileExists(filename string) bool {
	_, err := os.Stat(filename)

	return !os.IsNotExist(err)
}

func (w Wrapper) MkdirAll(name string, perm os.FileMode) error {
	return os.MkdirAll(name, perm)
}
