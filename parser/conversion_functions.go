package parser

import (
	"errors"
	"fmt"
	"go/build"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"github.com/underbek/datamapper/logger"
	"github.com/underbek/datamapper/models"
	"github.com/underbek/datamapper/utils"
	"golang.org/x/tools/go/packages"
)

var (
	ErrNotFoundType  = errors.New("not found type error")
	ErrNotFoundSign  = errors.New("not found signature error")
	ErrUndefinedType = errors.New("undefined type error")
)

func ParseConversionFunctionsByPackage(lg logger.Logger, source string) (models.Functions, error) {
	_, err := os.Stat(source)
	if err == nil {
		return ParseConversionFunctions(lg, source)
	}

	if !os.IsNotExist(err) {
		return nil, err
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	p, err := build.Import(source, wd, build.FindOnly)
	if err != nil {
		return nil, err
	}

	return ParseConversionFunctions(lg, p.Dir)
}
func ParseConversionFunctions(lg logger.Logger, source string) (models.Functions, error) {

	absSourcePath, err := filepath.Abs(source)
	if err != nil {
		return nil, err
	}

	pkg, err := utils.LoadPackage(lg, source)
	if err != nil {
		return nil, err
	}

	if pkg.Types == nil {
		return nil, fmt.Errorf("%w: package %s hasn't type", ErrNotFoundType, pkg.Name)
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

		for key, function := range currentFuncs {
			funcs[key] = function
		}
	}

	return funcs, nil
}

func parseFunction(pkg *packages.Package, f *types.Func) (models.Functions, error) {
	signature, ok := f.Type().(*types.Signature)
	if !ok {
		return nil, fmt.Errorf("%w: function %s hasn't signature", ErrNotFoundSign, f.Name())
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
	if signature.Results().Len() == 2 { //nolint:gomnd
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
				FromType: fromType.Type,
				ToType:   toType.Type,
			}

			cv := models.ConversionFunction{
				Name: f.Name(),
				Package: models.Package{
					Name: pkg.Name,
					Path: pkg.PkgPath,
				},
				FromType:  fromType.Type,
				ToType:    toType.Type,
				TypeParam: getTypeParam(fromType.generic, toType.generic),
				WithError: withError,
			}

			funcs[key] = cv
		}
	}

	return funcs, nil
}
