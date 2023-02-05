package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/underbek/datamapper/logger"
	"github.com/underbek/datamapper/models"
)

const testPath = "../_test_data/parser/"

func Test_IncorrectFile(t *testing.T) {
	_, err := ParseModels(logger.New(), "incorrect name")
	require.NoError(t, err)
}

func Test_ParseEmptyFile(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
	}{
		{
			name:     "Empty file",
			fileName: "empty_file.go",
		},
		{
			name:     "Empty structs",
			fileName: "empty_models.go",
		},
	}

	lg := logger.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ParseModels(lg, testPath+tt.fileName)
			assert.NoError(t, err)
			assert.Empty(t, res)
		})
	}
}

func Test_ParseModelsWithFunc(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		expected map[string]models.Struct
	}{
		{
			name:     "With function",
			fileName: "with_func.go",
			expected: map[string]models.Struct{
				"StructWithFunc": {
					Type: models.Type{
						Name:    "StructWithFunc",
						Package: models.Package{Name: "parser", Path: "github.com/underbek/datamapper/_test_data/parser"},
						Kind:    models.StructType,
					},
					Fields: []models.Field{},
				},
			},
		},
	}

	lg := logger.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ParseModels(lg, testPath+tt.fileName)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, res)
		})
	}
}

func Test_ParseModels(t *testing.T) {
	res, err := ParseModels(logger.New(), testPath+"models.go")
	assert.NoError(t, err)
	assert.Len(t, res, 3)
	expected := map[string]models.Struct{
		"Model": {
			Type: models.Type{
				Name:    "Model",
				Package: models.Package{Name: "parser", Path: "github.com/underbek/datamapper/_test_data/parser"},
				Kind:    models.StructType,
			},
			Fields: []models.Field{
				{Name: "ID", Type: models.Type{Name: "string"}},
			}},
		"TestModel": {
			Type: models.Type{
				Name:    "TestModel",
				Package: models.Package{Name: "parser", Path: "github.com/underbek/datamapper/_test_data/parser"},
				Kind:    models.StructType,
			}, Fields: []models.Field{
				{Name: "ID", Type: models.Type{Name: "int"}, Tags: []models.Tag{
					{Name: "json", Value: "id"},
					{Name: "map", Value: "id"},
				}},
				{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{
					{Name: "json", Value: "name"},
					{Name: "map", Value: "name"},
				}},
				{Name: "Empty", Type: models.Type{Name: "string"}},
			}},
		"TestModelTo": {
			Type: models.Type{
				Name:    "TestModelTo",
				Package: models.Package{Name: "parser", Path: "github.com/underbek/datamapper/_test_data/parser"},
				Kind:    models.StructType,
			},
			Fields: []models.Field{
				{Name: "UUID", Type: models.Type{Name: "string"}, Tags: []models.Tag{
					{Name: "db", Value: "uuid"},
					{Name: "map", Value: "id"},
				}},
				{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{
					{Name: "db", Value: "name"},
					{Name: "map", Value: "name"},
				}},
			}},
	}

	assert.Equal(t, expected, res)
}

func Test_ParseComplexModel(t *testing.T) {
	res, err := ParseModels(logger.New(), testPath+"complex_model.go")
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	expected := map[string]models.Struct{
		"ComplexModel": {
			Type: models.Type{
				Name:    "ComplexModel",
				Package: models.Package{Name: "parser", Path: "github.com/underbek/datamapper/_test_data/parser"},
				Kind:    models.StructType,
			},
			Fields: []models.Field{
				{
					Name: "ID", Type: models.Type{
						Name: "Model",
						Package: models.Package{
							Name: "parser",
							Path: "github.com/underbek/datamapper/_test_data/parser",
						},
						Kind: models.StructType,
					},
					Tags: []models.Tag{
						{Name: "json", Value: "id"},
						{Name: "map", Value: "id"},
					},
				},
				{
					Name: "Age",
					Type: models.Type{
						Name: "Decimal",
						Package: models.Package{
							Name: "decimal",
							Path: "github.com/shopspring/decimal",
						},
						Kind: models.StructType,
					},
					Tags: []models.Tag{
						{Name: "json", Value: "age"},
						{Name: "map", Value: "age"},
					},
				},
			}},
	}

	assert.Equal(t, expected, res)
}

