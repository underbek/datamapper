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
	if model.PackagePath != pkgPath {
		return fmt.Sprintf("%s.%s", model.PackageName, model.Name)
	}

	return model.Name
}

func getFullFieldName(filed models.Field, pkgPath string) string {
	if filed.Type.PackagePath != "" && filed.Type.PackagePath != pkgPath {
		pkgName := getPackageNameByPath(filed.Type.PackagePath)
		return fmt.Sprintf("%s.%s", pkgName, filed.Type.Name)
	}

	return filed.Type.Name
}

func getPackageNameByPath(path string) string {
	names := strings.Split(path, "/")
	return names[len(names)-1]
}

func filterAndSortImports(currentPkgPath string, imports []ImportType) []ImportType {
	set := make(map[ImportType]struct{})
	for _, imp := range imports {
		if imp != "" && imp != currentPkgPath {
			set[imp] = struct{}{}
		}
	}

	res := maps.Keys(set)
	slices.Sort(res)

	return res
}

func generateConvertorName(from, to models.Struct, pkgPath string) string {
	convertorName := "Convert"
	if from.PackagePath != pkgPath {
		convertorName += cases.Title(language.Und, cases.NoLower).String(from.PackageName)
	}
	convertorName += from.Name
	convertorName += "To"
	if to.PackagePath != pkgPath {
		convertorName += cases.Title(language.Und, cases.NoLower).String(to.PackageName)
	}
	convertorName += to.Name

	return convertorName
}

func generatePackageName(pkg *packages.Package) (string, error) {
	if pkg.Name != "" {
		return pkg.Name, nil
	}

	if pkg.PkgPath == "" {
		return "", fmt.Errorf("incorrect parsed destination package: %w", ErrParseError)
	}

	return getPackageNameByPath(pkg.PkgPath), nil
}
