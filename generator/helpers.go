package generator

import (
	"fmt"
	"strings"

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
	uniqConversions := make(map[string]struct{})
	var res []string
	for _, field := range fields {
		for _, conv := range field.Conversions {
			if _, ok := uniqConversions[conv]; !ok {
				uniqConversions[conv] = struct{}{}
				res = append(res, conv)
			}
		}
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

func getSkippedFieldsPointerCheckError(field models.Field, toTypeFullName, fromTypeName string) ([]string, error) {
	var res []string

	head := field.Head
	for head != nil {
		if head.Type.Pointer {
			conversion, err := getPointerCheck(
				createFieldPathWithPrefix(*head),
				toTypeFullName,
				fmt.Sprintf("errors.New(\"%s.%s is nil\")", fromTypeName, createFieldPath(*head)),
				true,
			)
			if err != nil {
				return nil, err
			}

			res = append(res, conversion)
		}

		head = head.Head
	}

	return res, nil
}

func findHead(fields []FieldsPair) TypeWithName {
	return fields[0].Types[0]
}

func createModelWithPairs(fields []FieldsPair, modelType TypeWithName) ModelWithPairs {
	return createModelWithPairsByLvl(fields, modelType, 0)
}

func createModelWithPairsByLvl(fields []FieldsPair, modelType TypeWithName, lvl int) ModelWithPairs {
	res := ModelWithPairs{
		Type: modelType,
	}

	var otherFields []FieldsPair
	var internalTypes []TypeWithName
	internalTypeMap := make(map[TypeWithName]struct{})

	for _, field := range fields {
		if field.Types[lvl] != modelType {
			continue
		}

		if len(field.Types) == lvl+1 {
			res.fields = append(res.fields, field)
			continue
		}

		if len(field.Types) == lvl+2 {
			currentType := field.Types[lvl+1]
			if _, ok := internalTypeMap[currentType]; !ok {
				internalTypes = append(internalTypes, currentType)
				internalTypeMap[currentType] = struct{}{}
			}
		}

		otherFields = append(otherFields, field)
	}

	if internalTypes == nil {
		return res
	}

	for _, internalType := range internalTypes {
		res.models = append(res.models, createModelWithPairsByLvl(otherFields, internalType, lvl+1))
	}

	return res
}

func makeFieldsPairByModel(pkg string, model ModelWithPairs) (FieldsPair, error) {
	for _, m := range model.models {
		field, err := makeFieldsPairByModel(pkg, m)
		if err != nil {
			return FieldsPair{}, err
		}

		model.fields = append(model.fields, field)
	}

	converter, err := fillResultStruct(
		strings.Replace(model.Type.Type.FullName(pkg), "*", "&", 1),
		model.fields,
	)
	if err != nil {
		return FieldsPair{}, err
	}

	return FieldsPair{
		Assignment: converter,
		ToName:     model.Type.FieldName,
	}, nil
}

func createResultConverter(pkg string, toName string, model ModelWithPairs) (string, error) {
	for _, m := range model.models {
		field, err := makeFieldsPairByModel(pkg, m)
		if err != nil {
			return "", err
		}

		model.fields = append(model.fields, field)
	}

	res, err := fillResultStruct(strings.Replace(toName, "*", "&", 1), model.fields)
	if err != nil {
		return "", err
	}

	return res, nil
}
