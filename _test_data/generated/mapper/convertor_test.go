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
	model := transport.User{
		UUID:       uuid.New(),
		Name:       "test_name",
		Age:        "12.58",
		ChildCount: "2",
	}

	expected := domain.User{
		ID:         int(model.UUID.ID()),
		Name:       "test_name",
		Age:        decimal.NewFromFloat(12.58),
		ChildCount: 2,
	}

	actual := ConvertTransportUserToDomainUser(model)

	assert.Equal(t, expected, actual)
}
