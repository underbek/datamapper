package generator

import (
	"fmt"
	"strings"

	"github.com/underbek/datamapper/models"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/tools/go/packages"
)

func getTypeParams(cf models.ConversionFunction, fromType, toType models.Type) string {
	// TODO: add packages to type params
	switch cf.TypeParam {
	case models.ToTypeParam:
		return fmt.Sprintf("[%s]", toType.Name)
	case models.FromToTypeParam:
		return fmt.Sprintf("[%s,%s]", fromType.Name, toType.Name)
	default:
		return ""
	}
}

func fillConversions(fields []FieldsPair) []string {
	var res []string
	for _, field := range fields {
		res = append(res, field.Conversions...)
	}

	return res
}

func isReturnError(fields []FieldsPair) bool {
	for _, field := range fields {
		if field.WithError || field.PointerToValue {
			return true
		}
	}

	return false
}

func getFullStructName(model models.Struct, pkgPath string) string {
	if model.Package.Path != pkgPath {
		if model.Package.Alias != "" {
			return fmt.Sprintf("%s.%s", model.Package.Alias, model.Name)
		}
		return fmt.Sprintf("%s.%s", model.Package.Name, model.Name)
	}

	return model.Name
}

func getFullFieldName(filed models.Field, pkgPath string) string {
	// TODO: set package alias if names equals
	if filed.Type.Package.Name != "" && filed.Type.Package.Path != pkgPath {
		return fmt.Sprintf("%s.%s", filed.Type.Package.Name, filed.Type.Name)
	}

	return filed.Type.Name
}

func filterAndSortImports(currentPkg models.Package, imports []ImportType) []ImportType {
	currentPkgPath := currentPkg.Import()

	set := make(map[ImportType]struct{})
	for _, imp := range imports {
		if imp != "\"\"" && imp != currentPkgPath {
			set[imp] = struct{}{}
		}
	}

	res := maps.Keys(set)
	slices.Sort(res)

	return res
}

func generateConvertorName(from, to models.Struct, pkgPath string) string {
	structNameGenerator := func(model models.Struct, pkgPath string) string {
		name := model.Name

		if model.Package.Path == pkgPath {
			return name
		}

		pkgName := model.Package.Name
		if model.Package.Alias != "" {
			pkgName = model.Package.Alias
		}
		return cases.Title(language.Und, cases.NoLower).String(pkgName) + name
	}

	return fmt.Sprintf(
		"Convert%sTo%s",
		structNameGenerator(from, pkgPath),
		structNameGenerator(to, pkgPath),
	)
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
