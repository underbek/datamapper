package generator

import (
	"fmt"

	"github.com/underbek/datamapper/models"
)

type FindFieldsPairError struct {
	From          models.Type
	To            models.Type
	fromFieldName string
}

func NewFindFieldsPairError(from, to models.Type, fromFieldName string) *FindFieldsPairError {
	return &FindFieldsPairError{
		From:          from,
		To:            to,
		fromFieldName: fromFieldName,
	}
}

func (e *FindFieldsPairError) Error() string {
	return fmt.Sprintf(
		"not found convertor function for types %s -> %s by %s field",
		e.From.Name,
		e.To.Name,
		e.fromFieldName,
	)
}
