package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/underbek/datamapper/models"
)

const testPath = "../test_data/parser/"

func Test_IncorrectFile(t *testing.T) {
	_, err := ParseStructs("incorrect name")
	require.Error(t, err)
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ParseStructs(testPath + tt.fileName)
			assert.NoError(t, err)
			assert.Empty(t, res)
		})
	}
}

func Test_ParseModels(t *testing.T) {
	res, err := ParseStructs(testPath + "models.go")
	assert.NoError(t, err)
	assert.Len(t, res, 3)
	expected := map[string]models.Struct{
		"Model": {Name: "Model", Fields: []models.Field{
			{Name: "ID", Type: models.Type{Name: "string"}},
		}, PackageName: "parser", PackagePath: "github.com/underbek/datamapper/test_data/parser"},
		"TestModel": {Name: "TestModel", Fields: []models.Field{
			{Name: "ID", Type: models.Type{Name: "int"}, Tags: []models.Tag{
				{Name: "json", Value: "id"},
				{Name: "map", Value: "id"},
			}},
			{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{
				{Name: "json", Value: "name"},
				{Name: "map", Value: "name"},
			}},
			{Name: "Empty", Type: models.Type{Name: "string"}},
		}, PackageName: "parser", PackagePath: "github.com/underbek/datamapper/test_data/parser"},
		"TestModelTo": {Name: "TestModelTo", Fields: []models.Field{
			{Name: "UUID", Type: models.Type{Name: "string"}, Tags: []models.Tag{
				{Name: "db", Value: "uuid"},
				{Name: "map", Value: "id"},
			}},
			{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{
				{Name: "db", Value: "name"},
				{Name: "map", Value: "name"},
			}},
		}, PackageName: "parser", PackagePath: "github.com/underbek/datamapper/test_data/parser"},
	}

	assert.Equal(t, expected, res)
}

func Test_ParseComplexModel(t *testing.T) {
	res, err := ParseStructs(testPath + "complex_model.go")
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	expected := map[string]models.Struct{
		"ComplexModel": {Name: "ComplexModel", Fields: []models.Field{
			{
				Name: "ID", Type: models.Type{
					Name:        "Model",
					PackagePath: "github.com/underbek/datamapper/test_data/parser",
				},
				Tags: []models.Tag{
					{Name: "json", Value: "id"},
					{Name: "map", Value: "id"},
				},
			},
			{
				Name: "Age",
				Type: models.Type{
					Name:        "Decimal",
					PackagePath: "github.com/shopspring/decimal",
				},
				Tags: []models.Tag{
					{Name: "json", Value: "age"},
					{Name: "map", Value: "age"},
				},
			},
		}, PackageName: "parser", PackagePath: "github.com/underbek/datamapper/test_data/parser"},
	}

	assert.Equal(t, expected, res)
}
