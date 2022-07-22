package utils

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

var ErrParseError = errors.New("parse error")

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

	absSourcePath, err := filepath.Abs(source)
	if err != nil {
		return nil, err
	}

	pkgs, err := packages.Load(cfg, absSourcePath)
	if err != nil {
		return nil, err
	}

	pkg := pkgs[0]

	if len(pkg.Errors) != 0 {
		errs := make([]string, 0, len(pkg.Errors))
		for _, err := range pkg.Errors {
			if !strings.Contains(err.Error(), fmt.Sprintf("no Go files in %s", absSourcePath)) {
				errs = append(errs, err.Error())
			}
		}
		if len(errs) != 0 {
			return nil, fmt.Errorf("%w: %s", ErrParseError, strings.Join(errs, "; "))
		}
	}

	return pkg, nil
}
