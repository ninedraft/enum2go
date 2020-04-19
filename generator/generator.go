package generator

import (
	"go/ast"
	"go/token"
	"io"
	"strconv"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

type Generator struct {
	targetType      string
	typePlaceholder string
	casts           map[string]*ast.File
}

type FileWriter interface {
	io.Writer
	io.Closer

	Open(filename string) error
}

func (gen *Generator) Generate(dst FileWriter) error {
	return nil
}

func (gen *Generator) patchAST(node ast.Node) {
	astutil.Apply(node, func(cursor *astutil.Cursor) bool {
		var ident, isIdent = cursor.Node().(*ast.Ident)
		if isIdent {
			ident.Name = strings.ReplaceAll(ident.Name, gen.typePlaceholder, gen.targetType)
		}
		return true
	}, nil)
}

func (gen *Generator) toPour() []string {
	return []string{
		"gadget.go",
		"helpers.go",
	}
}

type cSpec struct {
	Name  string
	Value interface{}
}

func (spec *cSpec) Ref() *ast.Ident { return ast.NewIdent(spec.Name) }

func (spec *cSpec) ValueSpec() *ast.ValueSpec {
	return &ast.ValueSpec{
		Names: []*ast.Ident{spec.Ref()},
	}
}

func mergeFiles(dst *ast.File, files []*ast.File) {
	for _, file := range files {
		dst.Decls = append(dst.Decls, file.Decls...)
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
	for _, file := range files {
		for _, imp := range file.Imports {
			var impPath, _ = strconv.Unquote(imp.Path.Value)
			switch {
			case imp.Name != nil:
				var name = imp.Name.String()
				astutil.AddNamedImport(fset, file, name, impPath)
			default:
				astutil.AddImport(fset, file, impPath)
			}
		}
	}
}
