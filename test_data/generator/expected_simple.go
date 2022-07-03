package generator

func convertModelToDAO(from Model) DAO {
	return DAO{
		Name: from.Name,
	}
}
