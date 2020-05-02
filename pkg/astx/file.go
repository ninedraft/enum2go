package astx

import "go/ast"

// FileAppend appends provided declarations to the file.
func FileAppend(file *ast.File, decl ...ast.Decl) {
	file.Decls = append(file.Decls, decl...)
}
