package convertors

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"golang.org/x/exp/constraints"
)

func CustomUUIDToInteger[T constraints.Integer](from uuid.UUID) T {
	return T(from.ID())
}

func CustomIntegerToUUID[T constraints.Integer](from T) uuid.UUID {
	return uuid.New()
}

func BrokenConvertor(from decimal.Decimal) int {
	return ""
}
