package generator

import (
	"fmt"
	"go/ast"
	"go/token"
	"sort"

	"github.com/ninedraft/enum2go/pkg/astx"
	"go.uber.org/multierr"
	"golang.org/x/tools/go/ast/inspector"
)

type specBinder struct {
	files []*ast.File
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
func bindSpecs(fset *token.FileSet, pkg string, files []*ast.File) (*specBinder, error) {
	var binder = &specBinder{
		files: append(files[:0:0], files...),
		index: make(map[string]*typeAlias, len(files)),
		specs: make([]*enumSpec, 0, len(files)),
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
		spec.Package = pkg
	}

	sort.Slice(binder.specs, func(i, j int) bool {
		return binder.specs[i].Type < binder.specs[j].Type
	})

	var errBind = binder.bake(fset)
	if errBind != nil {
		return nil, errBind
	}
	return binder, nil
}

func (binder *specBinder) bake(fset *token.FileSet) error {
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
			spec.Kind = alias.Kind
			spec.BaseType = alias.BaseType
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
	var fields = specStruct.Fields.List
	if len(fields) == 0 {
		return nil, false
	}
	var t = baseTypeName(fields[0].Type)
	if t == "" {
		return nil, false
	}
	var spec = &enumSpec{
		Type:   t,
		Pos:    specField.Pos(),
		Format: enumFormatEnum.Strict(),
	}
	for _, field := range fields {
		var tName = baseTypeName(field.Type)
		if tName != t {
			continue
		}
		spec.Names = append(spec.Names, field.Names...)
	}
	binder.parseOptions(spec, fields)
	return spec, true
}

func (binder *specBinder) parseOptions(spec *enumSpec, cfg []*ast.Field) {
	var matcher = astx.FieldMatcher{}
	var opts = matcher.Select("opt", cfg)
	if opts == nil {
		return
	}
	var optData, isInterface = opts.Type.(*ast.InterfaceType)
	if !isInterface {
		return
	}
	switch {
	case matcher.Match("snake", optData.Methods.List):
		spec.Format = enumFormatEnum.Snake()
	case matcher.Match("kebab", optData.Methods.List):
		spec.Format = enumFormatEnum.Kebab()
	}
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
		if baseTypeName(field.Type) == name {
			return field
		}
	}
	return nil
}

func baseTypeName(node ast.Expr) string {
	switch node := node.(type) {
	case *ast.Ident:
		return node.String()
	default:
		return ""
	}
}
