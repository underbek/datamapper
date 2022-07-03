package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testGeneratorPath = "test_data/generator/"
	generatedPath     = "generated/"
)

func Test_CreateModelsPair(t *testing.T) {
	fromModel := Struct{
		Name: "FromName",
		Fields: []Field{
			{Name: "ID", Type: "int", Tags: []Tag{{Name: "map", Value: "id"}}},
			{Name: "Name", Type: "string", Tags: []Tag{{Name: "map", Value: "name"}}},
		},
	}
	toModel := Struct{
		Name: "ToName",
		Fields: []Field{
			{Name: "UUID", Type: "string", Tags: []Tag{{Name: "map", Value: "id"}}},
			{Name: "Name", Type: "string", Tags: []Tag{{Name: "map", Value: "name"}}},
			{Name: "Data", Type: "string", Tags: []Tag{{Name: "map", Value: "data"}}},
		},
	}

	res := createModelsPair(fromModel, toModel)
	expected := []FieldsPair{
		{
			FromName:   "ID",
			FromType:   "int",
			ToName:     "UUID",
			ToType:     "string",
			Conversion: "from.ID",
		},
		{
			FromName:   "Name",
			FromType:   "string",
			ToName:     "Name",
			ToType:     "string",
			Conversion: "from.Name",
		},
	}

	assert.Equal(t, expected, res)
}

func Test_CreateSimpleConvertor(t *testing.T) {
	fromModel := Struct{
		Name: "Model",
		Fields: []Field{
			{Name: "Name", Type: "string", Tags: []Tag{{Name: "map", Value: "name"}}},
		},
	}
	toModel := Struct{
		Name: "DAO",
		Fields: []Field{
			{Name: "Name", Type: "string", Tags: []Tag{{Name: "map", Value: "name"}}},
		},
	}

	destination := generatedPath + "simple_convertor.go"

	err := createConvertor(fromModel, toModel, destination, "generator")
	require.NoError(t, err)

	expected, err := os.ReadFile(testGeneratorPath + "expected_simple.go")
	require.NoError(t, err)

	actual, err := os.ReadFile(destination)
	require.NoError(t, err)

	assert.Equal(t, string(expected), string(actual))
}
