package other_package_model

import (
	"github.com/underbek/datamapper/converts"

	"github.com/underbek/datamapper/_test_data/generator/other_package_model/other"
)

func ConvertOtherFromToTo(from other.From) To {
	return To{
		UUID: converts.ConvertNumericToString(from.ID),
		Name: from.Name,
	}
}