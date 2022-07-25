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
			Name: "ConvertIntToString",
			Package: models.Package{
				Name: "parser",
				Path: "github.com/underbek/datamapper/_test_data/parser",
			},
			FromType: models.Type{Name: "int"},
			ToType:   models.Type{Name: "string"},
		},
		res[models.ConversionFunctionKey{FromType: models.Type{Name: "int"}, ToType: models.Type{Name: "string"}}],
	)

	assert.Equal(t,
		models.ConversionFunction{
			Name: "ConvertFloatToString",
			Package: models.Package{
				Name: "parser",
				Path: "github.com/underbek/datamapper/_test_data/parser",
			},
			FromType: models.Type{Name: "float32"},
			ToType:   models.Type{Name: "string"},
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
						Name: tt.Name,
						Package: models.Package{
							Name: "parser",
							Path: "github.com/underbek/datamapper/_test_data/parser",
						},
						TypeParam: models.FromTypeParam,
						FromType:  models.Type{Name: name},
						ToType:    models.Type{Name: "string"},
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
						Name: tt.Name,
						Package: models.Package{
							Name: "parser",
							Path: "github.com/underbek/datamapper/_test_data/parser",
						},
						FromType:  models.Type{Name: "string"},
						ToType:    models.Type{Name: name},
						TypeParam: models.ToTypeParam,
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
					Name: "ConvertXFloatToIntegers",
					Package: models.Package{
						Name: "parser",
						Path: "github.com/underbek/datamapper/_test_data/parser",
					},
					FromType:  models.Type{Name: tt.FromTypeName},
					ToType:    models.Type{Name: tt.ToTypeName},
					TypeParam: models.FromToTypeParam,
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
		Name    string
		Package models.Package
	}{
		{
			Name: "ConvertModelsToString",
			Package: models.Package{
				Name: "parser",
				Path: "github.com/underbek/datamapper/_test_data/parser",
			},
		},
		{
			Name: "ConvertModelsToString",
			Package: models.Package{
				Name: "other",
				Path: "github.com/underbek/datamapper/_test_data/other",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			assert.Equal(t,
				models.ConversionFunction{
					Name: "ConvertModelsToString",
					Package: models.Package{
						Name: "parser",
						Path: "github.com/underbek/datamapper/_test_data/parser",
					},
					TypeParam: models.FromTypeParam,
					FromType:  models.Type{Name: "Model", Package: tt.Package},
					ToType:    models.Type{Name: "string"},
				},
				res[models.ConversionFunctionKey{FromType: models.Type{
					Name:    "Model",
					Package: tt.Package,
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
		Name        string
		FromPackage models.Package
		ToPackage   models.Package
	}{
		{
			Name: "ConvertCurrentModelToOther",
			FromPackage: models.Package{
				Name: "parser",
				Path: "github.com/underbek/datamapper/_test_data/parser",
			},
			ToPackage: models.Package{
				Name: "other",
				Path: "github.com/underbek/datamapper/_test_data/other",
			},
		},
		{
			Name: "ConvertOtherModelToCurrent",
			FromPackage: models.Package{
				Name: "other",
				Path: "github.com/underbek/datamapper/_test_data/other",
			},
			ToPackage: models.Package{
				Name: "parser",
				Path: "github.com/underbek/datamapper/_test_data/parser",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			assert.Equal(t,
				models.ConversionFunction{
					Name: tt.Name,
					Package: models.Package{
						Name: "parser",
						Path: "github.com/underbek/datamapper/_test_data/parser",
					},
					FromType: models.Type{Name: "Model", Package: tt.FromPackage},
					ToType:   models.Type{Name: "Model", Package: tt.ToPackage},
				},
				res[models.ConversionFunctionKey{FromType: models.Type{
					Name:    "Model",
					Package: tt.FromPackage,
				}, ToType: models.Type{
					Name:    "Model",
					Package: tt.ToPackage,
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
		Name      string
		TypeParam models.TypeParamType
		ToTypes   []string
		ToPackage models.Package
	}{
		{
			Name:      "ConvertStringToSigned",
			TypeParam: models.ToTypeParam,
			ToTypes:   []string{"int", "int8", "int16", "int32", "int64"},
		},
		{
			Name:    "ConvertStringToDecimal",
			ToTypes: []string{"Decimal"},
			ToPackage: models.Package{
				Name: "decimal",
				Path: "github.com/shopspring/decimal",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			for _, toTypeName := range tt.ToTypes {
				assert.Equal(t,
					models.ConversionFunction{
						Name: tt.Name,
						Package: models.Package{
							Name: "parser",
							Path: "github.com/underbek/datamapper/_test_data/parser",
						},
						FromType:  models.Type{Name: "string"},
						ToType:    models.Type{Name: toTypeName, Package: tt.ToPackage},
						TypeParam: tt.TypeParam,
						WithError: true,
					},
					res[models.ConversionFunctionKey{FromType: models.Type{
						Name: "string",
					}, ToType: models.Type{
						Name:    toTypeName,
						Package: tt.ToPackage,
					}}],
				)
			}
		})
	}
}

func Test_CFParseWithPointers(t *testing.T) {
	res, err := ParseConversionFunctions(testPath + "with_pointers.go")
	require.NoError(t, err)
	assert.Len(t, res, 5)

	tests := []struct {
		Name      string
		FromType  models.Type
		ToType    models.Type
		TypeParam models.TypeParamType
	}{
		{
			Name: "ConvertIntPtrToString",
			FromType: models.Type{
				Name:    "int",
				Pointer: true,
			},
			ToType: models.Type{
				Name: "string",
			},
		},
		{
			Name: "ConvertFloatToStringPtr",
			FromType: models.Type{
				Name: "float32",
			},
			ToType: models.Type{
				Name:    "string",
				Pointer: true,
			},
		},
		{
			Name: "ConvertFloatPtrToStringPtr",
			FromType: models.Type{
				Name:    "float32",
				Pointer: true,
			},
			ToType: models.Type{
				Name:    "string",
				Pointer: true,
			},
		},
		{
			Name:      "ConvertXFloatPointerToDecimal",
			TypeParam: models.FromTypeParam,
			FromType: models.Type{
				Name:    "float32",
				Pointer: true,
			},
			ToType: models.Type{
				Name: "Decimal",
				Package: models.Package{
					Name: "decimal",
					Path: "github.com/shopspring/decimal",
				},
			},
		},
		{
			Name:      "ConvertXFloatPointerToDecimal",
			TypeParam: models.FromTypeParam,
			FromType: models.Type{
				Name:    "float64",
				Pointer: true,
			},
			ToType: models.Type{
				Name: "Decimal",
				Package: models.Package{
					Name: "decimal",
					Path: "github.com/shopspring/decimal",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			cf, ok := res[models.ConversionFunctionKey{FromType: tt.FromType, ToType: tt.ToType}]
			assert.True(t, ok)
			assert.Equal(t,
				models.ConversionFunction{
					Name:      tt.Name,
					TypeParam: tt.TypeParam,
					Package: models.Package{
						Name: "parser",
						Path: "github.com/underbek/datamapper/_test_data/parser",
					},
					FromType: tt.FromType,
					ToType:   tt.ToType,
				}, cf)
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
			assert.Len(t, res, 11)
		})
	}
}
