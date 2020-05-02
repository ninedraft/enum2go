package embed

import (
	"bufio"
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strconv"

	"github.com/ninedraft/enum2go/pkg/astx"
)

func generateEmbed(dst *ast.File, key string, data []byte) {
	var literal = &bytes.Buffer{}
	literal.Grow(3 * len(data) / 2)
	_ = encodeStringLiteral(literal, data)
	var spec = &ast.ValueSpec{
		Names: []*ast.Ident{ast.NewIdent(key)},
		Values: []ast.Expr{
			&ast.BasicLit{Kind: token.STRING, Value: literal.String()},
		},
	}
	dst.Decls = append(dst.Decls, &ast.GenDecl{
		Tok:   token.CONST,
		Specs: []ast.Spec{spec},
	})
}

func encodeStringLiteral(dst io.Writer, data []byte) error {
	var re = bytes.NewReader(data)
	var scanner = bufio.NewScanner(re)
	var buf []byte
	var write = func(buf []byte) error {
		var _, err = dst.Write(buf)
		return err
	}
	var delim []byte
	for i := 0; scanner.Scan(); i++ {
		buf = strconv.AppendQuote(buf[:0], scanner.Text())
		var errWrite = anyErr(
			write(delim),
			write(buf),
		)
		if errWrite != nil {
			return errWrite
		}
		if i == 0 {
			delim = []byte("+\n")
		}
	}
	return nil
}

type fileFilter = func(file os.FileInfo) bool

func loadAST(dst io.Writer, dir string, filters ...fileFilter) error {
	var fset = token.NewFileSet()
	var pkgs, errParse = parser.ParseDir(fset, dir, fileFilters(filters), parser.AllErrors)
	if errParse != nil {
		return errParse
	}
	var pkg, errSelect = astx.SelectPkg(pkgs)
	if errSelect != nil {
		return errSelect
	}
	var files = astx.PkgFiles(pkg)
	var merged = &ast.File{
		Name: ast.NewIdent(pkg.Name),
	}
	astx.MergeFiles(merged, files)
	return format.Node(dst, fset, merged)
}

func fileFilters(filters []fileFilter) fileFilter {
	return func(file os.FileInfo) bool {
		for _, filter := range filters {
			if !filter(file) {
				return false
			}
		}
		return true
	}
}

func anyErr(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
