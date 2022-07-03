package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testPath = "test_data/parser/"

func Test_IncorrectFile(t *testing.T) {
	_, err := parseStructs("incorrect name")
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
			res, err := parseStructs(testPath + tt.fileName)
			assert.NoError(t, err)
			assert.Empty(t, res)
		})
	}
}

func Test_ParseModels(t *testing.T) {
	res, err := parseStructs(testPath + "models.go")
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	expected := map[string]Struct{
		"TestModel": {Name: "TestModel", Fields: []Field{
			{Name: "ID", Type: "int", Tags: []Tag{
				{Name: "json", Value: "id"},
				{Name: "map", Value: "id"},
			}},
			{Name: "Name", Type: "string", Tags: []Tag{
				{Name: "json", Value: "name"},
				{Name: "map", Value: "name"},
			}},
			{Name: "Empty", Type: "string"},
		}},
		"TestModelTo": {Name: "TestModelTo", Fields: []Field{
			{Name: "UUID", Type: "string", Tags: []Tag{
				{Name: "db", Value: "uuid"},
				{Name: "map", Value: "id"},
			}},
			{Name: "Name", Type: "string", Tags: []Tag{
				{Name: "db", Value: "name"},
				{Name: "map", Value: "name"},
			}},
		}},
	}

	assert.Equal(t, expected, res)
}
