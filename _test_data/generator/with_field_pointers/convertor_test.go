package with_filed_pointers

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_Convertor(t *testing.T) {
	name := "test_name"
	age := decimal.NewFromFloat(12.58)

	from := From{
		ID:   123,
		Name: &name,
		Age:  &age,
	}

	expected := To{
		UUID: &from.ID,
		Name: &name,
		Age:  age,
	}

	actual, err := ConvertFromToTo(from)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
