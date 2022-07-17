package generator

import (
	"fmt"

	"github.com/underbek/datamapper/models"
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

func filterImports(currentPkgPath string, imports []string) []string {
	res := make([]string, 0, len(imports))
	for _, imp := range imports {
		if imp != currentPkgPath {
			res = append(res, imp)
		}
	}

	return res
}

func fillConversations(fields []FieldsPair) []string {
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
