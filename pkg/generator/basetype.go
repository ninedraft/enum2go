package generator

import (
	"fmt"
	"go/ast"
)

func typeKind(t *ast.Ident) (enumKind, bool) {
	if !isBaseTypeAllowed(t) {
		return enumKindEnum.Empty(), false
	}
	var _, isInt = intTypes[t.Name]
	var _, isString = stringTypes[t.Name]
	switch {
	case isInt:
		return enumKindEnum.Int(), true
	case isString:
		return enumKindEnum.String(), true
	}
	return enumKindEnum.Empty(), false
}

func isBaseTypeAllowed(name *ast.Ident) bool {
	if name == nil {
		return false
	}
	var _, ok = baseTypeWhiteList[name.Name]
	return ok
}

var (
	baseTypeWhiteList = map[string]struct{}{}
	baseTypes         []string
	intTypes          = map[string]struct{}{}
	stringTypes       = map[string]struct{}{}
)

func init() {
	var allow = func(names ...string) {
		for _, name := range names {
			baseTypeWhiteList[name] = struct{}{}
			baseTypes = append(baseTypes, names...)
		}
	}

	// int types

	intTypes["byte"] = struct{}{}
	allow("byte")

	intTypes["int"] = struct{}{}
	allow("int")

	intTypes["rune"] = struct{}{}
	allow("rune")

	for _, integer := range []string{"int", "uint"} {
		allow(integer)
		for _, c := range []int{8, 16, 32, 64} {
			var t = fmt.Sprintf("%s%d", integer, c)
			allow(t)
			intTypes[t] = struct{}{}
		}
	}

	// string types
	stringTypes["string"] = struct{}{}
	allow(
		"string",
	)
}
