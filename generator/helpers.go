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
		name := model.Type.Name

		if model.Type.Package.Path == pkgPath {
			return name
		}

		pkgName := model.Type.Package.Name
		if model.Type.Package.Alias != "" {
			pkgName = model.Type.Package.Alias
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

func isSameTypesWithoutPointer(from, to models.Type) bool {
	from.Pointer = false
	to.Pointer = false

	return from == to
}

func getConversionFunction(fromType, toType models.Type, fromName string, functions models.Functions,
) (models.ConversionFunction, error) {

	if isSameTypesWithoutPointer(fromType, toType) {
		return models.ConversionFunction{}, nil
	}

	key := models.ConversionFunctionKey{
		FromType: fromType,
		ToType:   toType,
	}

	cf, ok := functions[key]
	if ok {
		return cf, nil
	}

	key.FromType.Pointer = false
	key.ToType.Pointer = false

	cf, ok = functions[key]
	if ok {
		return cf, nil
	}

	if !ok && fromType.Kind == models.SliceType && toType.Kind == models.SliceType {
		return getConversionFunction(
			fromType.Additional.(models.SliceAdditional).InType,
			toType.Additional.(models.SliceAdditional).InType,
			fromName,
			functions,
		)
	}

	return models.ConversionFunction{}, fmt.Errorf(
		"not found convertor function for types %s -> %s by %s field: %w",
		key.FromType.Name,
		key.ToType.Name,
		fromName,
		ErrNotFound,
	)
}

func isPointerToValue(fromType, fromCfType models.Type) bool {
	if fromType.Pointer && !fromCfType.Pointer {
		return true
	}

	return false
}

func getPointerSymbol(fromFieldType, cfFromType models.Type) string {
	if fromFieldType.Pointer && !cfFromType.Pointer {
		return "*"
	}

	return ""
}

func getConversionFunctionCall(cf models.ConversionFunction, fromFieldType, toFieldType models.Type, pkgPath,
	arg string) string {

	packageName := cf.Package.Name
	if cf.Package.Alias != "" {
		packageName = cf.Package.Alias
	}

	ptr := getPointerSymbol(fromFieldType, cf.FromType)
	typeParams := getTypeParams(cf, fromFieldType, toFieldType)

	if cf.Package.Path == pkgPath {
		return fmt.Sprintf("%s%s(%s%s)", cf.Name, typeParams, ptr, arg)
	}

	return fmt.Sprintf("%s.%s%s(%s%s)", packageName, cf.Name, typeParams, ptr, arg)
}
