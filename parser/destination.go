package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/underbek/datamapper/logger"
	"github.com/underbek/datamapper/models"
	"github.com/underbek/datamapper/utils"
	"golang.org/x/tools/go/packages"
)

var ErrParseError = errors.New("parse error")

func ParseDestinationPackage(lg logger.Logger, destination string) (models.Package, error) {
	pkg, err := utils.LoadPackage(lg, destination)
	if err != nil {
		return models.Package{}, err
	}

	return generateModelPackage(pkg)
}

func generateModelPackage(pkg *packages.Package) (models.Package, error) {
	if pkg.Name != "" {
		return models.Package{
			Name: pkg.Name,
			Path: pkg.PkgPath,
		}, nil
	}

	if pkg.PkgPath == "" {
		return models.Package{}, fmt.Errorf("incorrect parsed destination package: %w", ErrParseError)
	}

	return models.Package{
		Name: getPackageNameByPath(pkg.PkgPath),
		Path: pkg.PkgPath,
	}, nil
}

func getPackageNameByPath(path string) string {
	names := strings.Split(path, "/")
	return names[len(names)-1]
}