func Test_ParseModelWithPointerField(t *testing.T) {
	res, err := ParseModels(logger.New(), testPath+"pointer_model.go")
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	expected := map[string]models.Struct{
		"PointerModel": {
			Type: models.Type{
				Name:    "PointerModel",
				Package: models.Package{Name: "parser", Path: "github.com/underbek/datamapper/_test_data/parser"},
				Kind:    models.StructType,
			}, Fields: []models.Field{
				{
					Name: "ID",
					Type: models.Type{Name: "int", Pointer: true},
					Tags: []models.Tag{{Name: "map", Value: "id"}},
				},
				{
					Name: "Name",
					Type: models.Type{Name: "string", Pointer: true},
					Tags: []models.Tag{{Name: "map", Value: "name"}},
				},
				{
					Name: "Age",
					Type: models.Type{
						Name: "Decimal",
						Package: models.Package{
							Name: "decimal",
							Path: "github.com/shopspring/decimal",
						},
						Kind:    models.StructType,
						Pointer: true,
					},
					Tags: []models.Tag{{Name: "map", Value: "age"}},
				},
			}},
	}

	assert.Equal(t, expected, res)
}

func Test_ParseModelByPackage(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{
			name:   "Parse by package path",
			source: "github.com/underbek/datamapper/_test_data/mapper/transport",
		},
		{
			name:   "Parse by sources path",
			source: "../_test_data/mapper/transport",
		},
		{
			name:   "Parse by one source path",
			source: "../_test_data/mapper/transport/models.go",
		},
	}

	lg := logger.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ParseModelsByPackage(lg, tt.source)
			require.NoError(t, err)
			_, ok := res["User"]
			assert.True(t, ok)
		})
	}
}

func Test_ParseModelWithAlias(t *testing.T) {
	res, err := ParseModels(logger.New(), testPath+"alias_model.go")
	require.NoError(t, err)
	assert.Len(t, res, 2)

	expected := models.Struct{
		Type: models.Type{
			Name:    "WithAlias",
			Package: models.Package{Name: "parser", Path: "github.com/underbek/datamapper/_test_data/parser"},
			Kind:    models.StructType,
		},
		Fields: []models.Field{
			{
				Name: "String",
				Type: models.Type{
					Name:    "String",
					Package: models.Package{Name: "parser", Path: "github.com/underbek/datamapper/_test_data/parser"},
					Kind:    models.RedefinedType,
				},
			},
			{
				Name: "StringAlias",
				Type: models.Type{Name: "string"},
			},
			{
				Name: "Array",
				Type: models.Type{
					Name:    "Array",
					Package: models.Package{Name: "parser", Path: "github.com/underbek/datamapper/_test_data/parser"},
					Kind:    models.RedefinedType,
				},
			},
			{
				Name: "Slice",
				Type: models.Type{
					Name:    "Slice",
					Package: models.Package{Name: "parser", Path: "github.com/underbek/datamapper/_test_data/parser"},
					Kind:    models.RedefinedType,
				},
			},
			{
				Name: "Map",
				Type: models.Type{
					Name:    "Map",
					Package: models.Package{Name: "parser", Path: "github.com/underbek/datamapper/_test_data/parser"},
					Kind:    models.RedefinedType,
				},
			},
			{
				Name: "RawArray",
				Type: models.Type{
					Kind: models.ArrayType,
					Additional: models.ArrayAdditional{
						InType: models.Type{
							Name: "uint",
						},
						Len: 16,
					},
				},
			},
			{
				Name: "RawSlice",
				Type: models.Type{
					Kind: models.SliceType,
					Additional: models.SliceAdditional{
						InType: models.Type{
							Name: "uint",
						},
					},
				},
			},
			{
				Name: "RawMap",
				Type: models.Type{
					Kind: models.MapType,
					Additional: models.MapAdditional{
						KeyType: models.Type{
							Name: "int",
						},
						ValueType: models.Type{
							Name: "string",
						},
					},
				},
			},
			{
				Name: "ModelRedefinition",
				Type: models.Type{
					Name:    "ModelRedefinition",
					Package: models.Package{Name: "parser", Path: "github.com/underbek/datamapper/_test_data/parser"},
					Kind:    models.StructType,
				},
			},
			{
				Name: "UUID",
				Type: models.Type{
					Name:    "UUID",
					Package: models.Package{Name: "uuid", Path: "github.com/google/uuid"},
					Kind:    models.RedefinedType,
				},
			},
		},
	}

	assert.Equal(t, expected, res["WithAlias"])
}

