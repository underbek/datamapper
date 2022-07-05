package generator

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/underbek/datamapper/converts"
	"github.com/underbek/datamapper/models"
)

const (
	testGeneratorPath = "../test_data/generator/"
	generatedPath     = "../generated/"

	modelsFileName = "models.go"
)

func copySource(t *testing.T, fileName string) {
	data, err := ioutil.ReadFile(testGeneratorPath + fileName)
	require.NoError(t, err)

	err = ioutil.WriteFile(generatedPath+fileName, data, 0644)
	require.NoError(t, err)
}

func Test_CreateModelsPair(t *testing.T) {
	fromModel := models.Struct{
		Name: "FromName",
		Fields: []models.Field{
			{Name: "ID", Type: "int", Tags: []models.Tag{{Name: "map", Value: "id"}}},
			{Name: "Name", Type: "string", Tags: []models.Tag{{Name: "map", Value: "name"}}},
		},
	}
	toModel := models.Struct{
		Name: "ToName",
		Fields: []models.Field{
			{Name: "UUID", Type: "string", Tags: []models.Tag{{Name: "map", Value: "id"}}},
			{Name: "Name", Type: "string", Tags: []models.Tag{{Name: "map", Value: "name"}}},
			{Name: "Data", Type: "string", Tags: []models.Tag{{Name: "map", Value: "data"}}},
		},
	}

	g := New(converts.NewFactory())

	res, err := g.createModelsPair(fromModel, toModel)
	require.NoError(t, err)

	expected := result{
		fields: []FieldsPair{
			{
				FromName:   "ID",
				FromType:   "int",
				ToName:     "UUID",
				ToType:     "string",
				Conversion: "fmt.Sprint(from.ID)",
			},
			{
				FromName:   "Name",
				FromType:   "string",
				ToName:     "Name",
				ToType:     "string",
				Conversion: "from.Name",
			},
		},
		imports: []string{"fmt"},
	}

	assert.Equal(t, expected, res)
}

func Test_GenerateConvertorWithoutImports(t *testing.T) {
	fromModel := models.Struct{
		Name: "Model",
		Fields: []models.Field{
			{Name: "Name", Type: "string", Tags: []models.Tag{{Name: "map", Value: "name"}}},
		},
	}
	toModel := models.Struct{
		Name: "DAO",
		Fields: []models.Field{
			{Name: "Name", Type: "string", Tags: []models.Tag{{Name: "map", Value: "name"}}},
		},
	}

	g := New(converts.NewFactory())

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
			{Name: "ID", Type: "int", Tags: []models.Tag{{Name: "map", Value: "id"}}},
			{Name: "Name", Type: "string", Tags: []models.Tag{{Name: "map", Value: "name"}}},
		},
	}
	toModel := models.Struct{
		Name: "DAO",
		Fields: []models.Field{
			{Name: "UUID", Type: "string", Tags: []models.Tag{{Name: "map", Value: "id"}}},
			{Name: "Name", Type: "string", Tags: []models.Tag{{Name: "map", Value: "name"}}},
			{Name: "Data", Type: "string", Tags: []models.Tag{{Name: "map", Value: "data"}}},
		},
	}

	g := New(converts.NewFactory())

	actual, err := g.generateConvertor(fromModel, toModel, "generator")
	require.NoError(t, err)

	expected := `package generator

import "fmt"

func convertModelToDAO(from Model) DAO {
	return DAO{
		UUID: fmt.Sprint(from.ID),
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
			{Name: "ID", Type: "int", Tags: []models.Tag{{Name: "map", Value: "id"}}},
			{Name: "Name", Type: "string", Tags: []models.Tag{{Name: "map", Value: "name"}}},
		},
	}
	toModel := models.Struct{
		Name: "DAO",
		Fields: []models.Field{
			{Name: "UUID", Type: "string", Tags: []models.Tag{{Name: "map", Value: "id"}}},
			{Name: "Name", Type: "string", Tags: []models.Tag{{Name: "map", Value: "name"}}},
			{Name: "Data", Type: "string", Tags: []models.Tag{{Name: "map", Value: "data"}}},
		},
	}

	destination := generatedPath + "simple_convertor.go"

	g := New(converts.NewFactory())

	err := g.CreateConvertor(fromModel, toModel, destination, "generator")
	require.NoError(t, err)

	expected := `package generator

import "fmt"

func convertModelToDAO(from Model) DAO {
	return DAO{
		UUID: fmt.Sprint(from.ID),
		Name: from.Name,
	}
}
`

	actual, err := os.ReadFile(destination)
	require.NoError(t, err)

	assert.Equal(t, expected, string(actual))
}
