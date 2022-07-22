package parser

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/underbek/datamapper/models"
)

type Type struct {
	models.Type
	generic bool
}

func parseType(t types.Type) ([]Type, error) {
	switch t := t.(type) {
	case *types.Named:
		und := t.Underlying()
		if _, ok := und.(*types.Struct); ok {
			return []Type{{Type: models.Type{
				Name:        t.Obj().Name(),
				PackagePath: t.Obj().Pkg().Path(),
			}}}, nil
		}
		return parseType(t.Underlying())
	case *types.Interface:
		n := t.NumEmbeddeds()
		if n == 0 {
			return []Type{{Type: models.Type{Name: t.String()}}}, nil
		}

		res := make([]Type, 0, n)
		for i := 0; i < n; i++ {
			names, err := parseType(t.EmbeddedType(i))
			if err != nil {
				return nil, err
			}

			res = append(res, names...)
		}
		return res, nil

	case *types.Union:
		n := t.Len()
		res := make([]Type, 0, n)
		for i := 0; i < n; i++ {
			names, err := parseType(t.Term(i).Type())
			if err != nil {
				return nil, err
			}

			res = append(res, names...)
		}
		return res, nil
	case *types.Basic, *types.Struct:
		return []Type{{Type: models.Type{Name: t.String()}}}, nil
	case *types.TypeParam:
		res, err := parseType(t.Underlying())
		if err != nil {
			return nil, err
		}

		for i := range res {
			res[i].generic = true
		}
		return res, nil
	case *types.Array, *types.Slice:
		return []Type{{Type: models.Type{Name: t.String()}}}, nil
	case *types.Pointer:
		res, err := parseType(t.Elem())
		if err != nil {
			return nil, err
		}

		for i := range res {
			res[i].Pointer = true
		}

		return res, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUndefinedType, t.String())
	}
}

func getTypeParam(fromGeneric, toGeneric bool) models.TypeParamType {
	if fromGeneric && toGeneric {
		return models.FromToTypeParam
	}

	if fromGeneric {
		return models.FromTypeParam
	}

	if toGeneric {
		return models.ToTypeParam
	}

	return models.NoTypeParam
}

func parseTag(tag string) []models.Tag {
	if tag == "" {
		return nil
	}

	value := strings.Trim(tag, "`")
	textTags := strings.Split(value, " ")

	tags := make([]models.Tag, 0, len(textTags))
	for _, textTag := range textTags {
		sepIndex := strings.Index(textTag, ":")
		if sepIndex == -1 {
			continue
		}

		valueTag := strings.Trim(textTag[sepIndex+1:], "\"")

		tags = append(tags, models.Tag{
			Name:  textTag[:sepIndex],
			Value: strings.Split(valueTag, ",")[0],
		})
	}

	return tags
}

func isErrorType(t types.Type) (bool, error) {
	errTypes, err := parseType(t)
	if err != nil {
		return false, err
	}

	if len(errTypes) != 1 {
		return false, nil
	}

	if errTypes[0].generic {
		return false, nil
	}

	// TODO: parse custom error by errors interface
	if errTypes[0].PackagePath != "" {
		return false, nil
	}

	if errTypes[0].Name != "interface{Error() string}" {
		return false, nil
	}

	return true, nil
}
