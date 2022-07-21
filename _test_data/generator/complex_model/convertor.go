// Code generated by datamapper.
// https://github.com/underbek/datamapper

// Package complex_model is a generated datamapper package.
package complex_model

import "github.com/underbek/datamapper/converts"

// ConvertFromToTo convert From by tag map to To by tag map
func ConvertFromToTo(from From) To {
	return To{
		UUID: converts.ConvertNumericToString(from.ID),
		Name: from.Name,
		Age:  converts.ConvertFloatToDecimal(from.Age),
	}
}
