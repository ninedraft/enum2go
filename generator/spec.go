package generator

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/stoewer/go-strcase"
)

type enumSpec struct {
	Package  string
	Type     string
	Kind     enumKind
	Format   enumFormat
	BaseType *ast.Ident
	Names    []*ast.Ident
}

type (
	enumKind int

	_ struct {
		undefined, string, int enumKind
	}
)

type (
	enumFormat int

	_ struct {
		_, strict, snake, kebab enumFormat
	}
)

// ----------------------------------------------------------------

const (
	kindUndefined enumKind = iota
	kindString
	kindInt
)

const (
	strict enumFormat = iota
	snake
	kebab
)

// ----------------------------------------------------------------

func (enum *enumSpec) Pour(generator *Generator, cast *ast.File) {
	cast.Name = ast.NewIdent(enum.Package)
	var values = enum.allValuesList().Elts
	for i, name := range enum.Names {
		cast.Decls = append(cast.Decls, enum.factory(generator, name, values[i]))
	}
	enum.patchAST(generator, cast)
}

func (enum *enumSpec) factory(generator *Generator, name *ast.Ident, value ast.Expr) *ast.FuncDecl {
	var gadget = ast.NewIdent("_" + generator.typePlaceholder + "Enum")
	return &ast.FuncDecl{
		Name: name,
		Recv: &ast.FieldList{List: []*ast.Field{{Type: gadget}}},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{Type: ast.NewIdent(enum.Type)},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{Results: []ast.Expr{value}},
			},
		},
	}
}

func (enum *enumSpec) patchAST(generator *Generator, node ast.Node) {
	astutil.Apply(node, func(cursor *astutil.Cursor) bool {
		switch node := cursor.Node().(type) {
		case *ast.Ident:
			node.Name = strings.ReplaceAll(node.Name, generator.typePlaceholder, enum.Type)
		case *ast.FuncDecl:
			var name = node.Name.String()
			if methods[name] {
				node.Body = &ast.BlockStmt{
					List: enum.castMethod(name),
				}
			}
		}
		return true
	}, nil)
}

/*
	String
	IsValid
	AllValues
	AllNames
	Parse
*/

func (enum *enumSpec) castMethod(name string) []ast.Stmt {
	switch name {
	case "String":
		return enum.methodString()
	case "IsValid":
		return enum.methodIsValid()
	case "AllValues":
		return enum.methodAllValues()
	case "AllNames":
		return enum.methodAllNames()
	case "Parse":
		return enum.methodParse()
	}
	return nil
}

func (enum *enumSpec) methodIsValid() []ast.Stmt {
	var input = ast.NewIdent("v")
	var probe ast.Expr
	switch enum.Kind {
	case kindInt:
		var first = &ast.BasicLit{Kind: token.INT, Value: "1"}
		var last = &ast.BasicLit{Kind: token.INT, Value: strconv.Itoa(len(enum.Names) - 1)}
		probe = and(
			leq(first, input),
			leq(input, last),
		)
	case kindString:
		var strs = enum.allNamesList().Elts
		var first = strs[0]
		probe = eq(first, input)
		for _, str := range strs[1:] {
			probe = or(
				probe,
				eq(str, input),
			)
		}
	}
	return []ast.Stmt{
		&ast.ReturnStmt{
			Results: []ast.Expr{probe},
		},
	}
}

func (enum *enumSpec) methodString() []ast.Stmt {
	var input = ast.NewIdent("v")
	var sw = &ast.SwitchStmt{
		Tag:  input,
		Body: &ast.BlockStmt{},
	}
	var values = enum.allValuesList().Elts
	var strs = enum.allNamesList().Elts
	for i, value := range values {
		sw.Body.List = append(sw.Body.List,
			&ast.CaseClause{
				List: []ast.Expr{value},
				Body: []ast.Stmt{
					&ast.ReturnStmt{Results: []ast.Expr{strs[i]}},
				},
			})
	}

	var baseValue = &ast.CallExpr{
		Fun:  enum.BaseType,
		Args: []ast.Expr{input},
	}
	var err = makeStringFormat("unexpected value %v. Valid values: %v", baseValue, enum.allNamesList())
	sw.Body.List = append(sw.Body.List,
		&ast.CaseClause{
			Body: []ast.Stmt{
				&ast.ReturnStmt{Results: []ast.Expr{err}},
			},
		})
	return []ast.Stmt{sw}
}

func (enum *enumSpec) methodAllValues() []ast.Stmt {
	return []ast.Stmt{
		&ast.ReturnStmt{
			Results: []ast.Expr{enum.allValuesList()},
		},
	}
}

func (enum *enumSpec) allValuesList() *ast.CompositeLit {
	var list = &ast.CompositeLit{
		Type: &ast.ArrayType{Elt: ast.NewIdent(enum.Type)},
	}
	var converter = enum.valueConverter()
	for i := range enum.Names {
		list.Elts = append(list.Elts, converter(i))
	}
	return list
}

func (enum *enumSpec) methodAllNames() []ast.Stmt {
	return []ast.Stmt{
		&ast.ReturnStmt{
			Results: []ast.Expr{enum.allNamesList()},
		},
	}
}

