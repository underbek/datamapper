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
		res[models.ConversionFunctionKey{FromType: "int", ToType: "string"}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertFloatToString"},
		res[models.ConversionFunctionKey{FromType: "float32", ToType: "string"}],
	)
}

func Test_CFParseGenericFrom(t *testing.T) {
	res, err := ParseConversionFunctions(testPath + "generic_from.go")
	require.NoError(t, err)
	assert.Len(t, res, 8)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertAnyToString"},
		res[models.ConversionFunctionKey{FromType: "any", ToType: "string"}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertIntUintToString"},
		res[models.ConversionFunctionKey{FromType: "int", ToType: "string"}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertIntUintToString"},
		res[models.ConversionFunctionKey{FromType: "uint", ToType: "string"}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertIntegersToString"},
		res[models.ConversionFunctionKey{FromType: "int8", ToType: "string"}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertIntegersToString"},
		res[models.ConversionFunctionKey{FromType: "int16", ToType: "string"}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertIntegersToString"},
		res[models.ConversionFunctionKey{FromType: "int32", ToType: "string"}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertXFloatToString"},
		res[models.ConversionFunctionKey{FromType: "float32", ToType: "string"}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertXFloatToString"},
		res[models.ConversionFunctionKey{FromType: "float64", ToType: "string"}],
	)
}

func Test_CFParseGenericTo(t *testing.T) {
	res, err := ParseConversionFunctions(testPath + "generic_to.go")
	require.NoError(t, err)
	assert.Len(t, res, 8)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertStringToAny"},
		res[models.ConversionFunctionKey{FromType: "string", ToType: "any"}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertStringToIntUint"},
		res[models.ConversionFunctionKey{FromType: "string", ToType: "int"}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertStringToIntUint"},
		res[models.ConversionFunctionKey{FromType: "string", ToType: "uint"}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertStringToIntegers"},
		res[models.ConversionFunctionKey{FromType: "string", ToType: "int8"}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertStringToIntegers"},
		res[models.ConversionFunctionKey{FromType: "string", ToType: "int16"}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertStringToIntegers"},
		res[models.ConversionFunctionKey{FromType: "string", ToType: "int32"}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertStringToXFloat"},
		res[models.ConversionFunctionKey{FromType: "string", ToType: "float32"}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertStringToXFloat"},
		res[models.ConversionFunctionKey{FromType: "string", ToType: "float64"}],
	)
}

func Test_CFParseGenericFromTo(t *testing.T) {
	res, err := ParseConversionFunctions(testPath + "generic_from_to.go")
	require.NoError(t, err)
	assert.Len(t, res, 4)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertXFloatToIntegers"},
		res[models.ConversionFunctionKey{FromType: "float32", ToType: "int"}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertXFloatToIntegers"},
		res[models.ConversionFunctionKey{FromType: "float32", ToType: "uint"}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertXFloatToIntegers"},
		res[models.ConversionFunctionKey{FromType: "float64", ToType: "int"}],
	)

	assert.Equal(t,
		models.ConversionFunction{Name: "ConvertXFloatToIntegers"},
		res[models.ConversionFunctionKey{FromType: "float64", ToType: "uint"}],
	)
}
