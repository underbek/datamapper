package with_error

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_Convertor(t *testing.T) {
	from := From{
		UUID: "123",
		Name: "test_name",
		Age:  12,
	}

	expected := To{
		ID:   decimal.NewFromInt(123),
		Name: "test_name",
		Age:  decimal.NewFromFloat(12),
	}

	actual, err := ConvertFromToTo(from)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
