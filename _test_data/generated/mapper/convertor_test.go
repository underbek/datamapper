package mapper

import (
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/underbek/datamapper/_test_data/mapper/domain"
	"github.com/underbek/datamapper/_test_data/mapper/transport"
)

func Test_ConvertTransportUserToDomainUser(t *testing.T) {
	childrenString := "2"
	childrenInt := 2

	model := transport.User{
		UUID:       uuid.New(),
		Name:       "test_name",
		Age:        "12.58",
		ChildCount: &childrenString,
	}

	expected := domain.User{
		ID:         int(model.UUID.ID()),
		Name:       "test_name",
		Age:        decimal.NewFromFloat(12.58),
		ChildCount: &childrenInt,
	}

	actual, err := ConvertTransportUserToDomainUser(model)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
