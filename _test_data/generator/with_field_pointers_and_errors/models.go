package with_field_pointers_and_errors

import "github.com/shopspring/decimal"

type From struct {
	ID       string  `map:"id"`
	Age      *string `map:"age"`
	Children *string `map:"children"`
}

type To struct {
	UUID     *int             `map:"id"`
	Age      decimal.Decimal  `map:"age"`
	Children *decimal.Decimal `map:"children"`
}
