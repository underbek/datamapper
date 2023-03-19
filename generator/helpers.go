package generator

import (
	"fmt"

	"github.com/underbek/datamapper/models"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

func filterAndSortImports(currentPkgPath string, imports []ImportType) []ImportType {
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

func generateConvertorName(from, to models.Type, pkgPath string, kind models.KindOfType) string {
	structNameGenerator := func(t models.Type, pkgPath string) string {
		name := t.Name

		if t.Package.Path == pkgPath {
			return name
		}

		pkgName := t.Package.Name
		if t.Package.Alias != "" {
			pkgName = t.Package.Alias
		}
		return cases.Title(language.Und, cases.NoLower).String(pkgName) + name
	}

	prefix := ""
	switch kind {
	case models.SliceType:
		prefix = "Slice"
	}

	return fmt.Sprintf(
		"Convert%s%sTo%s%s",
		structNameGenerator(from, pkgPath),
		prefix,
		structNameGenerator(to, pkgPath),
		prefix,
	)
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

	return models.ConversionFunction{}, NewFindFieldsPairError(fromType, toType, fromName)
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

func getFieldPointerCheckError(fromModelName, toModelName, fromFieldName, toFieldName string) string {
	return fmt.Sprintf(`errors.New("cannot convert %s.%s -> %s.%s, field is nil")`,
		fromModelName,
		fromFieldName,
		toModelName,
		toFieldName,
	)
}
