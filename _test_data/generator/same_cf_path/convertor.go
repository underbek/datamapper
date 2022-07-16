package same_cf_path

func ConvertFromToTo(from From) To {
	return To{
		UUID: ConvertNumericToString(from.ID),
		Name: from.Name,
		Age:  ConvertFloatToDecimal(from.Age),
	}
}
