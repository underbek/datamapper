// Code generated by datamapper.
// https://github.com/underbek/datamapper

// Package mapper is a generated datamapper package.
package mapper

import (
	"fmt"

	f "github.com/underbek/datamapper/_test_data/mapper/recursive/from"
	t "github.com/underbek/datamapper/_test_data/mapper/recursive/to"
	"github.com/underbek/datamapper/converts"
)

// ConvertFAccountToTAccount convert f.Account by tag recursive to t.Account by tag recursive
func ConvertFAccountToTAccount(from f.Account) t.Account {
	return t.Account{
		ID:     from.ID,
		Amount: converts.ConvertDecimalToString(from.Amount),
	}
}

// ConvertTAccountToFAccount convert t.Account by tag recursive to f.Account by tag recursive
func ConvertTAccountToFAccount(from t.Account) (f.Account, error) {
	fromAmount, err := converts.ConvertStringToDecimal(from.Amount)
	if err != nil {
		return f.Account{}, fmt.Errorf("convert Account.Amount -> Account.Amount failed: %w", err)
	}

	return f.Account{
		ID:     from.ID,
		Amount: fromAmount,
	}, nil
}