package generator

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/underbek/datamapper/models"
	"github.com/underbek/datamapper/parser"
)

const (
	testGeneratorPath = "../test_data/generator/"
	generatedPath     = "../generated/"

	modelsFileName = "models.go"

	simpleConvertsSource  = "../converts/simple.go"
	decimalConvertsSource = "../converts/decimal.go"

	generatedPackagePath = "github.com/underbek/datamapper/generated"
	generatedPackageName = "generated"

	otherPackagePath = "github.com/underbek/datamapper/test_data/other"
	otherPackageName = "other"
)

func copySource(t *testing.T, fileName string) {
	data, err := ioutil.ReadFile(testGeneratorPath + fileName)
	require.NoError(t, err)

	err = ioutil.WriteFile(generatedPath+fileName, data, 0644)
	require.NoError(t, err)
}

func parseFunctions(t *testing.T, source string) models.Functions {
	funcs, err := parser.ParseConversionFunctions(source)
	require.NoError(t, err)
	return funcs
}

func Test_CreateModelsPair(t *testing.T) {
	fromModel := models.Struct{
		Name:        "FromName",
		PackageName: generatedPackageName,
		PackagePath: generatedPackagePath,
		Fields: []models.Field{
			{Name: "ID", Type: models.Type{Name: "int"}, Tags: []models.Tag{{Name: "map", Value: "id"}}},
			{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "name"}}},
			{Name: "Age", Type: models.Type{Name: "float64"}, Tags: []models.Tag{{Name: "map", Value: "age"}}},
		},
	}
	toModel := models.Struct{
		Name:        "ToName",
		PackageName: generatedPackageName,
		PackagePath: generatedPackagePath,
		Fields: []models.Field{
			{Name: "UUID", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "id"}}},
			{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "name"}}},
			{Name: "Data", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "data"}}},
			{Name: "Age", Type: models.Type{Name: "uint8"}, Tags: []models.Tag{{Name: "map", Value: "age"}}},
		},
	}

	g := New(parseFunctions(t, simpleConvertsSource))

	res, err := g.createModelsPair(fromModel, toModel)
	require.NoError(t, err)

	expected := result{
		fields: []FieldsPair{
			{
				FromName:   "ID",
				FromType:   "int",
				ToName:     "UUID",
				ToType:     "string",
				Conversion: "converts.ConvertNumericToString(from.ID)",
			},
			{
				FromName:   "Name",
				FromType:   "string",
				ToName:     "Name",
				ToType:     "string",
				Conversion: "from.Name",
			},
			{
				FromName:   "Age",
				FromType:   "float64",
				ToName:     "Age",
				ToType:     "uint8",
				Conversion: "converts.ConvertOrderedToOrdered[float64,uint8](from.Age)",
			},
		},
		imports: []string{"github.com/underbek/datamapper/converts"},
	}

	assert.Equal(t, expected, res)
}

func Test_GenerateConvertorWithoutImports(t *testing.T) {
	fromModel := models.Struct{
		Name:        "Model",
		PackageName: generatedPackageName,
		PackagePath: generatedPackagePath,
		Fields: []models.Field{
			{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "name"}}},
		},
	}
	toModel := models.Struct{
		Name:        "DAO",
		PackageName: generatedPackageName,
		PackagePath: generatedPackagePath,
		Fields: []models.Field{
			{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "name"}}},
		},
	}

	destination := generatedPath + "simple_convertor.go"

	g := New(parseFunctions(t, simpleConvertsSource))

	actual, err := g.generateConvertor(fromModel, toModel, destination)
	require.NoError(t, err)

	expected := `package generator

func ConvertModelToDAO(from Model) DAO {
	return DAO{
		Name: from.Name,
	}
}
`

	assert.Equal(t, expected, string(actual))
}

func Test_GenerateConvertorWithOneImport(t *testing.T) {
	fromModel := models.Struct{
		Name:        "Model",
		PackageName: generatedPackageName,
		PackagePath: generatedPackagePath,
		Fields: []models.Field{
			{Name: "ID", Type: models.Type{Name: "int"}, Tags: []models.Tag{{Name: "map", Value: "id"}}},
			{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "name"}}},
		},
	}
	toModel := models.Struct{
		Name:        "DAO",
		PackageName: generatedPackageName,
		PackagePath: generatedPackagePath,
		Fields: []models.Field{
			{Name: "UUID", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "id"}}},
			{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "name"}}},
			{Name: "Data", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "data"}}},
		},
	}

	destination := generatedPath + "simple_convertor.go"

	g := New(parseFunctions(t, simpleConvertsSource))

	actual, err := g.generateConvertor(fromModel, toModel, destination)
	require.NoError(t, err)

	expected := `package generator

import "github.com/underbek/datamapper/converts"

func ConvertModelToDAO(from Model) DAO {
	return DAO{
		UUID: converts.ConvertNumericToString(from.ID),
		Name: from.Name,
	}
}
`

	assert.Equal(t, expected, string(actual))
}

