package generator

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"testing"
)

func TestSpec(test *testing.T) {
	var fset = token.NewFileSet()
	var filter = func(file os.FileInfo) bool {
		return true
	}
	var pkgs, errPars = parser.ParseDir(fset, "cast", filter, parser.AllErrors)
	if errPars != nil {
		test.Fatal(errPars)
	}
	var pkg *ast.Package
	for _, p := range pkgs {
		pkg = p
	}
	var files []*ast.File
	for _, file := range pkg.Files {
		files = append(files, file)
	}
	var cast = &ast.File{}
	mergeFiles(cast, files)

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
	var gen = &Config{
		typePlaceholder: "Θ",
	}
	spec.Pour(gen, cast)
	if err := format.Node(ioutil.Discard, fset, cast); err != nil {
		test.Fatal(err)
	}
}
