package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/underbek/datamapper/_test_data"
	"github.com/underbek/datamapper/models"
	"github.com/underbek/datamapper/parser"
)

const (
	testGeneratorPath = "../_test_data/generator/"

	cfPath = "../converts"

	generatedPackagePath = "github.com/underbek/datamapper/_test_data/generated/generator"
	generatedPackageName = "generator"
)

func parseFunctions(t *testing.T, source string) models.Functions {
	funcs, err := parser.ParseConversionFunctions(source)
	require.NoError(t, err)
	return funcs
}

func Test_CreateModelsPair(t *testing.T) {
	fromModel := models.Struct{
		Type: models.Type{
			Name: "FromName",
			Package: models.Package{
				Name: generatedPackageName,
				Path: generatedPackagePath,
			},
		},
		Fields: []models.Field{
			{Name: "ID", Type: models.Type{Name: "int"}, Tags: []models.Tag{{Name: "map", Value: "id"}}},
			{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "name"}}},
			{Name: "Age", Type: models.Type{Name: "float64"}, Tags: []models.Tag{{Name: "map", Value: "age"}}},
		},
	}
	toModel := models.Struct{
		Type: models.Type{
			Name: "ToName",
			Package: models.Package{
				Name: generatedPackageName,
				Path: generatedPackagePath,
			},
		},
		Fields: []models.Field{
			{Name: "UUID", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "id"}}},
			{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "name"}}},
			{Name: "Data", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "data"}}},
			{Name: "Age", Type: models.Type{Name: "uint8"}, Tags: []models.Tag{{Name: "map", Value: "age"}}},
		},
	}

	res, err := createModelsPair(fromModel, toModel, "", parseFunctions(t, cfPath))
	require.NoError(t, err)

	expected := result{
		fields: []FieldsPair{
			{
				FromName:   "ID",
				FromType:   "int",
				ToName:     "UUID",
				ToType:     "string",
				Assignment: "converts.ConvertNumericToString(from.ID)",
			},
			{
				FromName:   "Name",
				FromType:   "string",
				ToName:     "Name",
				ToType:     "string",
				Assignment: "from.Name",
			},
			{
				FromName:   "Age",
				FromType:   "float64",
				ToName:     "Age",
				ToType:     "uint8",
				Assignment: "converts.ConvertOrderedToOrdered[float64,uint8](from.Age)",
			},
		},
		packages: map[models.Package]struct{}{
			{
				Name: "converts",
				Path: "github.com/underbek/datamapper/converts",
			}: {},
		},
	}

	assert.Equal(t, expected, res)
}

