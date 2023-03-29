package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/underbek/datamapper/_test_data"
	"github.com/underbek/datamapper/logger"
	"github.com/underbek/datamapper/models"
	"github.com/underbek/datamapper/parser"
)

const (
	testGeneratorPath = "../_test_data/generator/"

	cfPath = "../converts"
	cfFile = "convertor.go"

	generatedPackagePath = "github.com/underbek/datamapper/_test_data/generated/generator"
	generatedPackageName = "generator"

	defaultTag = "map"
)

func parseFunctions(t *testing.T, source string) models.Functions {
	funcs, err := parser.ParseConversionFunctions(logger.New(), source)
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
		Fields: models.NewFields([]models.Field{
			{Name: "ID", Type: models.Type{Name: "int"}, Tags: []models.Tag{{Name: "map", Value: "id"}}},
			{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "name"}}},
			{Name: "Age", Type: models.Type{Name: "float64"}, Tags: []models.Tag{{Name: "map", Value: "age"}}},
		}),
	}

	toType := models.Type{
		Name: "ToName",
		Package: models.Package{
			Name: generatedPackageName,
			Path: generatedPackagePath,
		},
	}

	toModel := models.Struct{
		Type: toType,
		Fields: models.NewFields([]models.Field{
			{Name: "UUID", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "id"}}},
			{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "name"}}},
			{Name: "Data", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "data"}}},
			{Name: "Age", Type: models.Type{Name: "uint8"}, Tags: []models.Tag{{Name: "map", Value: "age"}}},
		}),
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
				Types:      []TypeWithName{{Type: toType}},
			},
			{
				FromName:   "Name",
				FromType:   "string",
				ToName:     "Name",
				ToType:     "string",
				Assignment: "from.Name",
				Types:      []TypeWithName{{Type: toType}},
			},
			{
				FromName:   "Age",
				FromType:   "float64",
				ToName:     "Age",
				ToType:     "uint8",
				Assignment: "converts.ConvertOrderedToOrdered[float64,uint8](from.Age)",
				Types:      []TypeWithName{{Type: toType}},
			},
		},
		packages: models.Packages{
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
		name          string
		pathFrom      string
		pathTo        string
		generatePath  string
		cfPath        string
		isFromPointer bool
		isToPointer   bool
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
			name:         "With field pointers",
			pathFrom:     "with_field_pointers",
			pathTo:       "with_field_pointers",
			generatePath: "with_field_pointers",
			cfPath:       cfPath,
		},
		{
			name:         "With field pointers and convertors",
			pathFrom:     "with_field_pointers_and_convertors",
			pathTo:       "with_field_pointers_and_convertors",
			generatePath: "with_field_pointers_and_convertors",
			cfPath:       cfPath,
		},
		{
			name:         "With field pointers without error",
			pathFrom:     "with_field_pointers_without_error",
			pathTo:       "with_field_pointers_without_error",
			generatePath: "with_field_pointers_without_error",
			cfPath:       cfPath,
		},
		{
			name:         "With field pointers and errors",
			pathFrom:     "with_field_pointers_and_errors",
			pathTo:       "with_field_pointers_and_errors",
			generatePath: "with_field_pointers_and_errors",
			cfPath:       cfPath,
		},
		{
			name:         "Conversion functions with pointers",
			pathFrom:     "with_pointers",
			pathTo:       "with_pointers",
			generatePath: "with_pointers",
			cfPath:       testGeneratorPath + "with_pointers/cf",
		},
		{
			name:         "Conversion functions with pointers and errors",
			pathFrom:     "with_pointers_and_errors",
			pathTo:       "with_pointers_and_errors",
			generatePath: "with_pointers_and_errors",
			cfPath:       testGeneratorPath + "with_pointers_and_errors/cf",
		},
		{
			name:         "With field slice",
			pathFrom:     "with_field_slice",
			pathTo:       "with_field_slice",
			generatePath: "with_field_slice",
			cfPath:       testGeneratorPath + "with_field_slice/cf",
		},
		{
			name:         "With field slice and errors",
			pathFrom:     "with_field_slice_and_errors",
			pathTo:       "with_field_slice_and_errors",
			generatePath: "with_field_slice_and_errors",
			cfPath:       testGeneratorPath + "with_field_slice_and_errors/cf",
		},
		{
			name:         "With field slice and pointers",
			pathFrom:     "with_field_slice_and_pointers",
			pathTo:       "with_field_slice_and_pointers",
			generatePath: "with_field_slice_and_pointers",
			cfPath:       testGeneratorPath + "with_field_slice_and_pointers/cf",
		},
		{
			name:         "With field slice pointers and errors",
			pathFrom:     "with_field_slice_pointers_and_errors",
			pathTo:       "with_field_slice_pointers_and_errors",
			generatePath: "with_field_slice_pointers_and_errors",
			cfPath:       testGeneratorPath + "with_field_slice_pointers_and_errors/cf",
		},
		{
			name:          "With from pointer",
			pathFrom:      "with_from_pointer",
			pathTo:        "with_from_pointer",
			generatePath:  "with_from_pointer",
			cfPath:        cfPath,
			isFromPointer: true,
		},
		{
			name:         "With to pointer",
			pathFrom:     "with_to_pointer",
			pathTo:       "with_to_pointer",
			generatePath: "with_to_pointer",
			cfPath:       cfPath,
			isToPointer:  true,
		},
		{
			name:          "With from and to pointers",
			pathFrom:      "with_from_and_to_pointers",
			pathTo:        "with_from_and_to_pointers",
			generatePath:  "with_from_and_to_pointers",
			cfPath:        cfPath,
			isFromPointer: true,
			isToPointer:   true,
		},
		{
			name:          "With from and to pointers without errors",
			pathFrom:      "with_from_and_to_pointers_without_errors",
			pathTo:        "with_from_and_to_pointers_without_errors",
			generatePath:  "with_from_and_to_pointers_without_errors",
			cfPath:        cfPath,
			isFromPointer: true,
			isToPointer:   true,
		},
	}

	lg := logger.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modelsFrom, err := parser.ParseModels(lg, testGeneratorPath+tt.pathFrom+"/models.go")
			require.NoError(t, err)

			modelsTo, err := parser.ParseModels(lg, testGeneratorPath+tt.pathTo+"/models.go")
			require.NoError(t, err)

			funcs := parseFunctions(t, tt.cfPath)

			pkg, err := parser.ParseDestinationPackage(lg, testGeneratorPath+tt.generatePath)
			require.NoError(t, err)

			from := modelsFrom["From"]
			from.Type.Pointer = tt.isFromPointer

			to := modelsTo["To"]
			to.Type.Pointer = tt.isToPointer

			gcf, err := GenerateConvertor(from, to, defaultTag, defaultTag, pkg, funcs)
			require.NoError(t, err)

			actual, err := fillConvertorsSource(pkg, gcf.Packages, []string{gcf.Body})
			require.NoError(t, err)

			expected := _test_data.Generator(t, tt.generatePath+"/convertor.go")
			assert.Equal(t, expected, string(actual))
		})
	}
}

