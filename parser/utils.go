package parser

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

func loadPackagesByImports(imports []*ast.ImportSpec) (map[string]*packages.Package, error) {
	paths := make([]string, 0, len(imports))
	for _, spec := range imports {
		pkgPath := strings.Trim(spec.Path.Value, "\"")
		paths = append(paths, pkgPath)
	}

	cfg := &packages.Config{
		Mode: packages.NeedTypes,
	}

	pkgs, err := packages.Load(cfg, paths...)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*packages.Package)

	for _, pkg := range pkgs {
		result[pkg.Types.Name()] = pkg
	}

	return result, nil
}

func parseTypeParams(ctx context, fields []*ast.Field) (map[string][]string, error) {
	typeParams := make(map[string][]string)
	for _, field := range fields {
		if field.Type == nil {
			continue
		}

		fieldTypes, err := parseFiledType(ctx, field.Type)
		if err != nil {
			return nil, err
		}

		if len(fieldTypes) == 0 {
			continue
		}

		for _, name := range field.Names {
			typeParams[name.Name] = append(typeParams[name.Name], fieldTypes...)
		}
	}

	return typeParams, nil
}

func parseFiledType(ctx context, field ast.Expr) ([]string, error) {
	switch t := field.(type) {
	case *ast.Ident:
		return []string{t.Name}, nil
	case *ast.BinaryExpr:
		x, err := parseFiledType(ctx, t.X)
		if err != nil {
			return nil, err
		}
		y, err := parseFiledType(ctx, t.Y)
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

func getOriginalTypes(ctx context, packageName, typeName string) ([]string, error) {
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

func parseTypeFromPackage(t types.Type) ([]string, error) {
	switch t := t.(type) {
	case *types.Named:
		return parseTypeFromPackage(t.Underlying())
	case *types.Interface:
		n := t.NumEmbeddeds()
		res := make([]string, 0, n)
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
		res := make([]string, 0, n)
		for i := 0; i < n; i++ {
			names, err := parseTypeFromPackage(t.Term(i).Type())
			if err != nil {
				return nil, err
			}

			res = append(res, names...)
		}
		return res, nil
	case *types.Basic, *types.Struct:
		return []string{t.String()}, nil

	default:
		return nil, fmt.Errorf("undefined type %s", t.String())
	}
}
