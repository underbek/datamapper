package cf_with_slice_and_errors

import "github.com/shopspring/decimal"

type From struct {
	IDs []string `map:"ids"`
}

type To struct {
	UUIDs []decimal.Decimal `map:"ids"`
}
