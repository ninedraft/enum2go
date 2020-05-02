package astx

import (
	"go/ast"
	"go/token"
)

// ValBuilder provides utilities to build a single key-value expression.
type ValBuilder struct {
	name  *ast.Ident
	t     ast.Expr
	value ast.Expr
}

// Val creates a new new value builder.
func Val() *ValBuilder {
	return &ValBuilder{}
}

// WithName sets name for the value.
func (vb *ValBuilder) WithName(name string) *ValBuilder {
	vb.name = ast.NewIdent(name)
	return vb
}

// WithNameNode uses provided ident node as a name for the value.
func (vb *ValBuilder) WithNameNode(name *ast.Ident) *ValBuilder {
	vb.name = name
	return vb
}

// WithType sets type for the value.
func (vb *ValBuilder) WithType(t *ast.Ident) *ValBuilder {
	vb.t = t
	return vb
}

// WithValue sets values for the value spec.
func (vb *ValBuilder) WithValue(value ast.Expr) *ValBuilder {
	vb.value = value
	return vb
}

// Spec generates the key-value spec.
func (vb *ValBuilder) Spec() *ast.ValueSpec {
	var filled = false
	var spec = &ast.ValueSpec{}
	if vb.t != nil {
		filled = true
		spec.Type = vb.t
	}
	if vb.name != nil {
		filled = true
		spec.Names = []*ast.Ident{vb.name}
	}
	if vb.value != nil {
		filled = true
		spec.Values = []ast.Expr{vb.value}
	}
	if filled {
		return spec
	}
	return nil
}

// Decl generate a key-value declaration.
func (vb *ValBuilder) Decl() *ast.GenDecl {
	var decl = &ast.GenDecl{
		Tok: token.VAR,
	}
	var spec = vb.Spec()
	if spec != nil {
		decl.Specs = []ast.Spec{spec}
	}
	return decl
}

// Field generates a simple field type declaration from value.
func (vb *ValBuilder) Field() *ast.Field {
	var filled = false
	var field = &ast.Field{}
	if vb.t != nil {
		filled = true
		field.Type = vb.t
	}
	if vb.name != nil {
		filled = true
		field.Names = []*ast.Ident{vb.name}
	}
	if filled {
		return field
	}
	return nil
}
