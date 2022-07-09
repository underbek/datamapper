package converts

import "fmt"

type ConvertorType = string
type ImportType = string

type Factory struct{}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) GetConvertorFunctions(fromType, toType, fromFieldName string) (ConvertorType, ImportType) {
	if fromType == toType {
		return fmt.Sprintf("from.%s", fromFieldName), ""
	}

	if toType == stringType && (isNumeric(fromType) || isBoolean(toType)) {
		return fmt.Sprintf("fmt.Sprint(from.%s)", fromFieldName), "fmt"
	}

	return "", ""
}
