package with_filed_pointers_and_errors

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/underbek/datamapper/converts"
)

func ConvertFromToTo(from From) (To, error) {
	fromID, err := converts.ConvertStringToSigned[int](from.ID)
	if err != nil {
		return To{}, err
	}

	if from.Age == nil {
		return To{}, fmt.Errorf("cannot convert From.Age -> To.Age, field is nil")
	}

	fromAge, err := converts.ConvertStringToDecimal(*from.Age)
	if err != nil {
		return To{}, err
	}

	var fromChildren *decimal.Decimal
	if from.Children != nil {
		res, err := converts.ConvertStringToDecimal(*from.Children)
		if err != nil {
			return To{}, err
		}

		fromChildren = &res
	}

	return To{
		UUID:     &fromID,
		Age:      fromAge,
		Children: fromChildren,
	}, nil
}
