// +build mage

package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ninedraft/enum2go/pkg/astx"

	"go.uber.org/multierr"
)

func Generate() error {
	return multierr.Combine(
		embeddedCast(),
		embeddedUsage(),
	)
}

func embeddedCast() error {
	var fset = token.NewFileSet()
	var castPkgDir = filepath.Join("pkg", "generator", "enumcast")
	log.Printf("parsing template package %q", castPkgDir)

	var pkgs, errParseDir = parser.ParseDir(fset, castPkgDir, filterCastFiles, parser.AllErrors)
	if errParseDir != nil {
		return fmt.Errorf("unable to parse template package %q: %v", castPkgDir, errParseDir)
	}
	var pkg, errSelectPackage = astx.SelectPkg(pkgs)
	if errSelectPackage != nil {
		return fmt.Errorf("unable to select template package: %v", errSelectPackage)
	}
	var files = astx.PkgFiles(pkg)
	var c = &ast.File{
		Name: ast.NewIdent("cast"),
	}
	astx.MergeFiles(c, files)

	var data = &bytes.Buffer{}
	_ = format.Node(data, fset, c)

	var result = bytes.NewBufferString(
		"// code generate by mage script. DO NOT EDIT.\n\n" +
			"package static\n\n" +
			fmt.Sprintf("import %q\n", "github.com/ninedraft/enum2go/pkg/generator/cast") +
			fmt.Sprintf("import %q\n\n", "bytes"),
	)
	_, _ = fmt.Fprintf(result, "var Cast = cast.MustFromRe(%q, bytes.NewReader(%#v))\n\n", "cast", data.Bytes())
	return ioutil.WriteFile("pkg/static/cast.go", result.Bytes(), 0755)
}

func filterCastFiles(file os.FileInfo) bool {
	var isNotTest = !strings.HasSuffix(file.Name(), "_test.go")
	var isNotManifest = file.Name() != "pkg.go"
	return isNotTest && isNotManifest
}

func embeddedUsage() error {
	var usageFile = filepath.Join("doc", "usage.md")
	var usage, errRead = ioutil.ReadFile(usageFile)
	if errRead != nil {
		return errRead
	}
	usage = bytes.TrimSpace(usage)

	var result = bytes.NewBufferString(
		"// code generate by mage script. DO NOT EDIT.\n\n" +
			"package static\n\n",
	)
	_, _ = fmt.Fprintf(result, "const Usage = %q\n\n", usage)
	return ioutil.WriteFile("pkg/static/usage.go", result.Bytes(), 0755)
}
