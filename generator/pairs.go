package generator

import (
	"fmt"

	"github.com/underbek/datamapper/models"
	"golang.org/x/exp/maps"
)

func createModelsPair(from, to models.Struct, pkgPath string, functions models.Functions) (result, error) {
	var fields []FieldsPair
	packages := make(map[models.Package]struct{})

	fromFields := make(map[string]models.Field)
	for _, field := range from.Fields {
		fromFields[field.Tags[0].Value] = field
	}

	for _, toField := range to.Fields {
		fromField, ok := fromFields[toField.Tags[0].Value]
		if !ok {
			//TODO: warning or error politics
			continue
		}

		pair, packs, err := getFieldsPair(fromField, toField, from, to, pkgPath, functions)
		if err != nil {
			return result{}, err
		}

		maps.Copy(packages, packs)
		fields = append(fields, pair)
	}

	return result{
		fields:      fields,
		packages:    packages,
		conversions: fillConversions(fields),
	}, nil
}

func getFieldsPair(from, to models.Field, fromModel, toModel models.Struct, pkgPath string, functions models.Functions,
) (FieldsPair, map[models.Package]struct{}, error) {

	if from.Type.Name == to.Type.Name && from.Type.Package.Path == to.Type.Package.Path {
		res, pkg, err := getFieldsPairBySameTypes(from, to, fromModel.Type.Name, toModel.Type.Name)
		if err != nil {
			return FieldsPair{}, nil, err
		}

		if pkg != nil {
			return res, map[models.Package]struct{}{*pkg: {}}, nil
		}

		return res, nil, nil
	}

	key := models.ConversionFunctionKey{
		FromType: from.Type,
		ToType:   to.Type,
	}

	cf, ok := functions[key]

	if !ok {
		key.FromType.Pointer = false
		key.ToType.Pointer = false
		cf, ok = functions[key]
	}

	if !ok {
		return FieldsPair{}, nil, fmt.Errorf(
			"not found convertor function for types %s -> %s by %s field: %w",
			from.Type.Name,
			to.Type.Name,
			from.Name,
			ErrNotFound,
		)
	}

	res := FieldsPair{
		FromName:  from.Name,
		FromType:  from.Type.Name,
		ToName:    to.Name,
		ToType:    to.Type.Name,
		WithError: cf.WithError,
	}

	return fillConversionFunction(res, from, to, fromModel, toModel, cf, pkgPath)
}

func getFieldsPairBySameTypes(from, to models.Field, fromName, toName string) (FieldsPair, *models.Package, error) {
	res := FieldsPair{
		FromName: from.Name,
		FromType: from.Type.Name,
		ToName:   to.Name,
		ToType:   to.Type.Name,
	}

	if from.Type.Pointer == to.Type.Pointer {
		res.Assignment = fmt.Sprintf("from.%s", from.Name)
		return res, nil, nil
	}

	if to.Type.Pointer {
		res.Assignment = fmt.Sprintf("&from.%s", from.Name)
		return res, nil, nil
	}

	res.PointerToValue = true
	res.Assignment = fmt.Sprintf("*from.%s", from.Name)

	conversion, err := getPointerCheck(from, to, fromName, toName)
	if err != nil {
		return FieldsPair{}, nil, err
	}

	res.Conversions = append(res.Conversions, conversion)

	return res, &models.Package{
		Name: "fmt",
		Path: "fmt",
	}, nil
}

func fillConversionFunction(pair FieldsPair, fromField, toField models.Field, fromModel, toModel models.Struct,
	cf models.ConversionFunction, pkgPath string) (FieldsPair, map[models.Package]struct{}, error) {
	pkgs := map[models.Package]struct{}{cf.Package: {}}

	packageName := cf.Package.Name
	if cf.Package.Alias != "" {
		packageName = cf.Package.Alias
	}

	ptr := ""
	if fromField.Type.Pointer && !cf.FromType.Pointer {
		ptr = "*"
		pair.PointerToValue = true
	}

	typeParams := getTypeParams(cf, fromField.Type, toField.Type)

	cfCall := fmt.Sprintf("%s.%s%s(%sfrom.%s)", packageName, cf.Name, typeParams, ptr, fromField.Name)
	if cf.Package.Path == pkgPath {
		cfCall = fmt.Sprintf("%s%s(%sfrom.%s)", cf.Name, typeParams, ptr, fromField.Name)
	}

	refAssignment := fmt.Sprintf("&from%s", fromField.Name)
	valueAssignment := fmt.Sprintf("from%s", fromField.Name)

	if isNeedPointerCheckAndReturnError(fromField, toField, cf) {
		conversion, err := getPointerCheck(fromField, toField,
			fromModel.Type.FullName(pkgPath),
			toModel.Type.FullName(pkgPath),
		)
		if err != nil {
			return FieldsPair{}, nil, err
		}

		pkgs[models.Package{
			Name: "fmt",
			Path: "fmt",
		}] = struct{}{}

		pair.PointerToValue = true
		pair.Conversions = []string{conversion}
	}

	switch getConversionRule(fromField, toField, cf) {
	case NeedCallConversionFunctionRule:
		pair.Assignment = cfCall
		return pair, pkgs, nil

	case NeedCallConversionFunctionSeparatelyRule:
		conversion, err := getPointerConversion(fromField.Name, cfCall)
		if err != nil {
			return FieldsPair{}, nil, err
		}
		pair.Conversions = append(pair.Conversions, conversion)
		pair.Assignment = refAssignment
		return pair, pkgs, nil

	case PointerPoPointerConversionFunctionsRule:
		conversion, err := getPointerToPointerConversion(
			fromField.Name,
			toModel.Type.FullName(pkgPath),
			toField.Type.FullName(pkgPath),
			cfCall,
			cf.WithError,
		)
		if err != nil {
			return FieldsPair{}, nil, err
		}

		pkgs[toField.Type.Package] = struct{}{}

		// not use pointer check
		pair.Conversions = []string{conversion}
		pair.PointerToValue = false
		pair.Assignment = valueAssignment
		return pair, pkgs, nil

	case NeedCallConversionFunctionWithErrorRule:
		conversion, err := getErrorConversion(fromField.Name, toModel.Type.FullName(pkgPath), cfCall)
		if err != nil {
			return FieldsPair{}, nil, err
		}

		pair.Conversions = append(pair.Conversions, conversion)
		pair.Assignment = valueAssignment
		if toField.Type.Pointer && !cf.ToType.Pointer {
			pair.Assignment = refAssignment
		}
		return pair, pkgs, nil

	default:
		return FieldsPair{}, nil, fmt.Errorf(
			"%w: from field %s to field %s",
			ErrUndefinedConversionRule,
			fromField.Name,
			toField.Name,
		)
	}
}
