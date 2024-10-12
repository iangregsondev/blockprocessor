package io

import (
	"io"
)

type Wrapper struct{}

func NewIOWrapper() IO {
	return Wrapper{}
}

func (w Wrapper) ReadAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}
