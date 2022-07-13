package utils

import (
	"strings"

	"golang.org/x/tools/go/packages"
)

func LoadPackage(source string) (*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes,
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

	return pkgs[0], nil
}
