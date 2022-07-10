package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/underbek/datamapper/models"
)

func Test_CFIncorrectFile(t *testing.T) {
	_, err := ParseConversionFunctions("incorrect name")
	require.Error(t, err)
}

func Test_CFParseEmptyFile(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
	}{
		{
			name:     "Empty file",
			fileName: "empty_functions.go",
		},
		{
			name:     "Empty structs",
			fileName: "incorrect_functions.go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ParseConversionFunctions(testPath + tt.fileName)
			assert.NoError(t, err)
			assert.Empty(t, res)
		})
	}
}

func Test_CFParseSimpleFunctions(t *testing.T) {
	res, err := ParseConversionFunctions(testPath + "simple_conversions.go")
	require.NoError(t, err)
	require.Len(t, res, 2)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertIntToString"},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "int"}, ToType: models.Type{Name: "string"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertFloatToString"},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "float32"}, ToType: models.Type{Name: "string"}}],
	)
}

func Test_CFParseGenericFrom(t *testing.T) {
	res, err := ParseConversionFunctions(testPath + "generic_from.go")
	require.NoError(t, err)
	assert.Len(t, res, 8)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertAnyToString"},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "any"}, ToType: models.Type{Name: "string"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertIntUintToString"},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "int"}, ToType: models.Type{Name: "string"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertIntUintToString"},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "uint"}, ToType: models.Type{Name: "string"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertIntegersToString"},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "int8"}, ToType: models.Type{Name: "string"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertIntegersToString"},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "int16"}, ToType: models.Type{Name: "string"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertIntegersToString"},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "int32"}, ToType: models.Type{Name: "string"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertXFloatToString"},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "float32"}, ToType: models.Type{Name: "string"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertXFloatToString"},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "float64"}, ToType: models.Type{Name: "string"}}],
	)
}

func Test_CFParseGenericTo(t *testing.T) {
	res, err := ParseConversionFunctions(testPath + "generic_to.go")
	require.NoError(t, err)
	assert.Len(t, res, 8)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertStringToAny"},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "string"}, ToType: models.Type{Name: "any"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertStringToIntUint"},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "string"}, ToType: models.Type{Name: "int"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertStringToIntUint"},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "string"}, ToType: models.Type{Name: "uint"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertStringToIntegers"},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "string"}, ToType: models.Type{Name: "int8"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertStringToIntegers"},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "string"}, ToType: models.Type{Name: "int16"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertStringToIntegers"},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "string"}, ToType: models.Type{Name: "int32"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertStringToXFloat"},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "string"}, ToType: models.Type{Name: "float32"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertStringToXFloat"},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "string"}, ToType: models.Type{Name: "float64"}}],
	)
}

func Test_CFParseGenericFromTo(t *testing.T) {
	res, err := ParseConversionFunctions(testPath + "generic_from_to.go")
	require.NoError(t, err)
	assert.Len(t, res, 4)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertXFloatToIntegers"},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "float32"}, ToType: models.Type{Name: "int"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertXFloatToIntegers"},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "float32"}, ToType: models.Type{Name: "uint"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertXFloatToIntegers"},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "float64"}, ToType: models.Type{Name: "int"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertXFloatToIntegers"},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "float64"}, ToType: models.Type{Name: "uint"}}],
	)
}

func Test_CFParseGenericStruct(t *testing.T) {
	res, err := ParseConversionFunctions(testPath + "generic_struct.go")
	require.NoError(t, err)
	assert.Len(t, res, 2)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertModelsToString"},
		res[models.ConversionFunctionKey{FromType: models.Type{
			Name:    "Model",
			Package: "parser",
		}, ToType: models.Type{Name: "string"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertModelsToString"},
		res[models.ConversionFunctionKey{FromType: models.Type{
			Name:    "Model",
			Package: "other",
		}, ToType: models.Type{Name: "string"}}],
	)
}

func Test_CFParseWithStruct(t *testing.T) {
	res, err := ParseConversionFunctions(testPath + "with_struct.go")
	require.NoError(t, err)
	assert.Len(t, res, 2)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertCurrentModelToOther"},
		res[models.ConversionFunctionKey{FromType: models.Type{
			Name:    "Model",
			Package: "parser",
		}, ToType: models.Type{
			Name:    "Model",
			Package: "other",
		}}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertOtherModelToCurrent"},
		res[models.ConversionFunctionKey{FromType: models.Type{
			Name:    "Model",
			Package: "other",
		}, ToType: models.Type{
			Name:    "Model",
			Package: "parser",
		}}],
	)
}
