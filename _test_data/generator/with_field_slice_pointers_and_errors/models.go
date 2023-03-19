package with_field_slice_pointers_and_errors

import "github.com/shopspring/decimal"

type From struct {
	IDs     []decimal.Decimal  `map:"ids"`
	Ages    []*decimal.Decimal `map:"ages"`
	Counts  []*int             `map:"counts"`
	Origins []*float64         `map:"origs"`
}

type To struct {
	UUIDs   []*int             `map:"ids"`
	Ages    []int              `map:"ages"`
	Counts  []*decimal.Decimal `map:"counts"`
	Origins []*decimal.Decimal `map:"origs"`
}
