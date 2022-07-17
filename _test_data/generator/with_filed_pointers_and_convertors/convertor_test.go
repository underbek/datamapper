package with_filed_pointers_and_convertors

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_Convertor(t *testing.T) {
	idStr := "123"
	ageFloat := 12.58
	childrenInt := 5
	childrenDec := decimal.NewFromInt(5)

	from := From{
		ID:       123,
		Age:      &ageFloat,
		Children: &childrenInt,
	}

	expected := To{
		UUID:     &idStr,
		Age:      decimal.NewFromFloat(12.58),
		Children: &childrenDec,
	}

	actual, err := ConvertFromToTo(from)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
