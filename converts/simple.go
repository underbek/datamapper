package converts

import (
	"fmt"
	"strconv"

	"golang.org/x/exp/constraints"
)

const base = 10

func ConvertNumericToString[T constraints.Integer | constraints.Float](from T) string {
	return fmt.Sprint(from)
}

func ConvertComplexToString[T constraints.Complex](from T) string {
	return fmt.Sprint(from)
}

func ConvertOrderedToOrdered[T, V constraints.Integer | constraints.Float](from T) V {
	return V(from)
}

func ConvertStringToSigned[T constraints.Signed](from string) (T, error) {
	res, err := strconv.ParseInt(from, base, 0)
	if err != nil {
		return 0, err
	}
	return T(res), nil
}
