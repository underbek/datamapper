// Code generated by datamapper.
// https://github.com/underbek/datamapper

// Package with_field_slice_and_errors is a generated datamapper package.
package with_field_slice_and_errors

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/underbek/datamapper/_test_data/generator/with_field_slice_and_errors/cf"
)

// ConvertFromToTo convert From by tag map to To by tag map
func ConvertFromToTo(from From) (To, error) {
	fromIDs := make([]decimal.Decimal, 0, len(from.IDs))
	for _, item := range from.IDs {
		res, err := cf.ConvertStringToDecimal(item)
		if err != nil {
			return To{}, fmt.Errorf("convert From.IDs -> To.UUIDs failed: %w", err)
		}

		fromIDs = append(fromIDs, res)
	}

	return To{
		UUIDs: fromIDs,
	}, nil
}