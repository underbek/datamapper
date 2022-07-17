package with_filed_pointers_and_convertors

import (
	"github.com/underbek/datamapper/converts"

	"fmt"
)

func ConvertFromToTo(from From) (To, error) {
	fromID := converts.ConvertNumericToString(from.ID)

	if from.Age == nil {
		return To{}, fmt.Errorf("cannot convert From.Age -> To.Age, field is nil")
	}

	if from.Children == nil {
		return To{}, fmt.Errorf("cannot convert From.Children -> To.Children, field is nil")
	}

	fromChildren := converts.ConvertIntegerToDecimal(*from.Children)

	return To{
		UUID:     &fromID,
		Age:      converts.ConvertFloatToDecimal(*from.Age),
		Children: &fromChildren,
	}, nil
}
