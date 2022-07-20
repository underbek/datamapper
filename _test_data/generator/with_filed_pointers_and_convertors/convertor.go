package with_filed_pointers_and_convertors

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/underbek/datamapper/converts"
)

func ConvertFromToTo(from From) (To, error) {
	fromID := converts.ConvertNumericToString(from.ID)

	if from.Age == nil {
		return To{}, fmt.Errorf("cannot convert From.Age -> To.Age, field is nil")
	}

	var fromChildren *decimal.Decimal
	if from.Children != nil {
		res := converts.ConvertIntegerToDecimal(*from.Children)
		fromChildren = &res
	}

	return To{
		UUID:     &fromID,
		Age:      converts.ConvertFloatToDecimal(*from.Age),
		Children: fromChildren,
	}, nil
}
