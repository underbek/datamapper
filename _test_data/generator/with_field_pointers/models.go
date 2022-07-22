package with_filed_pointers

import "github.com/shopspring/decimal"

type From struct {
	ID   int              `map:"id"`
	Name *string          `map:"name"`
	Age  *decimal.Decimal `map:"age"`
}

type To struct {
	UUID *int            `map:"id"`
	Name *string         `map:"name"`
	Age  decimal.Decimal `map:"age"`
}
