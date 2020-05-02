package generator

import "bytes"

var _ FileWriter = (*Buffer)(nil)

// Buffer is an in memory FS mock.
// It is useful for testing and debugging.
type Buffer struct {
	currentFile string
	*bytes.Buffer
}

// NewBuffer creates a new in-memory FS buffer.
func NewBuffer() *Buffer {
	var buf = &bytes.Buffer{}
	return &Buffer{
		Buffer: buf,
	}
}

// Open creates a new in-memory file.
func (buffer *Buffer) Open(filename string) error {
	const delim = "\n---------------------------\n"
	var write = func(str string) { _, _ = buffer.WriteString(str) }
	write(delim)
	write("// START OF FILE: ")
	write(filename)
	write(delim)
	write("\n")
	buffer.currentFile = filename
	return nil
}

// Close stops writing the file.
func (buffer *Buffer) Close() error {
	const delim = "\n---------------------------\n"
	var write = func(str string) { _, _ = buffer.WriteString(str) }
	write(delim)
	write("// END OF FILE: ")
	write(buffer.currentFile)
	write(delim)
	write("\n")
	buffer.currentFile = ""
	return nil
}