func Test_GenerateConvertorWithAliases(t *testing.T) {
	lg := logger.New()

	modelsFrom, err := parser.ParseModels(lg, testGeneratorPath+"with_aliases/from")
	require.NoError(t, err)

	modelsTo, err := parser.ParseModels(lg, testGeneratorPath+"with_aliases/to")
	require.NoError(t, err)

	funcs := parseFunctions(t, testGeneratorPath+"with_aliases/cf")
	for key, cf := range funcs {
		cf.Package.Alias = "cfalias"
		funcs[key] = cf
	}

	from := modelsFrom["From"]
	from.Type.Package.Alias = "fromalias"

	to := modelsTo["To"]
	to.Type.Package.Alias = "toalias"

	pkg, err := parser.ParseDestinationPackage(lg, testGeneratorPath+"with_aliases")
	require.NoError(t, err)

	gcf, err := GenerateConvertor(from, to, defaultTag, defaultTag, pkg, funcs)
	require.NoError(t, err)

	actual, err := fillConvertorsSource(pkg, gcf.Packages, []string{gcf.Body})
	require.NoError(t, err)

	expected := _test_data.Generator(t, "with_aliases/convertor.go")
	assert.Equal(t, expected, string(actual))
}

func Test_GenerateConvertorWithSlice(t *testing.T) {
	tests := []struct {
		name          string
		pathFrom      string
		pathTo        string
		generatePath  string
		cfPath        string
		fromAlias     string
		toAlias       string
		isFromPointer bool
		isToPointer   bool
	}{
		{
			name:         "Without imports",
			pathFrom:     "without_imports",
			pathTo:       "without_imports",
			generatePath: "without_imports",
			cfPath:       cfFile,
		},
		{
			name:         "Other package model",
			pathFrom:     "other_package_model/other",
			pathTo:       "other_package_model",
			generatePath: "other_package_model",
			cfPath:       cfFile,
		},
		{
			name:         "With error",
			pathFrom:     "with_error",
			pathTo:       "with_error",
			generatePath: "with_error",
			cfPath:       cfFile,
		},
		{
			name:         "With aliases",
			pathFrom:     "with_aliases/from",
			pathTo:       "with_aliases/to",
			generatePath: "with_aliases",
			cfPath:       cfFile,
			fromAlias:    "fromalias",
			toAlias:      "toalias",
		},
		{
			name:          "With from pointer",
			pathFrom:      "with_from_pointer",
			pathTo:        "with_from_pointer",
			generatePath:  "with_from_pointer",
			cfPath:        cfFile,
			isFromPointer: true,
		},
		{
			name:         "With to pointer",
			pathFrom:     "with_to_pointer",
			pathTo:       "with_to_pointer",
			generatePath: "with_to_pointer",
			cfPath:       cfFile,
			isToPointer:  true,
		},
		{
			name:          "With from and to pointers",
			pathFrom:      "with_from_and_to_pointers",
			pathTo:        "with_from_and_to_pointers",
			generatePath:  "with_from_and_to_pointers",
			cfPath:        cfFile,
			isFromPointer: true,
			isToPointer:   true,
		},
		{
			name:          "With from and to pointers without errors",
			pathFrom:      "with_from_and_to_pointers_without_errors",
			pathTo:        "with_from_and_to_pointers_without_errors",
			generatePath:  "with_from_and_to_pointers_without_errors",
			cfPath:        cfFile,
			isFromPointer: true,
			isToPointer:   true,
		},
	}

	lg := logger.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modelsFrom, err := parser.ParseModels(lg, testGeneratorPath+tt.pathFrom+"/models.go")
			require.NoError(t, err)

			modelsTo, err := parser.ParseModels(lg, testGeneratorPath+tt.pathTo+"/models.go")
			require.NoError(t, err)

			funcs := parseFunctions(t, testGeneratorPath+tt.generatePath+"/"+tt.cfPath)

			pkg, err := parser.ParseDestinationPackage(lg, testGeneratorPath+tt.generatePath)
			require.NoError(t, err)

			from := modelsFrom["From"].Type
			from.Pointer = tt.isFromPointer

			to := modelsTo["To"].Type
			to.Pointer = tt.isToPointer

			cf := funcs[models.ConversionFunctionKey{
				FromType: from,
				ToType:   to,
			}]

			from.Package.Alias = tt.fromAlias
			to.Package.Alias = tt.toAlias

			gcf, err := GenerateSliceConvertor(from, to, pkg, cf)
			require.NoError(t, err)

			actual, err := fillConvertorsSource(pkg, gcf.Packages, []string{gcf.Body})
			require.NoError(t, err)

			expected := _test_data.Generator(t, tt.generatePath+"/convertor_with_slice.go")
			assert.Equal(t, expected, string(actual))
		})
	}
}
