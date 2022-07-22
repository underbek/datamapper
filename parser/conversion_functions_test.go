package parser

import (
	"fmt"
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

	tests := []struct {
		Name          string
		FromTypeNames []string
	}{
		{
			Name:          "ConvertAnyToString",
			FromTypeNames: []string{"any"},
		},
		{
			Name:          "ConvertIntUintToString",
			FromTypeNames: []string{"int", "uint"},
		},
		{
			Name:          "ConvertIntegersToString",
			FromTypeNames: []string{"int8", "int16", "int32"},
		},
		{
			Name:          "ConvertXFloatToString",
			FromTypeNames: []string{"float32", "float64"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			for _, name := range tt.FromTypeNames {
				assert.Equal(t,
					models.ConversionFunction{
						Name:        tt.Name,
						PackageName: "parser",
						PackagePath: "github.com/underbek/datamapper/_test_data/parser",
						TypeParam:   models.FromTypeParam,
					},
					res[models.ConversionFunctionKey{
						FromType: models.Type{Name: name},
						ToType:   models.Type{Name: "string"},
					}],
				)
			}
		})
	}
}

func Test_CFParseGenericTo(t *testing.T) {
	res, err := ParseConversionFunctions(testPath + "generic_to.go")
	require.NoError(t, err)
	assert.Len(t, res, 8)

	tests := []struct {
		Name        string
		ToTypeNames []string
	}{
		{
			Name:        "ConvertStringToAny",
			ToTypeNames: []string{"any"},
		},
		{
			Name:        "ConvertStringToIntUint",
			ToTypeNames: []string{"int", "uint"},
		},
		{
			Name:        "ConvertStringToIntegers",
			ToTypeNames: []string{"int8", "int16", "int32"},
		},
		{
			Name:        "ConvertStringToXFloat",
			ToTypeNames: []string{"float32", "float64"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			for _, name := range tt.ToTypeNames {
				assert.Equal(t,
					models.ConversionFunction{
						Name:        tt.Name,
						PackageName: "parser",
						PackagePath: "github.com/underbek/datamapper/_test_data/parser",
						TypeParam:   models.ToTypeParam,
					},
					res[models.ConversionFunctionKey{
						FromType: models.Type{Name: "string"},
						ToType:   models.Type{Name: name},
					}],
				)
			}
		})
	}
}

func Test_CFParseGenericFromTo(t *testing.T) {
	res, err := ParseConversionFunctions(testPath + "generic_from_to.go")
	require.NoError(t, err)
	assert.Len(t, res, 4)

	tests := []struct {
		FromTypeName string
		ToTypeName   string
	}{
		{
			FromTypeName: "float32",
			ToTypeName:   "int",
		},
		{
			FromTypeName: "float32",
			ToTypeName:   "uint",
		},
		{
			FromTypeName: "float64",
			ToTypeName:   "int",
		},
		{
			FromTypeName: "float64",
			ToTypeName:   "uint",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s->%s", tt.FromTypeName, tt.ToTypeName), func(t *testing.T) {
			assert.Equal(t,
				models.ConversionFunction{
					Name:        "ConvertXFloatToIntegers",
					PackageName: "parser",
					PackagePath: "github.com/underbek/datamapper/_test_data/parser",
					TypeParam:   models.FromToTypeParam,
				},
				res[models.ConversionFunctionKey{
					FromType: models.Type{Name: tt.FromTypeName},
					ToType:   models.Type{Name: tt.ToTypeName},
				}],
			)
		})
	}
}

func Test_CFParseGenericStruct(t *testing.T) {
	res, err := ParseConversionFunctions(testPath + "generic_struct.go")
	require.NoError(t, err)
	assert.Len(t, res, 2)

	tests := []struct {
		Name        string
		PackagePath string
	}{
		{
			Name:        "ConvertModelsToString",
			PackagePath: "github.com/underbek/datamapper/_test_data/parser",
		},
		{
			Name:        "ConvertModelsToString",
			PackagePath: "github.com/underbek/datamapper/_test_data/other",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			assert.Equal(t,
				models.ConversionFunction{
					Name:        "ConvertModelsToString",
					PackageName: "parser",
					PackagePath: "github.com/underbek/datamapper/_test_data/parser",
					TypeParam:   models.FromTypeParam,
				},
				res[models.ConversionFunctionKey{FromType: models.Type{
					Name:        "Model",
					PackagePath: tt.PackagePath,
				}, ToType: models.Type{Name: "string"}}],
			)
		})
	}
}

func Test_CFParseWithStruct(t *testing.T) {
	res, err := ParseConversionFunctions(testPath + "with_struct.go")
	require.NoError(t, err)
	assert.Len(t, res, 2)

	tests := []struct {
		Name            string
		FromPackagePath string
		ToPackagePath   string
	}{
		{
			Name:            "ConvertCurrentModelToOther",
			FromPackagePath: "github.com/underbek/datamapper/_test_data/parser",
			ToPackagePath:   "github.com/underbek/datamapper/_test_data/other",
		},
		{
			Name:            "ConvertOtherModelToCurrent",
			FromPackagePath: "github.com/underbek/datamapper/_test_data/other",
			ToPackagePath:   "github.com/underbek/datamapper/_test_data/parser",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			assert.Equal(t,
				models.ConversionFunction{
					Name:        tt.Name,
					PackageName: "parser",
					PackagePath: "github.com/underbek/datamapper/_test_data/parser",
				},
				res[models.ConversionFunctionKey{FromType: models.Type{
					Name:        "Model",
					PackagePath: tt.FromPackagePath,
				}, ToType: models.Type{
					Name:        "Model",
					PackagePath: tt.ToPackagePath,
				}}],
			)
		})
	}
}

func Test_CFParseWithError(t *testing.T) {
	res, err := ParseConversionFunctions(testPath + "with_error.go")
	require.NoError(t, err)
	assert.Len(t, res, 6)

	tests := []struct {
		Name          string
		TypeParam     models.TypeParamType
		ToTypes       []string
		ToTypePackage string
	}{
		{
			Name:      "ConvertStringToSigned",
			TypeParam: models.ToTypeParam,
			ToTypes:   []string{"int", "int8", "int16", "int32", "int64"},
		},
		{
			Name:          "ConvertStringToDecimal",
			ToTypes:       []string{"Decimal"},
			ToTypePackage: "github.com/shopspring/decimal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			for _, toTypeName := range tt.ToTypes {
				assert.Equal(t,
					models.ConversionFunction{
						Name:        tt.Name,
						PackageName: "parser",
						PackagePath: "github.com/underbek/datamapper/_test_data/parser",
						TypeParam:   tt.TypeParam,
						WithError:   true,
					},
					res[models.ConversionFunctionKey{FromType: models.Type{
						Name: "string",
					}, ToType: models.Type{
						Name:        toTypeName,
						PackagePath: tt.ToTypePackage,
					}}],
				)
			}
		})
	}
}

func Test_CFParseByPackage(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{
			name:   "Parse by package path",
			source: "github.com/underbek/datamapper/_test_data/mapper/convertors",
		},
		{
			name:   "Parse by sources path",
			source: "../_test_data/mapper/convertors",
		},
		{
			name:   "Parse by one source path",
			source: "../_test_data/mapper/convertors/user_conversion_functions.go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ParseConversionFunctionsByPackage(tt.source)
			require.NoError(t, err)
			assert.Len(t, res, 22)
		})
	}
}
