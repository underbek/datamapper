package with_aliases

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/underbek/datamapper/_test_data/generator/with_aliases/from"
	"github.com/underbek/datamapper/_test_data/generator/with_aliases/to"
)

func Test_Convertor(t *testing.T) {
	fromModel := from.From{
		ID:   123,
		Name: "test_name",
	}

	expected := to.To{
		UUID: "123",
		Name: "test_name",
	}

	actual := ConvertFromaliasFromToToaliasTo(fromModel)

	assert.Equal(t, expected, actual)
}