func Test_CreateConvertorInPackage(t *testing.T) {
	copySource(t, modelsFileName)

	fromModel := models.Struct{
		Name:        "Model",
		PackageName: generatedPackageName,
		PackagePath: generatedPackagePath,
		Fields: []models.Field{
			{Name: "ID", Type: models.Type{Name: "int"}, Tags: []models.Tag{{Name: "map", Value: "id"}}},
			{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "name"}}},
			{Name: "Age", Type: models.Type{Name: "float64"}, Tags: []models.Tag{{Name: "map", Value: "age"}}},
		},
	}
	toModel := models.Struct{
		Name:        "DAO",
		PackageName: generatedPackageName,
		PackagePath: generatedPackagePath,
		Fields: []models.Field{
			{Name: "UUID", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "id"}}},
			{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "name"}}},
			{Name: "Data", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "data"}}},
			{Name: "Age", Type: models.Type{Name: "uint8"}, Tags: []models.Tag{{Name: "map", Value: "age"}}},
		},
	}

	destination := generatedPath + "simple_convertor.go"

	g := New(parseFunctions(t, simpleConvertsSource))

	err := g.CreateConvertor(fromModel, toModel, destination)
	require.NoError(t, err)

	expected := `package generator

import "github.com/underbek/datamapper/converts"

func ConvertModelToDAO(from Model) DAO {
	return DAO{
		UUID: converts.ConvertNumericToString(from.ID),
		Name: from.Name,
		Age:  converts.ConvertOrderedToOrdered[float64, uint8](from.Age),
	}
}
`

	actual, err := os.ReadFile(destination)
	require.NoError(t, err)

	assert.Equal(t, expected, string(actual))
}

func Test_CreateConvertorByOtherPackage(t *testing.T) {
	fromModel := models.Struct{
		Name:        "Model",
		PackageName: otherPackageName,
		PackagePath: otherPackagePath,
		Fields: []models.Field{
			{Name: "ID", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "id"}}},
			{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "name"}}},
			{Name: "Age", Type: models.Type{Name: "float64"}, Tags: []models.Tag{{Name: "map", Value: "age"}}},
		},
	}
	toModel := models.Struct{
		Name:        "DAO",
		PackageName: generatedPackageName,
		PackagePath: generatedPackagePath,
		Fields: []models.Field{
			{Name: "UUID", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "id"}}},
			{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "name"}}},
			{Name: "Data", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "data"}}},
			{Name: "Age", Type: models.Type{Name: "uint8"}, Tags: []models.Tag{{Name: "map", Value: "age"}}},
		},
	}

	destination := generatedPath + "convertor_by_other_package.go"

	g := New(parseFunctions(t, simpleConvertsSource))

	err := g.CreateConvertor(fromModel, toModel, destination)
	require.NoError(t, err)

	expected := `package generator

import (
	"github.com/underbek/datamapper/converts"

	"github.com/underbek/datamapper/test_data/other"
)

func ConvertOtherModelToDAO(from other.Model) DAO {
	return DAO{
		UUID: from.ID,
		Name: from.Name,
		Age:  converts.ConvertOrderedToOrdered[float64, uint8](from.Age),
	}
}
`

	actual, err := os.ReadFile(destination)
	require.NoError(t, err)

	assert.Equal(t, expected, string(actual))
}

func Test_CreateConvertorByComplexModel(t *testing.T) {
	fromModel := models.Struct{
		Name:        "Model",
		PackageName: otherPackageName,
		PackagePath: otherPackagePath,
		Fields: []models.Field{
			{Name: "ID", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "id"}}},
			{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "name"}}},
			{Name: "Age", Type: models.Type{Name: "float64"}, Tags: []models.Tag{{Name: "map", Value: "age"}}},
		},
	}
	toModel := models.Struct{
		Name:        "DTO",
		PackageName: generatedPackageName,
		PackagePath: generatedPackagePath,
		Fields: []models.Field{
			{Name: "ID", Type: models.Type{
				Name:        "Decimal",
				PackagePath: "github.com/shopspring/decimal",
			}, Tags: []models.Tag{{Name: "map", Value: "id"}}},
			{Name: "Age", Type: models.Type{
				Name:        "Decimal",
				PackagePath: "github.com/shopspring/decimal",
			}, Tags: []models.Tag{{Name: "map", Value: "age"}}},
		},
	}

	destination := generatedPath + "convertor_with_complex_model.go"

	g := New(parseFunctions(t, decimalConvertsSource))

	err := g.CreateConvertor(fromModel, toModel, destination)
	require.NoError(t, err)

	expected := `package generator

import (
	"github.com/underbek/datamapper/converts"

	"github.com/underbek/datamapper/test_data/other"
)

func ConvertOtherModelToDTO(from other.Model) DTO {
	return DTO{
		ID:  converts.ConvertStringToDecimal(from.ID),
		Age: converts.ConvertFloatToDecimal(from.Age),
	}
}
`

	actual, err := os.ReadFile(destination)
	require.NoError(t, err)

	assert.Equal(t, expected, string(actual))
}
