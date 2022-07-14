package utils

import (
	"errors"
	"strings"

	"golang.org/x/tools/go/packages"
)

func LoadPackage(source string) (*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedDeps | packages.NeedImports,
	}

	index := strings.LastIndex(source, ".go")
	if index != -1 && index == len(source)-3 {
		index = strings.LastIndex(source, "/")
		if index == -1 {
			source = ""
		} else {
			source = source[:index]
		}
	}

	pkgs, err := packages.Load(cfg, source)
	if err != nil {
		return nil, err
	}

	pkg := pkgs[0]

	if len(pkg.Errors) != 0 {
		errs := make([]string, 0, len(pkg.Errors))
		for _, err := range pkg.Errors {
			errs = append(errs, err.Error())
		}
		return nil, errors.New(strings.Join(errs, "; "))
	}

	return pkg, nil
}
