package astx

import (
	"fmt"
	"go/ast"
	"sort"
	"strings"
)

// SelectPkg searches for non-test packages int the package set.
// Returns an error if no non-test package or multiple packages are found.
func SelectPkg(pkgs map[string]*ast.Package) (*ast.Package, error) {
	var isNonTests = func(name string) bool { return !strings.HasSuffix(name, "_test") }
	var nonTests []*ast.Package
	for name, pkg := range pkgs {
		if isNonTests(name) {
			nonTests = append(nonTests, pkg)
		}
	}
	if len(nonTests) > 1 {
		return nil, fmt.Errorf("found multiple non-tests packages: %v", PkgNames(pkgs))
	}
	return nonTests[0], nil
}

// PkgNames returns a sorted list of package names.
func PkgNames(pkgs map[string]*ast.Package) []string {
	var names = make([]string, 0, len(pkgs))
	for name := range pkgs {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// PkgFiles returns a list of package files sorted by filenames.
func PkgFiles(pkg *ast.Package) []*ast.File {
	var files = make([]*ast.File, 0, len(pkg.Files))
	var names = make(map[*ast.File]string, len(pkg.Files))
	for filename, file := range pkg.Files {
		files = append(files, file)
		names[file] = filename
	}
	sort.Slice(files, func(i, j int) bool {
		var fi, fj = files[i], files[j]
		return names[fi] < names[fj]
	})
	return files
}
