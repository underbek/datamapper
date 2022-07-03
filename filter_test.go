package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Filter(t *testing.T) {
	fields := []Field{
		{Name: "ID", Type: "int", Tags: []Tag{
			{Name: "json", Value: "id"},
			{Name: "map", Value: "id"},
		}},
		{Name: "Name", Type: "string", Tags: []Tag{
			{Name: "json", Value: "name"},
			{Name: "map", Value: "name"},
		}},
		{Name: "Empty", Type: "string"},
	}

	res := filterFields("map", fields)
	assert.Len(t, res, 2)

	expected := []Field{
		{Name: "ID", Type: "int", Tags: []Tag{
			{Name: "map", Value: "id"},
		}},
		{Name: "Name", Type: "string", Tags: []Tag{
			{Name: "map", Value: "name"},
		}},
	}
	assert.Equal(t, expected, res)
}
