package with_filed_pointers_and_errors

import (
	"github.com/underbek/datamapper/converts"

	"fmt"
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

	if from.Children == nil {
		return To{}, fmt.Errorf("cannot convert From.Children -> To.Children, field is nil")
	}

	fromChildren, err := converts.ConvertStringToDecimal(*from.Children)
	if err != nil {
		return To{}, err
	}

	return To{
		UUID:     &fromID,
		Age:      fromAge,
		Children: &fromChildren,
	}, nil
}
