package with_field_pointers_without_error

import "github.com/shopspring/decimal"

type From struct {
	ID   int      `map:"id"`
	Name *string  `map:"name"`
	Age  *float64 `map:"age"`
}

type To struct {
	UUID *int             `map:"id"`
	Name *string          `map:"name"`
	Age  *decimal.Decimal `map:"age"`
}