func Test_ParseModelWithCollections(t *testing.T) {
	res, err := ParseModels(logger.New(), testPath+"model_with_collections.go")
	require.NoError(t, err)
	assert.Len(t, res, 1)

	expected := models.Struct{
		Type: models.Type{
			Name:    "ModelWithCollections",
			Package: models.Package{Name: "parser", Path: "github.com/underbek/datamapper/_test_data/parser"},
			Kind:    models.StructType,
		},
		Fields: []models.Field{
			{
				Name: "Array",
				Type: models.Type{
					Kind: models.ArrayType,
					Additional: models.ArrayAdditional{
						InType: models.Type{
							Name: "uint64",
						},
						Len: 12,
					},
				},
			},
			{
				Name: "Slice",
				Type: models.Type{
					Kind: models.SliceType,
					Additional: models.SliceAdditional{
						InType: models.Type{
							Name: "string",
						},
					},
				},
			},
			{
				Name: "Map",
				Type: models.Type{
					Kind: models.MapType,
					Additional: models.MapAdditional{
						KeyType: models.Type{
							Name: "int",
						},
						ValueType: models.Type{
							Name: "string",
						},
					},
				},
			},
			{
				Name: "PointerArray",
				Type: models.Type{
					Kind:    models.ArrayType,
					Pointer: true,
					Additional: models.ArrayAdditional{
						InType: models.Type{
							Name: "uint64",
						},
						Len: 12,
					},
				},
			},
			{
				Name: "PointerSlice",
				Type: models.Type{
					Kind:    models.SliceType,
					Pointer: true,
					Additional: models.SliceAdditional{
						InType: models.Type{
							Name: "string",
						},
					},
				},
			},
			{
				Name: "PointerMap",
				Type: models.Type{
					Kind:    models.MapType,
					Pointer: true,
					Additional: models.MapAdditional{
						KeyType: models.Type{
							Name: "int",
						},
						ValueType: models.Type{
							Name: "string",
						},
					},
				},
			},
			{
				Name: "ArrayPointers",
				Type: models.Type{
					Kind: models.ArrayType,
					Additional: models.ArrayAdditional{
						InType: models.Type{
							Name:    "uint64",
							Pointer: true,
						},
						Len: 12,
					},
				},
			},
			{
				Name: "SlicePointers",
				Type: models.Type{
					Kind: models.SliceType,
					Additional: models.SliceAdditional{
						InType: models.Type{
							Name:    "string",
							Pointer: true,
						},
					},
				},
			},
			{
				Name: "MapPointers",
				Type: models.Type{
					Kind: models.MapType,
					Additional: models.MapAdditional{
						KeyType: models.Type{
							Name:    "int",
							Pointer: true,
						},
						ValueType: models.Type{
							Name:    "string",
							Pointer: true,
						},
					},
				},
			},
			{
				Name: "ArrayModel",
				Type: models.Type{
					Kind: models.ArrayType,
					Additional: models.ArrayAdditional{
						InType: models.Type{
							Name: "Model",
							Package: models.Package{
								Name: "parser",
								Path: "github.com/underbek/datamapper/_test_data/parser",
							},
							Kind: models.StructType,
						},
						Len: 12,
					},
				},
			},
			{
				Name: "SliceModel",
				Type: models.Type{
					Kind: models.SliceType,
					Additional: models.SliceAdditional{
						InType: models.Type{
							Name: "Model",
							Package: models.Package{
								Name: "parser",
								Path: "github.com/underbek/datamapper/_test_data/parser",
							},
							Kind: models.StructType,
						},
					},
				},
			},
			{
				Name: "MapModel",
				Type: models.Type{
					Kind: models.MapType,
					Additional: models.MapAdditional{
						KeyType: models.Type{
							Name: "Model",
							Package: models.Package{
								Name: "parser",
								Path: "github.com/underbek/datamapper/_test_data/parser",
							},
							Kind: models.StructType,
						},
						ValueType: models.Type{
							Name: "Model",
							Package: models.Package{
								Name: "parser",
								Path: "github.com/underbek/datamapper/_test_data/parser",
							},
							Kind: models.StructType,
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expected, res["ModelWithCollections"])
}

func Test_ParseModelByBrokenPackage(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{
			name:   "Parse by package path",
			source: "github.com/underbek/datamapper/_test_data/mapper/broken",
		},
		{
			name:   "Parse by sources path",
			source: "../_test_data/mapper/broken",
		},
		{
			name:   "Parse by one source path",
			source: "../_test_data/mapper/broken/models.go",
		},
	}

	lg := logger.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ParseModelsByPackage(lg, tt.source)
			require.NoError(t, err)
			_, ok := res["User"]
			assert.True(t, ok)
			assert.Len(t, res, 2)
		})
	}
}
