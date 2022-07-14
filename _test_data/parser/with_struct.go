package parser

import "github.com/underbek/datamapper/_test_data/other"

func ConvertCurrentModelToOther(from Model) other.Model {
	return other.Model{
		ID: from.ID,
	}
}

func ConvertOtherModelToCurrent(from other.Model) Model {
	return Model{
		ID: from.ID,
	}
}
