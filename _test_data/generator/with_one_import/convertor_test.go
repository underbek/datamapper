package with_one_import

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Convertor(t *testing.T) {
	from := From{
		ID:   123,
		Name: "test_name",
	}

	expected := To{
		UUID: "123",
		Name: "test_name",
	}

	actual := ConvertFromToTo(from)

	assert.Equal(t, expected, actual)
}
