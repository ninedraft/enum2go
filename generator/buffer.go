package generator

import "bytes"

var _ FileWriter = (*Buffer)(nil)

type Buffer struct {
	currentFile string
	*bytes.Buffer
}

func NewBuffer() *Buffer {
	var buf = &bytes.Buffer{}
	return &Buffer{
		Buffer: buf,
	}
}

func (buffer *Buffer) Open(filename string) error {
	const delim = "\n---------------------------\n"
	var write = func(str string) { _, _ = buffer.WriteString(str) }
	write(delim)
	write("// START OF FILE: ")
	write(filename)
	write(delim)
	buffer.currentFile = filename
	return nil
}

func (buffer *Buffer) Close() error {
	const delim = "\n---------------------------\n"
	var write = func(str string) { _, _ = buffer.WriteString(str) }
	write(delim)
	write("// END OF FILE: ")
	write(buffer.currentFile)
	write(delim)
	buffer.currentFile = ""
	return nil
}
