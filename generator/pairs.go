package generator

import (
	"fmt"

	"github.com/underbek/datamapper/models"
	"golang.org/x/exp/maps"
)

func createModelsPair(from, to models.Struct, pkgPath string, functions models.Functions) (result, error) {
	var fields []FieldsPair
	imports := make(map[string]struct{})

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

		for _, pack := range packs {
			if pack != "" {
				imports[pack] = struct{}{}
			}
		}

		fields = append(fields, pair)
	}

	return result{
		fields:        fields,
		imports:       maps.Keys(imports),
		conversations: fillConversations(fields),
	}, nil
}

func getFieldsPair(from, to models.Field, fromModel, toModel models.Struct, pkgPath string, functions models.Functions,
) (FieldsPair, []ImportType, error) {

	// TODO: check package
	if from.Type.Name == to.Type.Name {
		res, pkg, err := getFieldsPairBySameTypes(from, to, fromModel.Name, toModel.Name)
		if err != nil {
			return FieldsPair{}, nil, err
		}
		return res, []ImportType{pkg}, nil
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
			"not found convertor function for types %s -> %s by %s field",
			from.Type.Name,
			to.Type.Name,
			from.Name,
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

func getFieldsPairBySameTypes(from, to models.Field, fromName, toName string) (FieldsPair, ImportType, error) {
	res := FieldsPair{
		FromName: from.Name,
		FromType: from.Type.Name,
		ToName:   to.Name,
		ToType:   to.Type.Name,
	}

	if from.Type.Pointer == to.Type.Pointer {
		res.Assignment = fmt.Sprintf("from.%s", from.Name)
		return res, "", nil
	}

	if to.Type.Pointer {
		res.Assignment = fmt.Sprintf("&from.%s", from.Name)
		return res, "", nil
	}

	res.PointerToValue = true
	res.Assignment = fmt.Sprintf("*from.%s", from.Name)

	conversion, err := getPointerCheck(from, to, fromName, toName)
	if err != nil {
		return FieldsPair{}, "", err
	}

	res.Conversions = append(res.Conversions, conversion)

	return res, "fmt", nil
}

func fillConversionFunction(pair FieldsPair, fromFiled, toFiled models.Field, fromModel, toModel models.Struct,
	cf models.ConversionFunction, pkgPath string) (FieldsPair, []ImportType, error) {
	imports := []ImportType{cf.PackagePath}

	ptr := ""
	if fromFiled.Type.Pointer {
		pointerCheck, err := getPointerCheck(fromFiled, toFiled,
			getFullStructName(fromModel, pkgPath),
			getFullStructName(toModel, pkgPath),
		)
		if err != nil {
			return FieldsPair{}, nil, err
		}

		pair.Conversions = append(pair.Conversions, pointerCheck)
		ptr = "*"
		imports = append(imports, "fmt")
		pair.PointerToValue = true
	}

	typeParams := getTypeParams(cf, fromFiled.Type, toFiled.Type)
	conversion := fmt.Sprintf("%s.%s%s(%sfrom.%s)", cf.PackageName, cf.Name, typeParams, ptr, fromFiled.Name)
	if cf.PackagePath == pkgPath {
		conversion = fmt.Sprintf("%s%s(%sfrom.%s)", cf.Name, typeParams, ptr, fromFiled.Name)
	}

	if !toFiled.Type.Pointer && !cf.WithError {
		pair.Assignment = conversion
		return pair, imports, nil
	}

	var err error

	refAssignment := fmt.Sprintf("&from%s", fromFiled.Name)
	valueAssignment := fmt.Sprintf("from%s", fromFiled.Name)

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
		return pair, imports, nil
	}

	pair.Assignment = valueAssignment
	return pair, imports, nil
}
