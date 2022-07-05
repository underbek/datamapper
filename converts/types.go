package converts

type basicType = string

const (
	intType   basicType = "int"
	int8Type  basicType = "int8"
	int16Type basicType = "int16"
	int32Type basicType = "int32"
	int64Type basicType = "int64"

	uintType   basicType = "uint"
	uint8Type  basicType = "uint8"
	uint16Type basicType = "uint16"
	uint32Type basicType = "uint32"
	uint64Type basicType = "uint64"

	uintptrType basicType = "uintptr"

	float32Type basicType = "float32"
	float64Type basicType = "float64"

	complex64Type  basicType = "complex64"
	complex128Type basicType = "complex128"

	byteType    basicType = "byte"
	runeType    basicType = "rune"
	stringType  basicType = "string"
	booleanType basicType = "bool"
)

func isSigned(fieldType string) bool {
	switch fieldType {
	case intType, int8Type, int16Type, int32Type, int64Type:
		return true
	}
	return false
}

func isUnsigned(fieldType string) bool {
	switch fieldType {
	case uintType, uint8Type, uint16Type, uint32Type, uint64Type, uintptrType:
		return true
	}
	return false
}

func isFloat(fieldType string) bool {
	switch fieldType {
	case float32Type, float64Type:
		return true
	}
	return false
}

func isComplex(fieldType string) bool {
	switch fieldType {
	case complex64Type, complex128Type:
		return true
	}
	return false
}

func isNumeric(fieldType string) bool {
	return isSigned(fieldType) || isUnsigned(fieldType) || isFloat(fieldType) || isComplex(fieldType)
}

func isBoolean(fieldType string) bool {
	return fieldType == booleanType
}
