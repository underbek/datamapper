package other_package_model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/underbek/datamapper/_test_data/generator/other_package_model/other"
)

func Test_Convertor(t *testing.T) {
	from := other.From{
		ID:   123,
		Name: "test_name",
	}

	expected := To{
		UUID: "123",
		Name: "test_name",
	}

	actual := ConvertOtherFromToTo(from)

	assert.Equal(t, expected, actual)
}
