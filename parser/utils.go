package parser

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"github.com/underbek/datamapper/models"
	"golang.org/x/tools/go/packages"
)

type Type struct {
	models.Type
	GenericName string
}

func parseFunction(ctx context, funcD *ast.FuncDecl) (models.Functions, error) {
	if funcD.Type.Params.NumFields() != 1 {
		return nil, nil
	}

	if funcD.Type.Results.NumFields() != 1 {
		return nil, nil
	}

	if funcD.Type.Params.List[0].Type == nil {
		return nil, nil
	}

	fromTypes, err := parseFiledTypes(ctx, funcD.Type.Params.List[0].Type)
	if err != nil {
		return nil, err
	}

	if funcD.Type.Results.List[0].Type == nil {
		return nil, nil
	}

	toTypes, err := parseFiledTypes(ctx, funcD.Type.Results.List[0].Type)
	if err != nil {
		return nil, err
	}

	funcs := make(models.Functions)

	for _, fromType := range fromTypes {
		for _, toType := range toTypes {
			key := models.ConversionFunctionKey{
				FromType: models.Type{
					Name:    fromType.Name,
					Package: fromType.Package,
				},
				ToType: models.Type{
					Name:    toType.Name,
					Package: toType.Package,
				},
			}

			cv := models.ConversionFunction{
				Name:        funcD.Name.Name,
				PackageName: ctx.currentPkg.Name,
				Import:      ctx.currentPkg.PkgPath,
				TypeParam:   getTypeParam(fromType.GenericName, toType.GenericName),
			}

			funcs[key] = cv
		}
	}

	return funcs, nil
}

func loadCurrentPackage(source string) (*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes,
	}

	index := strings.LastIndex(source, "/")
	if index == -1 {
		source = ""
	} else {
		source = source[:index]
	}

	pkgs, err := packages.Load(cfg, source)
	if err != nil {
		return nil, err
	}

	return pkgs[0], nil
}

func loadPackagesByImports(imports []*ast.ImportSpec) (Packages, error) {
	paths := make([]string, 0, len(imports))
	for _, spec := range imports {
		pkgPath := strings.Trim(spec.Path.Value, "\"")
		paths = append(paths, pkgPath)
	}

	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes,
	}

	pkgs, err := packages.Load(cfg, paths...)
	if err != nil {
		return nil, err
	}

	result := make(Packages)

	for _, pkg := range pkgs {
		result[pkg.Types.Name()] = pkg
	}

	return result, nil
}

func parseFiledTypes(ctx context, field ast.Expr) ([]Type, error) {
	switch t := field.(type) {
	case *ast.Ident:
		info := types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
		}

		err := types.CheckExpr(ctx.currentPkg.Fset, ctx.currentPkg.Types, t.Pos(), t, &info)
		if err != nil {
			if strings.Contains(err.Error(), "undeclared name:") {
				return parseGeneric(ctx, t)
			}
			return nil, err
		}

		tt, ok := info.Types[t]
		if !ok {
			return nil, fmt.Errorf("not found type %s from package %s", t.Name, ctx.currentPkg.Types.Name())
		}

		return parseTypeFromPackage(tt.Type)
	case *ast.BinaryExpr:
		x, err := parseFiledTypes(ctx, t.X)
		if err != nil {
			return nil, err
		}
		y, err := parseFiledTypes(ctx, t.Y)
		if err != nil {
			return nil, err
		}

		return append(x, y...), nil
	case *ast.SelectorExpr:
		pkg, ok := t.X.(*ast.Ident)
		if !ok {
			return nil, fmt.Errorf("cannot convert type.X to *ast.SelectorExpr")
		}

		return getOriginalTypes(ctx, pkg.Name, t.Sel.Name)

	}
	return nil, nil
}

func getOriginalTypes(ctx context, packageName, typeName string) ([]Type, error) {
	pkg, ok := ctx.packages[packageName]
	if !ok {
		return nil, fmt.Errorf("package %s not found", packageName)
	}

	if pkg.Types == nil {
		return nil, fmt.Errorf("types is nil for package %s", packageName)
	}

	obj := pkg.Types.Scope().Lookup(typeName)

	if obj == nil {
		return nil, fmt.Errorf("not found type %s from package %s", typeName, packageName)
	}

	return parseTypeFromPackage(obj.Type())
}

func parseTypeFromPackage(t types.Type) ([]Type, error) {
	switch t := t.(type) {
	case *types.Named:
		und := t.Underlying()
		if _, ok := und.(*types.Struct); ok {
			return []Type{{Type: models.Type{
				Name:    t.Obj().Name(),
				Package: t.Obj().Pkg().Name(),
			}}}, nil
		}
		return parseTypeFromPackage(t.Underlying())
	case *types.Interface:
		if t.String() == "any" {
			return []Type{{Type: models.Type{Name: "any"}}}, nil
		}
		n := t.NumEmbeddeds()
		res := make([]Type, 0, n)
		for i := 0; i < n; i++ {
			names, err := parseTypeFromPackage(t.EmbeddedType(i))
			if err != nil {
				return nil, err
			}

			res = append(res, names...)
		}
		return res, nil

	case *types.Union:
		n := t.Len()
		res := make([]Type, 0, n)
		for i := 0; i < n; i++ {
			names, err := parseTypeFromPackage(t.Term(i).Type())
			if err != nil {
				return nil, err
			}

			res = append(res, names...)
		}
		return res, nil
	case *types.Basic, *types.Struct:
		return []Type{{Type: models.Type{Name: t.String()}}}, nil
	default:
		return nil, fmt.Errorf("undefined type %s", t.String())
	}
}

// I don't know how it parse (((
func parseGeneric(ctx context, idn *ast.Ident) ([]Type, error) {
	if idn.Obj == nil {
		return []Type{{Type: models.Type{Name: idn.Name}}}, nil
	}

	if idn.Obj.Decl == nil {
		return []Type{{Type: models.Type{Name: idn.Name}}}, nil
	}

	field, ok := idn.Obj.Decl.(*ast.Field)
	if !ok {
		return []Type{{Type: models.Type{Name: idn.Name}}}, nil
	}

	fields, err := parseFiledTypes(ctx, field.Type)
	if err != nil {
		return nil, err
	}

	for i := range fields {
		fields[i].GenericName = idn.Name
	}

	return fields, nil
}

func getTypeParam(fromGenericName, toGenericName string) models.TypeParamType {
	if fromGenericName != "" && toGenericName != "" {
		return models.FromToTypeParam
	}

	if fromGenericName != "" {
		return models.FromTypeParam
	}

	if toGenericName != "" {
		return models.ToTypeParam
	}

	return models.NoTypeParam
}
