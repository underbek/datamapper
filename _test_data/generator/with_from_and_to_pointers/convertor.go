// Code generated by datamapper.
// https://github.com/underbek/datamapper

// Package with_from_and_to_pointers is a generated datamapper package.
package with_from_and_to_pointers

import "github.com/underbek/datamapper/converts"

// ConvertFromToTo convert *From by tag map to *To by tag map
func ConvertFromToTo(from *From) (*To, error) {
	if from == nil {
		return nil, nil
	}

	fromUUID, err := converts.ConvertStringToDecimal(from.UUID)
	if err != nil {
		return nil, err
	}

	fromAge, err := converts.ConvertStringToSigned[int8](from.Age)
	if err != nil {
		return nil, err
	}

	return &To{
		ID:   fromUUID,
		Name: from.Name,
		Age:  fromAge,
	}, nil
}
