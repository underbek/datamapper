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
		models.ConversionFunction{
			Name:        "ConvertIntToString",
			PackageName: "parser",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
		},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "int"}, ToType: models.Type{Name: "string"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{
			Name:        "ConvertFloatToString",
			PackageName: "parser",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
		},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "float32"}, ToType: models.Type{Name: "string"}}],
	)
}

func Test_CFParseGenericFrom(t *testing.T) {
	res, err := ParseConversionFunctions(testPath + "generic_from.go")
	require.NoError(t, err)
	assert.Len(t, res, 8)

	assert.Equal(t,
		models.ConversionFunction{
			Name:        "ConvertAnyToString",
			PackageName: "parser",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
			TypeParam:   models.FromTypeParam,
		},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "any"}, ToType: models.Type{Name: "string"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{
			Name:        "ConvertIntUintToString",
			PackageName: "parser",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
			TypeParam:   models.FromTypeParam,
		},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "int"}, ToType: models.Type{Name: "string"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{
			Name:        "ConvertIntUintToString",
			PackageName: "parser",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
			TypeParam:   models.FromTypeParam,
		},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "uint"}, ToType: models.Type{Name: "string"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{
			Name:        "ConvertIntegersToString",
			PackageName: "parser",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
			TypeParam:   models.FromTypeParam,
		},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "int8"}, ToType: models.Type{Name: "string"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{
			Name:        "ConvertIntegersToString",
			PackageName: "parser",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
			TypeParam:   models.FromTypeParam,
		},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "int16"}, ToType: models.Type{Name: "string"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{
			Name:        "ConvertIntegersToString",
			PackageName: "parser",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
			TypeParam:   models.FromTypeParam,
		},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "int32"}, ToType: models.Type{Name: "string"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{
			Name:        "ConvertXFloatToString",
			PackageName: "parser",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
			TypeParam:   models.FromTypeParam,
		},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "float32"}, ToType: models.Type{Name: "string"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{
			Name:        "ConvertXFloatToString",
			PackageName: "parser",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
			TypeParam:   models.FromTypeParam,
		},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "float64"}, ToType: models.Type{Name: "string"}}],
	)
}

func Test_CFParseGenericTo(t *testing.T) {
	res, err := ParseConversionFunctions(testPath + "generic_to.go")
	require.NoError(t, err)
	assert.Len(t, res, 8)

	assert.Equal(t,
		models.ConversionFunction{
			Name:        "ConvertStringToAny",
			PackageName: "parser",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
			TypeParam:   models.ToTypeParam,
		},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "string"}, ToType: models.Type{Name: "any"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{
			Name:        "ConvertStringToIntUint",
			PackageName: "parser",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
			TypeParam:   models.ToTypeParam,
		},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "string"}, ToType: models.Type{Name: "int"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{
			Name:        "ConvertStringToIntUint",
			PackageName: "parser",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
			TypeParam:   models.ToTypeParam,
		},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "string"}, ToType: models.Type{Name: "uint"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{
			Name:        "ConvertStringToIntegers",
			PackageName: "parser",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
			TypeParam:   models.ToTypeParam,
		},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "string"}, ToType: models.Type{Name: "int8"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{
			Name:        "ConvertStringToIntegers",
			PackageName: "parser",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
			TypeParam:   models.ToTypeParam,
		},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "string"}, ToType: models.Type{Name: "int16"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{
			Name:        "ConvertStringToIntegers",
			PackageName: "parser",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
			TypeParam:   models.ToTypeParam,
		},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "string"}, ToType: models.Type{Name: "int32"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{
			Name:        "ConvertStringToXFloat",
			PackageName: "parser",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
			TypeParam:   models.ToTypeParam,
		},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "string"}, ToType: models.Type{Name: "float32"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{
			Name:        "ConvertStringToXFloat",
			PackageName: "parser",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
			TypeParam:   models.ToTypeParam,
		},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "string"}, ToType: models.Type{Name: "float64"}}],
	)
}

