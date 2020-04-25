package generator

import (
	"io"
	"os"
	"path/filepath"
)

var _ FileWriter = (*FS)(nil)

// FS is a gadget, which provides write access to OS filesystem.
type FS struct {
	fsProxy

	dir  string
	mode os.FileMode
	flag int
}

type fsProxy interface {
	io.Writer
	io.Closer
}

// NewFS creates a new filesystem gadget.
func NewFS(dir string) *FS {
	return &FS{
		dir:  dir,
		mode: 0755,
		flag: os.O_CREATE | os.O_WRONLY | os.O_TRUNC,
	}
}

// Open opens file and setups the gadget to consume incoming data.
// Teh clos method must be called to ensure data is written to storage device.
func (fs *FS) Open(name string) error {
	var path = filepath.Join(fs.dir, name)
	var file, errOpen = os.OpenFile(path, fs.flag, fs.mode)
	if errOpen != nil {
		return errOpen
	}
	fs.fsProxy = file
	return nil
}
