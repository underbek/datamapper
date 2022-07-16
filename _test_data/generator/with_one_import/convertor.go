package with_one_import

import "github.com/underbek/datamapper/converts"

func ConvertFromToTo(from From) To {
	return To{
		UUID: converts.ConvertNumericToString(from.ID),
		Name: from.Name,
	}
}
