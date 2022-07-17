package parser

import (
	"fmt"
	"go/types"
	"path/filepath"
	"strings"

	"github.com/underbek/datamapper/models"
	"github.com/underbek/datamapper/utils"
	"golang.org/x/exp/maps"
	"golang.org/x/tools/go/packages"
)

func ParseConversionFunctions(source string) (models.Functions, error) {
	absSourcePath, err := filepath.Abs(source)
	if err != nil {
		return nil, err
	}

	pkg, err := utils.LoadPackage(source)
	if err != nil {
		return nil, err
	}

	if pkg.Types == nil {
		return nil, fmt.Errorf("package %s haven't type", pkg.Name)
	}

	funcs := make(models.Functions)

	names := pkg.Types.Scope().Names()
	for _, name := range names {
		obj := pkg.Types.Scope().Lookup(name)

		fset := pkg.Fset.Position(obj.Pos())
		if !strings.Contains(fset.Filename, absSourcePath) {
			continue
		}

		f, ok := obj.(*types.Func)
		if !ok {
			continue
		}

		if !f.Exported() {
			continue
		}

		currentFuncs, err := parseFunction(pkg, f)
		if err != nil {
			return nil, err
		}
		maps.Copy(funcs, currentFuncs)
	}

	return funcs, nil
}

func parseFunction(pkg *packages.Package, f *types.Func) (models.Functions, error) {
	signature, ok := f.Type().(*types.Signature)
	if !ok {
		return nil, fmt.Errorf("function %s hasn't signature", f.Name())
	}

	if signature.Params().Len() != 1 {
		return nil, nil
	}

	if signature.Results().Len() == 0 || signature.Results().Len() > 2 {
		return nil, nil
	}

	if signature.Params().At(0).Type() == nil {
		return nil, nil
	}

	if signature.Results().At(0).Type() == nil {
		return nil, nil
	}

	fromTypes, err := parseType(signature.Params().At(0).Type())
	if err != nil {
		return nil, err
	}

	toTypes, err := parseType(signature.Results().At(0).Type())
	if err != nil {
		return nil, err
	}

	withError := false
	if signature.Results().Len() == 2 {
		isError, err := isErrorType(signature.Results().At(1).Type())
		if err != nil {
			return nil, err
		}
		if !isError {
			return nil, nil
		}

		withError = true
	}

	funcs := make(models.Functions)

	for _, fromType := range fromTypes {
		for _, toType := range toTypes {
			key := models.ConversionFunctionKey{
				FromType: models.Type{
					Name:        fromType.Name,
					PackagePath: fromType.PackagePath,
				},
				ToType: models.Type{
					Name:        toType.Name,
					PackagePath: toType.PackagePath,
				},
			}

			cv := models.ConversionFunction{
				Name:        f.Name(),
				PackageName: pkg.Name,
				PackagePath: pkg.PkgPath,
				TypeParam:   getTypeParam(fromType.generic, toType.generic),
				WithError:   withError,
			}

			funcs[key] = cv
		}
	}

	return funcs, nil
}