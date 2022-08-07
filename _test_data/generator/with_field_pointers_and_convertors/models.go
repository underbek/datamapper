package with_field_pointers_and_convertors

import "github.com/shopspring/decimal"

type From struct {
	ID       int      `map:"id"`
	Age      *float64 `map:"age"`
	Children *int     `map:"children"`
}

type To struct {
	UUID     *string          `map:"id"`
	Age      decimal.Decimal  `map:"age"`
	Children *decimal.Decimal `map:"children"`
}
