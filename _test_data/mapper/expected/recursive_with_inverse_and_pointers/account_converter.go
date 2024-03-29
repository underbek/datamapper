// Code generated by datamapper.
// https://github.com/underbek/datamapper

// Package mapper is a generated datamapper package.
package mapper

import (
	"errors"
	"fmt"

	f "github.com/underbek/datamapper/_test_data/mapper/recursive/from"
	t "github.com/underbek/datamapper/_test_data/mapper/recursive/to"
	"github.com/underbek/datamapper/converts"
)

// ConvertFAccountToTAccount convert *f.Account by tag recursive to t.Account by tag recursive
func ConvertFAccountToTAccount(from *f.Account) (t.Account, error) {
	if from == nil {
		return t.Account{}, errors.New("Account is nil")
	}

	return t.Account{
		ID:     from.ID,
		Amount: converts.ConvertDecimalToString(from.Amount),
	}, nil
}

// ConvertTAccountToFAccount convert t.Account by tag recursive to *f.Account by tag recursive
func ConvertTAccountToFAccount(from t.Account) (*f.Account, error) {
	fromAmount, err := converts.ConvertStringToDecimal(from.Amount)
	if err != nil {
		return nil, fmt.Errorf("convert Account.Amount -> Account.Amount failed: %w", err)
	}

	return &f.Account{
		ID:     from.ID,
		Amount: fromAmount,
	}, nil
}