func Test_CFParseGenericFromTo(t *testing.T) {
	res, err := ParseConversionFunctions(testPath + "generic_from_to.go")
	require.NoError(t, err)
	assert.Len(t, res, 4)

	assert.Equal(t,
		models.ConversionFunction{
			Name:        "ConvertXFloatToIntegers",
			PackageName: "parser",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
			TypeParam:   models.FromToTypeParam,
		},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "float32"}, ToType: models.Type{Name: "int"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{
			Name:        "ConvertXFloatToIntegers",
			PackageName: "parser",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
			TypeParam:   models.FromToTypeParam,
		},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "float32"}, ToType: models.Type{Name: "uint"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{
			Name:        "ConvertXFloatToIntegers",
			PackageName: "parser",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
			TypeParam:   models.FromToTypeParam,
		},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "float64"}, ToType: models.Type{Name: "int"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{
			Name:        "ConvertXFloatToIntegers",
			PackageName: "parser",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
			TypeParam:   models.FromToTypeParam,
		},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "float64"}, ToType: models.Type{Name: "uint"}}],
	)
}

func Test_CFParseGenericStruct(t *testing.T) {
	res, err := ParseConversionFunctions(testPath + "generic_struct.go")
	require.NoError(t, err)
	assert.Len(t, res, 2)

	assert.Equal(t,
		models.ConversionFunction{
			Name:        "ConvertModelsToString",
			PackageName: "parser",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
			TypeParam:   models.FromTypeParam,
		},
		res[models.ConversionFunctionKey{FromType: models.Type{
			Name:        "Model",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
		}, ToType: models.Type{Name: "string"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{
			Name:        "ConvertModelsToString",
			PackageName: "parser",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
			TypeParam:   models.FromTypeParam,
		},
		res[models.ConversionFunctionKey{FromType: models.Type{
			Name:        "Model",
			PackagePath: "github.com/underbek/datamapper/_test_data/other",
		}, ToType: models.Type{Name: "string"}}],
	)
}

func Test_CFParseWithStruct(t *testing.T) {
	res, err := ParseConversionFunctions(testPath + "with_struct.go")
	require.NoError(t, err)
	assert.Len(t, res, 2)

	assert.Equal(t,
		models.ConversionFunction{
			Name:        "ConvertCurrentModelToOther",
			PackageName: "parser",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
		},
		res[models.ConversionFunctionKey{FromType: models.Type{
			Name:        "Model",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
		}, ToType: models.Type{
			Name:        "Model",
			PackagePath: "github.com/underbek/datamapper/_test_data/other",
		}}],
	)

	assert.Equal(t,
		models.ConversionFunction{
			Name:        "ConvertOtherModelToCurrent",
			PackageName: "parser",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
		},
		res[models.ConversionFunctionKey{FromType: models.Type{
			Name:        "Model",
			PackagePath: "github.com/underbek/datamapper/_test_data/other",
		}, ToType: models.Type{
			Name:        "Model",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
		}}],
	)
}

func Test_CFParseWithError(t *testing.T) {
	res, err := ParseConversionFunctions(testPath + "with_error.go")
	require.NoError(t, err)
	assert.Len(t, res, 6)

	for _, toTypeName := range []string{"int", "int8", "int16", "int32", "int64"} {
		assert.Equal(t,
			models.ConversionFunction{
				Name:        "ConvertStringToSigned",
				PackageName: "parser",
				PackagePath: "github.com/underbek/datamapper/_test_data/parser",
				TypeParam:   models.ToTypeParam,
				WithError:   true,
			},
			res[models.ConversionFunctionKey{FromType: models.Type{
				Name: "string",
			}, ToType: models.Type{
				Name: toTypeName,
			}}],
		)
	}

	assert.Equal(t,
		models.ConversionFunction{
			Name:        "ConvertStringToDecimal",
			PackageName: "parser",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
			WithError:   true,
		},
		res[models.ConversionFunctionKey{FromType: models.Type{
			Name: "string",
		}, ToType: models.Type{
			Name:        "Decimal",
			PackagePath: "github.com/shopspring/decimal",
		}}],
	)
}