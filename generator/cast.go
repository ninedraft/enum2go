package generator

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
)

type astCast interface {
	File() (*ast.File, *token.FileSet)
}

type serializedCast struct {
	filename string
	data     []byte
}

func serializedCastFromRe(filename string, src io.Reader) (*serializedCast, error) {
	var fset = token.NewFileSet()
	var data = &bytes.Buffer{}
	var _, errParse = parser.ParseFile(fset, filename, io.TeeReader(src, data), parser.AllErrors)
	if errParse != nil {
		return nil, errParse
	}
	return &serializedCast{
		filename: filename,
		data:     data.Bytes(),
	}, nil
}

func serializedCastFromAST(filename string, file *ast.File) *serializedCast {
	var fset = token.NewFileSet()
	var data = &bytes.Buffer{}
	var errFormat = format.Node(data, fset, file)
	if errFormat != nil {
		panic(errFormat)
	}
	return &serializedCast{
		filename: filename,
		data:     data.Bytes(),
	}
}

func (cast *serializedCast) AST() (*ast.File, *token.FileSet) {
	var fset = token.NewFileSet()
	var src = bytes.NewReader(cast.data)
	var file, errParse = parser.ParseFile(fset, cast.filename, src, parser.AllErrors)
	if errParse != nil {
		panic(errParse)
	}
	return file, fset
}
