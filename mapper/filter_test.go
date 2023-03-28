package mapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/underbek/datamapper/models"
)

func Test_Filter(t *testing.T) {
	fields := []models.Field{
		{Name: "ID", Type: models.Type{Name: "int"}, Tags: []models.Tag{
			{Name: "json", Value: "id"},
			{Name: "map", Value: "id"},
		}},
		{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{
			{Name: "json", Value: "name"},
			{Name: "map", Value: "name"},
		}},
		{Name: "Empty", Type: models.Type{Name: "string"}},
	}

	res := filterFields("map", models.NewFields(fields))
	assert.Equal(t, res.Len(), 2)

	expected := models.NewFields([]models.Field{
		{Name: "ID", Type: models.Type{Name: "int"}, Tags: []models.Tag{
			{Name: "map", Value: "id"},
		}},
		{Name: "Name", Type: models.Type{Name: "string"}, Tags: []models.Tag{
			{Name: "map", Value: "name"},
		}},
	})
	assert.Equal(t, expected, res)
}
