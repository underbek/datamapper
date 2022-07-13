package parser

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/underbek/datamapper/models"
	"github.com/underbek/datamapper/utils"
	"golang.org/x/exp/maps"
	"golang.org/x/tools/go/packages"
)

type Packages = map[string]*packages.Package

type context struct {
	packages   Packages
	currentPkg *packages.Package
	fset       *token.FileSet
}

func ParseConversionFunctions(source string) (models.Functions, error) {
	ctx := context{}

	ctx.fset = token.NewFileSet()

	node, err := parser.ParseFile(ctx.fset, source, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	funcs := make(models.Functions)

	ctx.currentPkg, err = utils.LoadPackage(source)
	if err != nil {
		return nil, err
	}

	ctx.packages, err = loadPackagesByImports(node.Imports)
	if err != nil {
		return nil, err
	}

	for _, f := range node.Decls {
		funcD, ok := f.(*ast.FuncDecl)
		if !ok {
			continue
		}

		currentFuncs, err := parseFunction(ctx, funcD)
		if err != nil {
			return nil, err
		}

		maps.Copy(funcs, currentFuncs)
	}

	return funcs, nil
}
