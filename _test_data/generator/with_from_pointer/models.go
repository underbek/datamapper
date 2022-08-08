package with_from_pointer

import "github.com/shopspring/decimal"

type From struct {
	UUID string `map:"id"`
	Name string `map:"name"`
	Age  string `map:"age"`
}

type To struct {
	ID   decimal.Decimal `map:"id"`
	Name string          `map:"name"`
	Age  int8            `map:"age"`
}
