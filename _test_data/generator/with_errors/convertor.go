package with_errors

import "github.com/underbek/datamapper/converts"

func ConvertFromToTo(from From) (To, error) {
	fromUUID, err := converts.ConvertStringToDecimal(from.UUID)
	if err != nil {
		return To{}, err
	}

	fromAge, err := converts.ConvertStringToSigned[int8](from.Age)
	if err != nil {
		return To{}, err
	}

	return To{
		ID:   fromUUID,
		Name: from.Name,
		Age:  fromAge,
	}, nil
}
