package generator

import (
	"io"
)

type Config struct {
	typePlaceholder string
}

type FileWriter interface {
	io.Writer
	io.Closer

	Open(filename string) error
}