func (enum *enumSpec) allNamesList() *ast.CompositeLit {
	var list = &ast.CompositeLit{
		Type: &ast.ArrayType{Elt: ast.NewIdent("string")},
	}
	var stringer = enum.stringer()
	for i := range enum.Names {
		var str = stringer(i)
		list.Elts = append(list.Elts, &ast.BasicLit{
			Kind:  token.STRING,
			Value: strconv.Quote(str),
		})
	}
	return list
}

func (enum *enumSpec) methodParse() []ast.Stmt {
	var input = ast.NewIdent("str")
	var nilVal = ast.NewIdent("nil")
	var empty = ast.NewIdent("empty")
	var sw = &ast.SwitchStmt{
		Tag:  input,
		Body: &ast.BlockStmt{},
	}
	var values = enum.allValuesList().Elts
	var strs = enum.allNamesList().Elts
	for i, value := range values {
		sw.Body.List = append(sw.Body.List,
			&ast.CaseClause{
				List: []ast.Expr{strs[i]},
				Body: []ast.Stmt{
					&ast.ReturnStmt{Results: []ast.Expr{value, nilVal}},
				},
			})
	}
	var err = makeErrorFormat("unexpected value %q. Valid inputs: %v", input, enum.allNamesList())
	sw.Body.List = append(sw.Body.List,
		&ast.CaseClause{
			Body: []ast.Stmt{
				&ast.ReturnStmt{Results: []ast.Expr{empty, err}},
			},
		})
	return []ast.Stmt{
		&ast.DeclStmt{
			Decl: emptyDecl(empty, ast.NewIdent(enum.Type)),
		},
		sw,
	}
}

func (enum *enumSpec) valueConverter() func(i int) ast.Expr {
	var converter func(i int) ast.Expr
	var stringer = enum.stringer()
	switch enum.Kind {
	case kindString:
		converter = func(i int) ast.Expr {
			var name = stringer(i)
			var lit = strconv.Quote(name)
			return &ast.BasicLit{
				Kind:  token.STRING,
				Value: lit,
			}
		}
	case kindInt:
		converter = func(i int) ast.Expr {
			var lit = strconv.Itoa(i)
			return &ast.BasicLit{
				Kind:  token.INT,
				Value: lit,
			}
		}
	}
	return converter
}

func (enum *enumSpec) stringer() func(i int) string {
	var stringer = func(i int) string {
		var name = enum.Names[i]
		return name.Name
	}
	switch enum.Format {
	case kebab:
		stringer = func(i int) string {
			var name = enum.Names[i].Name
			return strcase.KebabCase(name)
		}
	case snake:
		stringer = func(i int) string {
			var name = enum.Names[i].Name
			return strcase.SnakeCase(name)
		}
	}
	return stringer
}

func makeErrorString(str ast.Expr) *ast.CallExpr {
	var newErr = &ast.SelectorExpr{
		X:   ast.NewIdent("errors"),
		Sel: ast.NewIdent("New"),
	}
	return &ast.CallExpr{
		Fun:  newErr,
		Args: []ast.Expr{str},
	}
}

func makeStringFormat(format string, args ...ast.Expr) *ast.CallExpr {
	var newErr = &ast.SelectorExpr{
		X:   ast.NewIdent("fmt"),
		Sel: ast.NewIdent("Sprintf"),
	}
	var formatLit = &ast.BasicLit{
		Kind:  token.STRING,
		Value: strconv.Quote(format),
	}
	return &ast.CallExpr{
		Fun:  newErr,
		Args: append([]ast.Expr{formatLit}, args...),
	}
}

func makeErrorFormat(format string, args ...ast.Expr) *ast.CallExpr {
	var newErr = &ast.SelectorExpr{
		X:   ast.NewIdent("fmt"),
		Sel: ast.NewIdent("Errorf"),
	}
	var formatLit = &ast.BasicLit{
		Kind:  token.STRING,
		Value: strconv.Quote(format),
	}
	return &ast.CallExpr{
		Fun:  newErr,
		Args: append([]ast.Expr{formatLit}, args...),
	}
}

var methods = strSet([]string{
	"AllValues",
	"AllNames",
	"String",
	"IsValid",
	"Parse",
})

func strSet(items []string) map[string]bool {
	var set = make(map[string]bool, len(items))
	for _, item := range items {
		set[item] = true
	}
	return set
}

func emptyDecl(name *ast.Ident, t ast.Expr) *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{name},
				Type:  t,
			},
		},
	}
}

func eq(left, right ast.Expr) *ast.BinaryExpr {
	return op(left, token.EQL, right)
}

func or(left, right ast.Expr) *ast.BinaryExpr {
	return op(left, token.OR, right)
}

func and(left, right ast.Expr) *ast.BinaryExpr {
	return op(left, token.LAND, right)
}

func leq(left, right ast.Expr) *ast.BinaryExpr {
	return op(left, token.LEQ, right)
}

func op(left ast.Expr, op token.Token, right ast.Expr) *ast.BinaryExpr {
	return &ast.BinaryExpr{
		Op: op,
		X:  left,
		Y:  right,
	}
}
