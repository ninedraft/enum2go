package astx

import (
	"go/ast"
)

// FieldMatcher is a configurable matcher for field lists.
type FieldMatcher struct {
	FieldNames func(fields []*ast.Field) []string
	NameEq     func(a, b string) bool
}

// Match returns true if any field matches the given name.
func (fm FieldMatcher) Match(name string, fields []*ast.Field) bool {
	var eq = fm.eq()
	var names = fm.names()
	for _, method := range names(fields) {
		if eq(method, name) {
			return true
		}
	}
	return false
}

// Filter returns fields which match the given name.
func (fm FieldMatcher) Filter(name string, fields []*ast.Field) []*ast.Field {
	var eq = fm.eq()
	var names = fm.names()
	var filtered []*ast.Field
	for i, method := range names(fields) {
		if eq(method, name) {
			filtered = append(filtered, fields[i])
		}
	}
	return filtered
}

// Select returns field, which matches the given name or nil.
func (fm FieldMatcher) Select(name string, fields []*ast.Field) *ast.Field {
	var eq = fm.eq()
	var names = fm.names()
	for i, method := range names(fields) {
		if eq(method, name) {
			return fields[i]
		}
	}
	return nil
}

func (fm FieldMatcher) eq() func(a, b string) bool {
	if fm.NameEq != nil {
		return fm.NameEq
	}
	return func(a, b string) bool { return a == b }
}

func (fm FieldMatcher) names() func(fields []*ast.Field) []string {
	if fm.FieldNames != nil {
		return fm.FieldNames
	}
	return FieldsNames
}

// FieldsNames returns all field names flattened into a single list.
func FieldsNames(fields []*ast.Field) []string {
	var names = make([]string, 0, len(fields))
	for _, field := range fields {
		for _, name := range field.Names {
			names = append(names, name.String())
		}
	}
	return names
}

// FieldsNamesEmbed returns all field names or names of embedded types (if not anonymous).
func FieldsNamesEmbed(fields []*ast.Field) []string {
	var names = make([]string, 0, len(fields))
	for _, field := range fields {
		for _, name := range field.Names {
			names = append(names, name.String())
		}
		if len(field.Names) > 0 {
			continue
		}
		var name, isNamed = nodeName(field.Type)
		if isNamed {
			names = append(names, name)
		}
	}
	return names
}

func nodeName(node ast.Node) (_ string, ok bool) {
	switch node := node.(type) {
	case *ast.Ident:
		return node.Name, true
	case *ast.StarExpr:
		return nodeName(node.X)
	default:
		return "", false
	}
}
