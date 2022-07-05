package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/underbek/datamapper/models"
)

func Test_Filter(t *testing.T) {
	fields := []models.Field{
		{Name: "ID", Type: "int", Tags: []models.Tag{
			{Name: "json", Value: "id"},
			{Name: "map", Value: "id"},
		}},
		{Name: "Name", Type: "string", Tags: []models.Tag{
			{Name: "json", Value: "name"},
			{Name: "map", Value: "name"},
		}},
		{Name: "Empty", Type: "string"},
	}

	res := filterFields("map", fields)
	assert.Len(t, res, 2)

	expected := []models.Field{
		{Name: "ID", Type: "int", Tags: []models.Tag{
			{Name: "map", Value: "id"},
		}},
		{Name: "Name", Type: "string", Tags: []models.Tag{
			{Name: "map", Value: "name"},
		}},
	}
	assert.Equal(t, expected, res)
}
