// Code generated by datamapper.
// https://github.com/underbek/datamapper

// Package with_from_and_to_pointers_without_errors is a generated datamapper package.
package with_from_and_to_pointers_without_errors

import "github.com/underbek/datamapper/converts"

// ConvertFromToTo convert *From by tag map to *To by tag map
func ConvertFromToTo(from *From) *To {
	if from == nil {
		return nil
	}

	return &To{
		Name: from.Name,
		Age:  converts.ConvertOrderedToOrdered[int, uint](from.Age),
	}
}
