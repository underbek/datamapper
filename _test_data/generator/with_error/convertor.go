package with_error

import "github.com/underbek/datamapper/converts"

func ConvertFromToTo(from From) (To, error) {
	fromUUID, err := converts.ConvertStringToDecimal(from.UUID)
	if err != nil {
		return To{}, err
	}

	return To{
		ID:   fromUUID,
		Name: from.Name,
		Age:  converts.ConvertIntegerToDecimal(from.Age),
	}, nil
}
