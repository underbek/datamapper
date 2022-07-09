package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/underbek/datamapper/models"
	"golang.org/x/tools/go/packages"
)

type context struct {
	packages map[string]*packages.Package
}

func ParseConversionFunctions(source string) (map[models.ConversionFunctionKey]models.ConversionFunction, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, source, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	funcs := make(map[models.ConversionFunctionKey]models.ConversionFunction)

	ctx := context{}

	ctx.packages, err = loadPackagesByImports(node.Imports)
	if err != nil {
		return nil, err
	}

	for _, f := range node.Decls {
		funcD, ok := f.(*ast.FuncDecl) //ast.GenDecl
		if !ok {
			continue
		}

		if funcD.Type.Params.NumFields() != 1 {
			continue
		}

		if funcD.Type.Results.NumFields() != 1 {
			continue
		}

		if funcD.Type.Params.List[0].Type == nil {
			continue
		}

		fromType, ok := funcD.Type.Params.List[0].Type.(*ast.Ident)
		if !ok {
			continue
		}

		if funcD.Type.Results.List[0].Type == nil {
			continue
		}

		toType, ok := funcD.Type.Results.List[0].Type.(*ast.Ident)
		if !ok {
			continue
		}

		if funcD.Type.TypeParams.NumFields() == 0 {
			funcs[models.ConversionFunctionKey{
				FromType: fromType.Name,
				ToType:   toType.Name,
			}] = models.ConversionFunction{
				Name: funcD.Name.Name,
			}

			continue
		}

		typeParams, err := parseTypeParams(ctx, funcD.Type.TypeParams.List)
		if err != nil {
			return nil, fmt.Errorf("parse func %s failed: %w", funcD.Name.Name, err)
		}

		fromTypes, ok := typeParams[fromType.Name]
		if !ok {
			fromTypes = []string{fromType.Name}
		}

		toTypes, ok := typeParams[toType.Name]
		if !ok {
			toTypes = []string{toType.Name}
		}

		for _, fromType := range fromTypes {
			for _, toType := range toTypes {
				funcs[models.ConversionFunctionKey{
					FromType: fromType,
					ToType:   toType,
				}] = models.ConversionFunction{
					Name: funcD.Name.Name,
				}
			}
		}
	}

	return funcs, nil
}
