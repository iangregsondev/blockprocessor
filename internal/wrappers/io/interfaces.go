package io

import (
	"io"
)

type IO interface {
	ReadAll(r io.Reader) ([]byte, error)
}
