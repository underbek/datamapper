package domain

import "github.com/shopspring/decimal"

type User struct {
	ID         int             `map:"id"`
	Name       string          `map:"name"`
	Age        decimal.Decimal `map:"age"`
	ChildCount int             `map:"children"`
	Empty      string
}
