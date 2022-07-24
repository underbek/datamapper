package generator

import "github.com/underbek/datamapper/models"

type ConversionRule int

const (
	UndefinedRule ConversionRule = iota
	NeedOnlyAssigmentRule
	NeedCallConversionFunctionRule
	NeedCallConversionFunctionSeparatelyRule
	NeedCallConversionFunctionWithErrorRule
	PointerPoPointerConversionFunctionsRule
)

func getConversionRule(fromField, toField models.Field, cf models.ConversionFunction) ConversionRule {
	if isNeedOnlyAssigmentRule(fromField, toField, cf) {
		return NeedOnlyAssigmentRule
	}

	if isNeedCallConversionFunctionRule(fromField, toField, cf) {
		return NeedCallConversionFunctionRule
	}

	if isPointerPoPointerConversionFunctionsRule(fromField, toField, cf) {
		return PointerPoPointerConversionFunctionsRule
	}

	if isNeedCallConversionFunctionWithErrorRule(fromField, toField, cf) {
		return NeedCallConversionFunctionWithErrorRule
	}

	if isNeedCallConversionFunctionSeparatelyRule(fromField, toField, cf) {
		return NeedCallConversionFunctionSeparatelyRule
	}

	return UndefinedRule
}

func isNeedPointerCheckAndReturnError(fromField, toField models.Field, cf models.ConversionFunction) bool {
	if fromField.Type.Pointer && !cf.FromType.Pointer {
		return true
	}

	return false
}

func isNeedOnlyAssigmentRule(fromField, toField models.Field, cf models.ConversionFunction) bool {
	if fromField.Type == toField.Type {
		return true
	}

	if fromField.Type.Package.Path == toField.Type.Package.Path &&
		fromField.Type.Name == toField.Type.Name &&
		!fromField.Type.Pointer && toField.Type.Pointer {

		return true
	}

	return false
}

func isNeedCallConversionFunctionRule(fromField, toField models.Field, cf models.ConversionFunction) bool {
	if cf.WithError {
		return false
	}

	if toField.Type.Pointer == cf.ToType.Pointer {
		return true
	}

	//TODO: if !fromField.Type.Pointer && cf.FromType.Pointer
	return false
}

func isPointerPoPointerConversionFunctionsRule(fromField, toField models.Field, cf models.ConversionFunction) bool {
	if fromField.Type.Pointer && toField.Type.Pointer && !cf.FromType.Pointer && !cf.ToType.Pointer {
		return true
	}

	return false
}

func isNeedCallConversionFunctionWithErrorRule(fromField, toField models.Field, cf models.ConversionFunction) bool {
	return cf.WithError
}

func isNeedCallConversionFunctionSeparatelyRule(fromField, toField models.Field, cf models.ConversionFunction) bool {
	if cf.WithError {
		return false
	}

	if !cf.ToType.Pointer && toField.Type.Pointer {
		return true
	}

	return false
}
