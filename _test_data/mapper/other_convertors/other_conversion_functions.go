package other_convertors

import (
	"github.com/google/uuid"
	"golang.org/x/exp/constraints"
)

func CustomIntegerToUUID[T constraints.Integer](from T) uuid.UUID {
	return uuid.New()
}
