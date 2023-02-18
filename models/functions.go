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
	Name      string        `yaml:"name"`
	Package   Package       `yaml:"package"`
	FromType  Type          `yaml:"from_type"`
	ToType    Type          `yaml:"to_type"`
	TypeParam TypeParamType `yaml:"type_param"`
	WithError bool          `yaml:"with_error"`
}

type Functions = map[ConversionFunctionKey]ConversionFunction
