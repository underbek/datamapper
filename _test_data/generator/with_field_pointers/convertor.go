package with_filed_pointers

import "fmt"

func ConvertFromToTo(from From) (To, error) {
	if from.Age == nil {
		return To{}, fmt.Errorf("cannot convert From.Age -> To.Age, field is nil")
	}

	return To{
		UUID: &from.ID,
		Name: from.Name,
		Age:  *from.Age,
	}, nil
}
