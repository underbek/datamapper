package complex_model

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_Convertor(t *testing.T) {
	from := From{
		ID:   123,
		Name: "test_name",
		Age:  12.58,
	}

	expected := To{
		UUID: "123",
		Name: "test_name",
		Age:  decimal.NewFromFloat(12.58),
	}

	actual := ConvertFromToTo(from)

	assert.Equal(t, expected, actual)
}
