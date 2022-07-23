package models

type TypeParamType int

const (
	NoTypeParam = TypeParamType(iota)
	FromTypeParam
	ToTypeParam
	FromToTypeParam
)

type ConversionFunctionKey struct {
	FromType, ToType Type
}

type ConversionFunction struct {
	Name      string
	Package   Package
	TypeParam TypeParamType
	WithError bool
}

type Functions = map[ConversionFunctionKey]ConversionFunction
