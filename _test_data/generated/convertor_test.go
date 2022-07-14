package generator

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/underbek/datamapper/_test_data/other"
)

func Test_ConvertModelToDao(t *testing.T) {
	model := Model{
		ID:    123,
		Name:  "test_name",
		Empty: "empty",
		Age:   12.58,
	}

	expected := DAO{
		UUID: "123",
		Name: "test_name",
		Age:  12,
	}

	actual := ConvertModelToDAO(model)

	assert.Equal(t, expected, actual)
}

func Test_ConvertOtherModelToDTO(t *testing.T) {
	model := other.Model{
		ID:   "123",
		Name: "test_name",
		Age:  12.58,
	}

	expected := DTO{
		ID:  decimal.NewFromInt(123),
		Age: decimal.NewFromFloat(12.58),
	}

	actual := ConvertOtherModelToDTO(model)

	assert.Equal(t, expected, actual)
}

func Test_ConvertDTOToOtherModel(t *testing.T) {
	dto := DTO{
		ID:  decimal.NewFromInt(123),
		Age: decimal.NewFromFloat(12.58),
	}

	expected := other.Model{
		ID:   "123",
		Name: "",
		Age:  12.58,
	}

	actual := ConvertDTOToOtherModel(dto)

	assert.Equal(t, expected, actual)
}
