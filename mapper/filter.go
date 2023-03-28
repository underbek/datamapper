package mapper

import (
	"fmt"

	"github.com/underbek/datamapper/logger"
	"github.com/underbek/datamapper/models"
	"github.com/underbek/datamapper/parser"
)

const dash = "-"

func TransformAndFilterFields(lg logger.Logger, tagName string, structure models.Struct, head *models.Field) (models.Struct, error) {
	structure.Fields = filterFields(tagName, structure.Fields)
	newStruct := models.Struct{
		Type: structure.Type,
	}

	err := structure.Fields.Each(func(field *models.Field) error {
		field.Head = head
		if field.Tags[0].Value != dash {
			field.CurrentStruct = &newStruct
			newStruct.Fields.Add(*field)
			return nil
		}

		if field.Type.Kind != models.StructType {
			lg.Infof("skip base type %s from %s", field.Type.Name, fullFieldNameForSkipComment(*field))
			return nil
		}

		newField, err := transformAndFilterField(lg, *field)
		if err != nil {
			lg.Errorf("transform field error: %s", err)
			return err
		}

		newField.CurrentStruct = &newStruct
		newStruct.Fields.Add(newField)
		return nil
	})
	if err != nil {
		return models.Struct{}, fmt.Errorf("transform and filter fields error: %w", err)
	}

	return newStruct, nil
}

func findTag(tagName string, tags []models.Tag) (models.Tag, bool) {
	for _, tag := range tags {
		if tag.Name == tagName {
			return tag, true
		}
	}

	return models.Tag{}, false
}

func filterFields(tagName string, fields models.Fields) models.Fields {
	return fields.Filter(func(field *models.Field) bool {
		tag, ok := findTag(tagName, field.Tags)
		if !ok {
			return false
		}

		field.Tags = []models.Tag{tag}
		return true
	})
}

func transformAndFilterField(lg logger.Logger, field models.Field) (models.Field, error) {
	structs, err := parser.ParseModelsByPackage(lg, field.Type.Package.Path)
	if err != nil {
		return models.Field{}, fmt.Errorf("parse dash field error: %w", err)
	}

	fieldStruct, ok := structs[field.Type.Name]
	if !ok {
		return models.Field{}, fmt.Errorf("struct %s not found in package %s", field.Type.Name, field.Type.Package.Path)
	}

	newStruct, err := TransformAndFilterFields(lg, field.Tags[0].Name, fieldStruct, &field)
	if err != nil {
		return models.Field{}, fmt.Errorf("transform dash field error: %w", err)
	}
	field.SkippedStruct = &fieldStruct
	field.SkippedStruct.Fields = newStruct.Fields

	return field, nil
}

func fullFieldNameForSkipComment(field models.Field) string {
	res := field.Name

	for field.Head != nil {
		res = field.Head.Name + "." + res
		field = *field.Head
	}

	return res
}
