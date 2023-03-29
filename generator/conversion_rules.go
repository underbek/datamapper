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
	NeedRangeBySlice
)

func getConversionRule(fromType, toType models.Type, cf models.ConversionFunction) ConversionRule {
	if isSameTypesWithoutPointer(fromType, toType) {
		return NeedOnlyAssigmentRule
	}

	if isNeedRangeBySlice(fromType, toType, cf) {
		return NeedRangeBySlice
	}

	if isNeedCallConversionFunctionRule(fromType, toType, cf) {
		return NeedCallConversionFunctionRule
	}

	if isPointerPoPointerConversionFunctionsRule(fromType, toType, cf) {
		return PointerPoPointerConversionFunctionsRule
	}

	if isNeedCallConversionFunctionWithErrorRule(fromType, toType, cf) {
		return NeedCallConversionFunctionWithErrorRule
	}

	if isNeedCallConversionFunctionSeparatelyRule(fromType, toType, cf) {
		return NeedCallConversionFunctionSeparatelyRule
	}

	return UndefinedRule
}

func isNeedPointerCheckAndReturnError(fromType, toType models.Type, cf models.ConversionFunction) bool {
	// if conversion by same types
	if fromType.Pointer == toType.Pointer {
		return false
	}

	defaultCf := models.ConversionFunction{}
	if cf == defaultCf {
		return fromType.Pointer && !toType.Pointer
	}

	if fromType.Pointer && !cf.FromType.Pointer {
		return true
	}

	return false
}

func isNeedCallConversionFunctionRule(fromType, toType models.Type, cf models.ConversionFunction) bool {
	if cf.WithError {
		return false
	}

	if toType.Pointer == cf.ToType.Pointer {
		return true
	}

	//TODO: if !fromType.Pointer && cf.FromType.Pointer
	return false
}

func isPointerPoPointerConversionFunctionsRule(fromType, toType models.Type, cf models.ConversionFunction) bool {
	if fromType.Pointer && toType.Pointer && !cf.FromType.Pointer && !cf.ToType.Pointer {
		return true
	}

	return false
}

func isNeedCallConversionFunctionWithErrorRule(fromType, toType models.Type, cf models.ConversionFunction) bool {
	return cf.WithError
}

func isNeedCallConversionFunctionSeparatelyRule(fromType, toType models.Type, cf models.ConversionFunction) bool {
	if cf.WithError {
		return false
	}

	if !cf.ToType.Pointer && toType.Pointer {
		return true
	}

	return false
}

func isNeedRangeBySlice(fromType, toType models.Type, cf models.ConversionFunction) bool {
	if fromType.Kind != models.SliceType {
		return false
	}

	if toType.Kind != models.SliceType {
		return false
	}

	if cf.FromType.Kind == models.SliceType {
		return false
	}

	if cf.ToType.Kind == models.SliceType {
		return false
	}

	return true
}

func isNeedPointerCheckSkippedFields(field models.Field) bool {
	head := field.Head
	for head != nil {
		if head.Type.Pointer {
			return true
		}

		head = head.Head
	}

	return false
}
