package with_filed_pointers_and_errors

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_Convertor(t *testing.T) {
	idInt := 123
	ageString := "12.58"
	childrenString := "5"
	childrenDec := decimal.NewFromInt(5)

	from := From{
		ID:       "123",
		Age:      &ageString,
		Children: &childrenString,
	}

	expected := To{
		UUID:     &idInt,
		Age:      decimal.NewFromFloat(12.58),
		Children: &childrenDec,
	}

	actual, err := ConvertFromToTo(from)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
