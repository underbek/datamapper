package utils

import (
	"path/filepath"
	"strings"

	"github.com/underbek/datamapper/logger"
	"golang.org/x/tools/go/packages"
)

var (
	packageCache = make(map[string]*packages.Package)
)

func LoadPackage(lg logger.Logger, source string) (*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedDeps | packages.NeedImports,
	}

	source = ClearFileName(source)

	absSourcePath, err := filepath.Abs(source)
	if err != nil {
		return nil, err
	}

	if pkg, ok := packageCache[absSourcePath]; ok {
		return pkg, nil
	}

	pkgs, err := packages.Load(cfg, absSourcePath)
	if err != nil {
		return nil, err
	}

	pkg := pkgs[0]

	for _, err := range pkg.Errors {
		lg.Warn(err)
	}

	packageCache[absSourcePath] = pkg

	return pkg, nil
}

func ClearFileName(source string) string {
	index := strings.LastIndex(source, ".go")
	if index != -1 && index == len(source)-3 {
		index = strings.LastIndex(source, "/")
		if index == -1 {
			return ""
		} else {
			return source[:index]
		}
	}

	return source
}
