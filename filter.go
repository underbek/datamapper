package main

func findTag(tagName string, tags []Tag) (Tag, bool) {
	for _, tag := range tags {
		if tag.Name == tagName {
			return tag, true
		}
	}

	return Tag{}, false
}

func filterFields(tagName string, fields []Field) []Field {
	var res []Field
	for _, field := range fields {
		tag, ok := findTag(tagName, field.Tags)
		if !ok {
			continue
		}

		field.Tags = []Tag{tag}
		res = append(res, field)
	}

	return res
}
