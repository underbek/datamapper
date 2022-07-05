package main

import "github.com/underbek/datamapper/models"

func findTag(tagName string, tags []models.Tag) (models.Tag, bool) {
	for _, tag := range tags {
		if tag.Name == tagName {
			return tag, true
		}
	}

	return models.Tag{}, false
}

func filterFields(tagName string, fields []models.Field) []models.Field {
	var res []models.Field
	for _, field := range fields {
		tag, ok := findTag(tagName, field.Tags)
		if !ok {
			continue
		}

		field.Tags = []models.Tag{tag}
		res = append(res, field)
	}

	return res
}
