package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Convert(t *testing.T) {
	model := Model{
		ID:    123,
		Name:  "test_name",
		Empty: "empty",
	}

	expected := DAO{
		UUID: "123",
		Name: "test_name",
	}

	actual := convertModelToDAO(model)

	assert.Equal(t, expected, actual)
}
