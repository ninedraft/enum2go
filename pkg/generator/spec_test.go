package generator

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ninedraft/enum2go/pkg/astx"
)

func TestSpec(test *testing.T) {
	var fset = token.NewFileSet()
	var filter = func(file os.FileInfo) bool {
		return true
	}
	var pkgs, errPars = parser.ParseDir(fset, "enumcast", filter, parser.AllErrors)
	if errPars != nil {
		test.Fatal(errPars)
	}
	var pkg *ast.Package
	for _, p := range pkgs {
		pkg = p
	}
	var files = make([]*ast.File, 0, len(pkg.Files))
	for _, file := range pkg.Files {
		files = append(files, file)
	}
	var cast = &ast.File{}
	astx.MergeFiles(cast, files)

	var spec = enumSpec{
		Package:  "result",
		Type:     "Kind",
		Kind:     enumKindEnum.Int(),
		Format:   enumFormatEnum.Kebab(),
		BaseType: ast.NewIdent("int"),
		Names: []*ast.Ident{
			ast.NewIdent("AaA"),
			ast.NewIdent("bbA"),
			ast.NewIdent("C"),
		},
	}
	var gen = &Config{}
	spec.Pour(gen, cast)
	if err := format.Node(ioutil.Discard, fset, cast); err != nil {
		test.Fatal(err)
	}
}
