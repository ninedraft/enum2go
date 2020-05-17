package cast

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
)

// Cast describes a generic AST template.
type Cast interface {
	File() (*ast.File, *token.FileSet)
}

// Serialized is byte buffer baked AST template.
type Serialized struct {
	name string
	data []byte
}

// MustFromRe creates a serialized AST from given byte source.
// Panics on error
func MustFromRe(filename string, src io.Reader) *Serialized {
	var cast, err = FromRe(filename, src)
	if err != nil {
		panic(err)
	}
	return cast
}

// FromRe creates a new serialized AST template from given byte source.
// Returns any error encountered.
func FromRe(filename string, src io.Reader) (*Serialized, error) {
	var fset = token.NewFileSet()
	var data = &bytes.Buffer{}
	switch src := src.(type) {
	case *io.LimitedReader:
		var size = int(src.N)
		data.Grow(size)
	case interface{ Len() int }:
		data.Grow(src.Len())
	case interface{ Size() int64 }:
		var size = int(src.Size())
		data.Grow(size)
	}
	var _, errParse = parser.ParseFile(fset, filename, io.TeeReader(src, data), parser.AllErrors)
	if errParse != nil {
		return nil, errParse
	}
	return &Serialized{
		name: filename,
		data: data.Bytes(),
	}, nil
}

// FromAST creates a new serialized AST template from given AST.
// If the AST template is invalid, this function can panic.
func FromAST(filename string, file *ast.File) *Serialized {
	var fset = token.NewFileSet()
	var data = &bytes.Buffer{}
	var errFormat = format.Node(data, fset, file)
	if errFormat != nil {
		panic(errFormat)
	}
	return &Serialized{
		name: filename,
		data: data.Bytes(),
	}
}

// AST return the internal AST template.
func (cast *Serialized) AST() (*ast.File, *token.FileSet) {
	var fset = token.NewFileSet()
	var src = bytes.NewReader(cast.data)
	var file, errParse = parser.ParseFile(fset, cast.name, src, parser.AllErrors)
	if errParse != nil {
		panic(errParse)
	}
	return file, fset
}

// Reader returns the template as a byte reader.
func (cast *Serialized) Reader() *bytes.Reader {
	return bytes.NewReader(cast.data)
}

// Name returns the filename of the serialized template.
func (cast *Serialized) Name() string { return cast.name }
