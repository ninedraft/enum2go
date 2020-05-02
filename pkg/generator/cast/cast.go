package cast

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
)

type Cast interface {
	File() (*ast.File, *token.FileSet)
}

type Serialized struct {
	name string
	data []byte
}

func MustFromRe(filename string, src io.Reader) *Serialized {
	var cast, err = FromRe(filename, src)
	if err != nil {
		panic(err)
	}
	return cast
}

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

func (cast *Serialized) AST() (*ast.File, *token.FileSet) {
	var fset = token.NewFileSet()
	var src = bytes.NewReader(cast.data)
	var file, errParse = parser.ParseFile(fset, cast.name, src, parser.AllErrors)
	if errParse != nil {
		panic(errParse)
	}
	return file, fset
}

func (cast *Serialized) Reader() *bytes.Reader {
	return bytes.NewReader(cast.data)
}

func (cast *Serialized) Name() string { return cast.name }
