package generator

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ninedraft/enum2go/pkg/astx"
	"github.com/ninedraft/enum2go/pkg/static"
)

// Config describes some
type Config struct {
	Dir        string
	TargetFile string
}

// Run generates enum definitions using config data and specs in user code.
func Run(cfg *Config) error {
	var fset = token.NewFileSet()
	var pkgs, errParse = parser.ParseDir(fset, cfg.Dir, fileFilter, parser.AllErrors)
	if errParse != nil {
		return errParse
	}
	var pkg, errSelect = astx.SelectPkg(pkgs)
	if errSelect != nil {
		return errSelect
	}
	var bind, errBind = bindSpecs(fset, pkg)
	if errBind != nil {
		return errBind
	}
	var results = make([]*ast.File, 0, len(bind.specs))
	for _, spec := range bind.specs {
		var cast, _ = static.Cast.AST()
		spec.Pour(cfg, cast)
		results = append(results, cast)
	}
	var dst = &ast.File{
		Name: ast.NewIdent(pkg.Name),
	}
	astx.MergeFiles(dst, results)
	var generated = bytes.NewBufferString("// generated by enum2go. DO NOT EDIT.\n\n")
	if err := format.Node(generated, fset, dst); err != nil {
		return err
	}
	var target = filepath.Join(cfg.Dir, cfg.TargetFile)
	return ioutil.WriteFile(target, generated.Bytes(), 0600)
}

func fileFilter(os.FileInfo) bool { return true }
