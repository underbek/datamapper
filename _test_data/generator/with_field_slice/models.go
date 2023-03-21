package with_field_slice

import "github.com/shopspring/decimal"

type From struct {
	IDs  []int    `map:"ids"`
	Ages []string `map:"ages"`
}

type To struct {
	UUIDs []decimal.Decimal `map:"ids"`
	Ages  []decimal.Decimal `map:"ages"`
}
