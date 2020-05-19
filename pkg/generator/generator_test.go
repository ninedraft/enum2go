package generator_test

import (
	"bytes"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/ninedraft/enum2go/pkg/astx"
	"github.com/ninedraft/enum2go/pkg/generator"
	"gopkg.in/yaml.v2"
)

func TestGenerator(test *testing.T) {
	var issues = loadIssues(test, "testdata")
	if len(issues) == 0 {
		test.Fatalf("no issues to test")
	}
	for _, issue := range issues {
		issue := issue
		test.Run(issue.Name+" "+issue.Name, func(test *testing.T) {
			var tio = newTestFileIO()
			var errGenerate = generator.Run(&generator.Config{
				Dir:        issue.Dir,
				TargetFile: "test_file.go",
				FileIO:     tio,
			})
			if errGenerate != nil {
				test.Errorf("unable to generate: %v", errGenerate)
				return
			}
			validateGeneratedCode(test, issue, tio.result)
		})
	}
}

type Issue struct {
	Name     string
	Dir      string
	Expected Expected
}

type Expected struct {
	Methods []string
}

func loadExpected(test *testing.T, tdata, issueID string) Issue {
	var fpath = filepath.Join(tdata, issueID, "issue.yml")
	// nolint:gosec // no path injection is possible
	var data, errRead = ioutil.ReadFile(fpath)
	if errRead != nil {
		test.Fatalf("reading expected data: %v", errRead)
	}
	var issue Issue
	var errUnmarshal = yaml.Unmarshal(data, &issue)
	if errUnmarshal != nil {
		test.Fatalf("decoding expected data: %v", errUnmarshal)
	}
	issue.Dir = filepath.Join(tdata, issueID)
	return issue
}

var isIssueName = regexp.MustCompile(`^#[0-9]+$`).MatchString

func loadIssues(test *testing.T, tdata string) []Issue {
	var files, errReadDir = ioutil.ReadDir(tdata)
	if errReadDir != nil {
		test.Fatalf("parsing issues: %v", errReadDir)
	}
	var issues = make([]Issue, 0, len(files))
	for _, file := range files {
		if !file.IsDir() || !isIssueName(file.Name()) {
			continue
		}
		issues = append(issues, loadExpected(test, tdata, file.Name()))
	}
	return issues
}

type testFileIO struct {
	result *bytes.Buffer
}

func newTestFileIO() *testFileIO {
	return &testFileIO{
		result: &bytes.Buffer{},
	}
}

func (nfio *testFileIO) ParsePkg(fset *token.FileSet, dir string) (pkgName string, _ []*ast.File, _ error) {
	var pkgs, errParse = parser.ParseDir(fset, dir, nfio.fileFilter, parser.AllErrors)
	if errParse != nil {
		return "", nil, errParse
	}
	var pkg, errSelect = astx.SelectPkg(pkgs)
	if errSelect != nil {
		return "", nil, errSelect
	}
	return pkg.Name, astx.PkgFiles(pkg), nil
}

func (nfio *testFileIO) WriteFile(_ string, data io.Reader) error {
	nfio.result.Reset()
	var _, err = io.Copy(nfio.result, data)
	return err
}

func (*testFileIO) fileFilter(os.FileInfo) bool { return true }

func validateGeneratedCode(test *testing.T, issue Issue, data io.Reader) {
	var fset = token.NewFileSet()
	var f, errParse = parser.ParseFile(fset, "hello.go", data, parser.AllErrors)
	if errParse != nil {
		test.Errorf("parsing test result: %v", errParse)
		return
	}

	var conf = types.Config{Importer: importer.Default()}

	var pkg, errCheck = conf.Check("test", fset, []*ast.File{f}, nil)
	if errParse != nil {
		test.Errorf("type check: %v", errCheck)
		return
	}

	var object = pkg.Scope().Lookup("_FooEnum")
	if object == nil || object.Type() == nil {
		test.Errorf("generated enum type not found. Available names: %v", pkg.Scope().Names())
		return
	}
	var named, isNamed = object.Type().(*types.Named)
	if !isNamed {
		test.Errorf("expected Foo to be a type declaration")
		return
	}
	if len(issue.Expected.Methods) > named.NumMethods() {
		test.Errorf("expected %d methods, got %d", len(issue.Expected.Methods), named.NumMethods())
	}
	issue.Expected.checkMethods(test, named)
}

func (expected Expected) checkMethods(test *testing.T, t *types.Named) {
	var methods = strSet(expected.Methods)
	for i := 0; i < t.NumMethods(); i++ {
		var method = t.Method(i)
		if methods[method.Name()] {
			test.Logf("✅ method %q", method.Name())
		}
		delete(methods, method.Name())
	}
	for method := range methods {
		test.Errorf("❌ method %s", method)
	}
}

func strSet(items []string) map[string]bool {
	var set = make(map[string]bool, len(items))
	for _, item := range items {
		set[item] = true
	}
	return set
}
