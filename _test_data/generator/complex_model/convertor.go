package complex_model

import "github.com/underbek/datamapper/converts"

func ConvertFromToTo(from From) To {
	return To{
		UUID: converts.ConvertNumericToString(from.ID),
		Name: from.Name,
		Age:  converts.ConvertFloatToDecimal(from.Age),
	}
}
