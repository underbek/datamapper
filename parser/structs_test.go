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
	assert.Len(t, res, 2)
	expected := map[string]models.Struct{
		"TestModel": {Name: "TestModel", Fields: []models.Field{
			{Name: "ID", Type: "int", Tags: []models.Tag{
				{Name: "json", Value: "id"},
				{Name: "map", Value: "id"},
			}},
			{Name: "Name", Type: "string", Tags: []models.Tag{
				{Name: "json", Value: "name"},
				{Name: "map", Value: "name"},
			}},
			{Name: "Empty", Type: "string"},
		}},
		"TestModelTo": {Name: "TestModelTo", Fields: []models.Field{
			{Name: "UUID", Type: "string", Tags: []models.Tag{
				{Name: "db", Value: "uuid"},
				{Name: "map", Value: "id"},
			}},
			{Name: "Name", Type: "string", Tags: []models.Tag{
				{Name: "db", Value: "name"},
				{Name: "map", Value: "name"},
			}},
		}},
	}

	assert.Equal(t, expected, res)
}
