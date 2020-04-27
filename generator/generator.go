package generator

import (
	"go/ast"
	"go/token"
	"io"
	"strconv"

	"golang.org/x/tools/go/ast/astutil"
)

type Config struct {
	typePlaceholder string
}

type FileWriter interface {
	io.Writer
	io.Closer

	Open(filename string) error
}

func mergeFiles(dst *ast.File, files []*ast.File) {
	var imports []*ast.ImportSpec
	for _, file := range files {
		dst.Decls = append(dst.Decls, file.Decls...)
		imports = append(imports, file.Imports...)
	}

	// removing import declarations from merged files
	astutil.Apply(dst, func(cursor *astutil.Cursor) bool {
		var decl, isGenDecl = cursor.Node().(*ast.GenDecl)
		if isGenDecl && decl.Tok == token.IMPORT {
			cursor.Delete()
		}
		return true
	}, nil)

	// merging imports
	var fset = token.NewFileSet()
	for _, imp := range imports {
		var impPath, _ = strconv.Unquote(imp.Path.Value)
		switch {
		case imp.Name != nil:
			var name = imp.Name.String()
			astutil.AddNamedImport(fset, dst, name, impPath)
		default:
			astutil.AddImport(fset, dst, impPath)
		}
	}
}
