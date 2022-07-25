// Code generated by datamapper.
// https://github.com/underbek/datamapper

// Package cf_with_pointers is a generated datamapper package.
package cf_with_pointers

import "github.com/underbek/datamapper/_test_data/generator/cf_with_pointers/cf"

// ConvertFromToTo convert From by tag map to To by tag map
func ConvertFromToTo(from From) To {
	return To{
		UUID:  cf.ConvertIntToDecimalPtr(from.ID),
		Age:   cf.ConvertIntPtrToDecimal(from.Age),
		Count: cf.ConvertIntPtrToDecimalPtr(from.Count),
		Orig:  cf.ConvertIntToDecimal(from.Orig),
	}
}