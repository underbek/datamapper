package without_imports

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Convertor(t *testing.T) {
	from := From{
		Name: "test_name",
	}

	expected := To{
		Name: "test_name",
	}

	actual := ConvertFromToTo(from)

	assert.Equal(t, expected, actual)
}
