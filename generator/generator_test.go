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

	convertsSource = "../converts/simple.go"
)

func copySource(t *testing.T, fileName string) {
	data, err := ioutil.ReadFile(testGeneratorPath + fileName)
	require.NoError(t, err)

	err = ioutil.WriteFile(generatedPath+fileName, data, 0644)
	require.NoError(t, err)
}

func parseFunctions(t *testing.T) models.Functions {
	funcs, err := parser.ParseConversionFunctions(convertsSource)
	require.NoError(t, err)
	return funcs
}

func Test_CreateModelsPair(t *testing.T) {
	fromModel := models.Struct{
		Name: "FromName",
		Fields: []models.Field{
			{Name: "ID", Type: models.Type{Name: "int"}, Tags: []models.Tag{{Name: "map", Value: "id"}}},
			{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "name"}}},
			{Name: "Age", Type: models.Type{Name: "float64"}, Tags: []models.Tag{{Name: "map", Value: "age"}}},
		},
	}
	toModel := models.Struct{
		Name: "ToName",
		Fields: []models.Field{
			{Name: "UUID", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "id"}}},
			{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "name"}}},
			{Name: "Data", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "data"}}},
			{Name: "Age", Type: models.Type{Name: "uint8"}, Tags: []models.Tag{{Name: "map", Value: "age"}}},
		},
	}

	g := New(parseFunctions(t))

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
		Name: "Model",
		Fields: []models.Field{
			{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "name"}}},
		},
	}
	toModel := models.Struct{
		Name: "DAO",
		Fields: []models.Field{
			{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "name"}}},
		},
	}

	g := New(parseFunctions(t))

	actual, err := g.generateConvertor(fromModel, toModel, "generator")
	require.NoError(t, err)

	expected := `package generator

func convertModelToDAO(from Model) DAO {
	return DAO{
		Name: from.Name,
	}
}
`

	assert.Equal(t, expected, string(actual))
}

func Test_GenerateConvertorWithOneImport(t *testing.T) {
	fromModel := models.Struct{
		Name: "Model",
		Fields: []models.Field{
			{Name: "ID", Type: models.Type{Name: "int"}, Tags: []models.Tag{{Name: "map", Value: "id"}}},
			{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "name"}}},
		},
	}
	toModel := models.Struct{
		Name: "DAO",
		Fields: []models.Field{
			{Name: "UUID", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "id"}}},
			{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "name"}}},
			{Name: "Data", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "data"}}},
		},
	}

	g := New(parseFunctions(t))

	actual, err := g.generateConvertor(fromModel, toModel, "generator")
	require.NoError(t, err)

	expected := `package generator

import "github.com/underbek/datamapper/converts"

func convertModelToDAO(from Model) DAO {
	return DAO{
		UUID: converts.ConvertNumericToString(from.ID),
		Name: from.Name,
	}
}
`

	assert.Equal(t, expected, string(actual))
}

func Test_CreateConvertor(t *testing.T) {
	copySource(t, modelsFileName)

	fromModel := models.Struct{
		Name: "Model",
		Fields: []models.Field{
			{Name: "ID", Type: models.Type{Name: "int"}, Tags: []models.Tag{{Name: "map", Value: "id"}}},
			{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "name"}}},
			{Name: "Age", Type: models.Type{Name: "float64"}, Tags: []models.Tag{{Name: "map", Value: "age"}}},
		},
	}
	toModel := models.Struct{
		Name: "DAO",
		Fields: []models.Field{
			{Name: "UUID", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "id"}}},
			{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "name"}}},
			{Name: "Data", Type: models.Type{Name: "string"}, Tags: []models.Tag{{Name: "map", Value: "data"}}},
			{Name: "Age", Type: models.Type{Name: "uint8"}, Tags: []models.Tag{{Name: "map", Value: "age"}}},
		},
	}

	destination := generatedPath + "simple_convertor.go"

	g := New(parseFunctions(t))

	err := g.CreateConvertor(fromModel, toModel, destination, "generator")
	require.NoError(t, err)

	expected := `package generator

import "github.com/underbek/datamapper/converts"

func convertModelToDAO(from Model) DAO {
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
