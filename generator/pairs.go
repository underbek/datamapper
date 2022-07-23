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

	// TODO: check package
	if from.Type.Name == to.Type.Name {
		res, pkg, err := getFieldsPairBySameTypes(from, to, fromModel.Name, toModel.Name)
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

	//TODO: Use conversion functions with pointers
	key.FromType.Pointer = false
	key.ToType.Pointer = false

	cf, ok := functions[key]

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

func fillConversionFunction(pair FieldsPair, fromFiled, toFiled models.Field, fromModel, toModel models.Struct,
	cf models.ConversionFunction, pkgPath string) (FieldsPair, map[models.Package]struct{}, error) {
	pkgs := map[models.Package]struct{}{cf.Package: {}}

	ptr := ""
	if fromFiled.Type.Pointer {
		ptr = "*"
	}

	packageName := cf.Package.Name
	if cf.Package.Alias != "" {
		packageName = cf.Package.Alias
	}

	typeParams := getTypeParams(cf, fromFiled.Type, toFiled.Type)
	conversion := fmt.Sprintf("%s.%s%s(%sfrom.%s)", packageName, cf.Name, typeParams, ptr, fromFiled.Name)
	if cf.Package.Path == pkgPath {
		conversion = fmt.Sprintf("%s%s(%sfrom.%s)", cf.Name, typeParams, ptr, fromFiled.Name)
	}

	if fromFiled.Type.Pointer && !toFiled.Type.Pointer {
		pointerCheck, err := getPointerCheck(fromFiled, toFiled,
			getFullStructName(fromModel, pkgPath),
			getFullStructName(toModel, pkgPath),
		)
		if err != nil {
			return FieldsPair{}, nil, err
		}

		pair.Conversions = append(pair.Conversions, pointerCheck)
		pkgs[models.Package{
			Name: "fmt",
			Path: "fmt",
		}] = struct{}{}
		pair.PointerToValue = true
	}

	if !toFiled.Type.Pointer && !cf.WithError {
		pair.Assignment = conversion
		return pair, pkgs, nil
	}

	var err error

	refAssignment := fmt.Sprintf("&from%s", fromFiled.Name)
	valueAssignment := fmt.Sprintf("from%s", fromFiled.Name)

	if fromFiled.Type.Pointer && toFiled.Type.Pointer {
		pointerToPointer, err := getPointerToPointerConversion(
			fromFiled.Name,
			getFullStructName(toModel, pkgPath),
			getFullFieldName(toFiled, pkgPath),
			conversion,
			cf.WithError,
		)
		if err != nil {
			return FieldsPair{}, nil, err
		}

		pkgs[toFiled.Type.Package] = struct{}{}
		pair.Conversions = append(pair.Conversions, pointerToPointer)
		pair.Assignment = valueAssignment
		return pair, pkgs, nil
	}

	if toFiled.Type.Pointer && !cf.WithError {
		conversion, err = getPointerConversion(fromFiled.Name, conversion)
		if err != nil {
			return FieldsPair{}, nil, err
		}

		pair.Conversions = append(pair.Conversions, conversion)
	}

	if cf.WithError {
		conversion, err = getErrorConversion(fromFiled.Name, getFullStructName(toModel, pkgPath), conversion)
		if err != nil {
			return FieldsPair{}, nil, err
		}

		pair.Conversions = append(pair.Conversions, conversion)
	}

	if toFiled.Type.Pointer {
		pair.Assignment = refAssignment
		return pair, pkgs, nil
	}

	pair.Assignment = valueAssignment
	return pair, pkgs, nil
}
