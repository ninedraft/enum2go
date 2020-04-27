package generator

import (
	"fmt"
	"go/ast"
	"go/token"
	"sort"

	"go.uber.org/multierr"
	"golang.org/x/tools/go/ast/inspector"
)

type specBinder struct {
	files     []*ast.File
	filenames map[*ast.File]string

	// index of type enum declarations
	/*
		type Kind int
	*/
	index map[string]*typeAlias

	// index of enum specs
	/*

		type _ struct {
			Enum struct { A, B, C Kind }
		}

	*/
	specs []*enumSpec
}

// bindSpecs searches for enum types and enum specs and binds them together.
// It can return an error on every malformed enum spec.
func bindSpecs(fset *token.FileSet, pkg *ast.Package) (*specBinder, error) {
	var binder = &specBinder{
		files:     make([]*ast.File, 0, len(pkg.Files)),
		filenames: make(map[*ast.File]string, len(pkg.Files)),
		index:     make(map[string]*typeAlias, len(pkg.Files)),
		specs:     make([]*enumSpec, 0, len(pkg.Files)),
	}
	for filename, file := range pkg.Files {
		binder.filenames[file] = filename
		binder.files = append(binder.files, file)
	}

	var filter = []ast.Node{
		(*ast.TypeSpec)(nil),
	}

	var insp = inspector.New(binder.files)
	insp.Preorder(filter, func(node ast.Node) {
		var spec, isTypeSpec = node.(*ast.TypeSpec)
		if !isTypeSpec {
			return
		}

		var alias, isEnumAlias = binder.ParseAlias(spec)
		if isEnumAlias {
			binder.index[spec.Name.String()] = alias
			return
		}

		var enum, isEnumSpec = binder.ParseSpec(spec)
		if isEnumSpec {
			binder.specs = append(binder.specs, enum)
			return
		}
	})
	for _, spec := range binder.specs {
		spec.Package = pkg.Name
	}

	sort.Slice(binder.specs, func(i, j int) bool {
		return binder.specs[i].Type < binder.specs[j].Type
	})

	var errBind = binder.Validate(fset)
	if errBind != nil {
		return nil, errBind
	}
	return binder, nil
}

func (binder *specBinder) Validate(fset *token.FileSet) error {
	var errBind error
	var visited = make(map[string]*enumSpec, len(binder.specs))
	for _, spec := range binder.specs {
		var alias = binder.index[spec.Type]
		var position = fset.Position(spec.Pos)
		var err error
		switch {
		case alias == nil:
			err = fmt.Errorf("enum spec %q (%s) uses an invalid or not declared enum type. "+
				"Enum type must be declared as type wrapper over a certain base type."+
				"Valid base types: %v", spec.Type, position, baseTypes)
		case !alias.Valid:
			err = fmt.Errorf("enum spec %q (%s) uses an invalid enum type %s with base type %s. "+
				"Enum type must be declared as type wrapper over one the following base types: %v",
				spec.Type, position,
				alias.Type.Name, alias.BaseType, baseTypes)
		}
		multierr.AppendInto(&errBind, err)

		var previous, isVisited = visited[spec.Type]
		switch {
		case isVisited:
			var duplicatePosition = fset.Position(previous.Pos)
			var err = fmt.Errorf("enum spec %q (%s) is redeclared at %s", spec.Type, duplicatePosition, position)
			multierr.AppendInto(&errBind, err)
		default:
			visited[spec.Type] = spec
		}
	}
	return errBind
}

type typeAlias struct {
	Type     *ast.TypeSpec
	BaseType *ast.Ident
	Kind     enumKind
	Valid    bool
}

func (binder *specBinder) ParseAlias(node *ast.TypeSpec) (*typeAlias, bool) {
	var baseType *ast.Ident
	switch node := node.Type.(type) {
	case *ast.Ident:
		baseType = node
	default:
		return nil, false
	}
	var kind, ok = typeKind(baseType)
	return &typeAlias{
		Type:     node,
		BaseType: baseType,
		Kind:     kind,
		Valid:    ok,
	}, true
}

func (binder *specBinder) ParseSpec(node *ast.TypeSpec) (*enumSpec, bool) {
	var specContainer, isStructSpec = node.Type.(*ast.StructType)
	if !isStructSpec || node.Name.String() != "_" {
		return nil, false
	}
	var specField = getField(specContainer, "Enum")
	if specField == nil {
		return nil, false
	}
	var specStruct, isSpecStruct = specField.Type.(*ast.StructType)
	if !isSpecStruct {
		return nil, false
	}
	var spec *enumSpec
	for _, field := range specStruct.Fields.List {
		spec = &enumSpec{
			Pos:    specField.Pos(),
			Format: enumFormatEnum.Strict(),
		}
		var tName = typeName(field.Type)
		if tName == "" {
			return nil, false
		}
		spec.Type = tName
		spec.Names = append(spec.Names, field.Names...)
		break
	}
	return spec, spec != nil
}

func getField(spec *ast.StructType, name string) *ast.Field {
	for _, field := range spec.Fields.List {
		for _, n := range field.Names {
			if n.String() == name {
				return field
			}
		}
		if len(field.Names) != 0 {
			continue
		}
		if typeName(field.Type) == name {
			return field
		}
	}
	return nil
}

func typeName(node ast.Expr) string {
	switch node := node.(type) {
	case *ast.Ident:
		return node.String()
	case *ast.StarExpr:
		return typeName(node.X)
	default:
		return ""
	}
}