func Test_GenerateConvertor(t *testing.T) {
	tests := []struct {
		name         string
		pathFrom     string
		pathTo       string
		generatePath string
		cfPath       string
	}{
		{
			name:         "Without imports",
			pathFrom:     "without_imports",
			pathTo:       "without_imports",
			generatePath: "without_imports",
			cfPath:       cfPath,
		},
		{
			name:         "With one import",
			pathFrom:     "with_one_import",
			pathTo:       "with_one_import",
			generatePath: "with_one_import",
			cfPath:       cfPath,
		},
		{
			name:         "Other package model",
			pathFrom:     "other_package_model/other",
			pathTo:       "other_package_model",
			generatePath: "other_package_model",
			cfPath:       cfPath,
		},
		{
			name:         "Complex model",
			pathFrom:     "complex_model",
			pathTo:       "complex_model",
			generatePath: "complex_model",
			cfPath:       cfPath,
		},
		{
			name:         "Same conversion functions path",
			pathFrom:     "same_cf_path",
			pathTo:       "same_cf_path",
			generatePath: "same_cf_path",
			cfPath:       testGeneratorPath + "same_cf_path/convertors.go",
		},
		{
			name:         "With error",
			pathFrom:     "with_error",
			pathTo:       "with_error",
			generatePath: "with_error",
			cfPath:       cfPath,
		},
		{
			name:         "With some errors",
			pathFrom:     "with_errors",
			pathTo:       "with_errors",
			generatePath: "with_errors",
			cfPath:       cfPath,
		},
		{
			name:         "With filed pointers",
			pathFrom:     "with_field_pointers",
			pathTo:       "with_field_pointers",
			generatePath: "with_field_pointers",
			cfPath:       cfPath,
		},
		{
			name:         "With filed pointers and convertors",
			pathFrom:     "with_filed_pointers_and_convertors",
			pathTo:       "with_filed_pointers_and_convertors",
			generatePath: "with_filed_pointers_and_convertors",
			cfPath:       cfPath,
		},
		{
			name:         "With filed pointers and errors",
			pathFrom:     "with_filed_pointers_and_errors",
			pathTo:       "with_filed_pointers_and_errors",
			generatePath: "with_filed_pointers_and_errors",
			cfPath:       cfPath,
		},
		{
			name:         "Conversion functions with pointers",
			pathFrom:     "cf_with_pointers",
			pathTo:       "cf_with_pointers",
			generatePath: "cf_with_pointers",
			cfPath:       testGeneratorPath + "cf_with_pointers/cf",
		},
		{
			name:         "Conversion functions with pointers and errors",
			pathFrom:     "cf_with_pointers_and_errors",
			pathTo:       "cf_with_pointers_and_errors",
			generatePath: "cf_with_pointers_and_errors",
			cfPath:       testGeneratorPath + "cf_with_pointers_and_errors/cf",
		},
		{
			name:         "With slice",
			pathFrom:     "cf_with_slice",
			pathTo:       "cf_with_slice",
			generatePath: "cf_with_slice",
			cfPath:       testGeneratorPath + "cf_with_slice/cf",
		},
		{
			name:         "With slice and errors",
			pathFrom:     "cf_with_slice_and_errors",
			pathTo:       "cf_with_slice_and_errors",
			generatePath: "cf_with_slice_and_errors",
			cfPath:       testGeneratorPath + "cf_with_slice_and_errors/cf",
		},
		{
			name:         "With slice and pointers",
			pathFrom:     "cf_with_slice_and_pointers",
			pathTo:       "cf_with_slice_and_pointers",
			generatePath: "cf_with_slice_and_pointers",
			cfPath:       testGeneratorPath + "cf_with_slice_and_pointers/cf",
		},
		//TODO: slices with pinter and errors add pointer slices
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modelsFrom, err := parser.ParseModels(testGeneratorPath + tt.pathFrom + "/models.go")
			require.NoError(t, err)

			modelsTo, err := parser.ParseModels(testGeneratorPath + tt.pathTo + "/models.go")
			require.NoError(t, err)

			funcs := parseFunctions(t, tt.cfPath)

			destination := testGeneratorPath + tt.generatePath + "/convertor.go"

			actual, err := generateConvertor(modelsFrom["From"], modelsTo["To"], destination, funcs)
			require.NoError(t, err)

			expected := _test_data.Generator(t, tt.generatePath+"/convertor.go")
			assert.Equal(t, expected, string(actual))
		})
	}
}

func Test_GenerateConvertorWithAliases(t *testing.T) {
	modelsFrom, err := parser.ParseModels(testGeneratorPath + "with_aliases/from")
	require.NoError(t, err)

	modelsTo, err := parser.ParseModels(testGeneratorPath + "with_aliases/to")
	require.NoError(t, err)

	funcs := parseFunctions(t, testGeneratorPath+"with_aliases/cf")
	for key, cf := range funcs {
		cf.Package.Alias = "cfalias"
		funcs[key] = cf
	}

	destination := testGeneratorPath + "with_aliases/convertor.go"

	from := modelsFrom["From"]
	from.Type.Package.Alias = "fromalias"

	to := modelsTo["To"]
	to.Type.Package.Alias = "toalias"

	actual, err := generateConvertor(from, to, destination, funcs)
	require.NoError(t, err)

	expected := _test_data.Generator(t, "with_aliases/convertor.go")
	assert.Equal(t, expected, string(actual))
}
