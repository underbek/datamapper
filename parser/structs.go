package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/underbek/datamapper/models"
)

func ParseStructs(source string) (map[string]models.Struct, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, source, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	structs := make(map[string]models.Struct)

	for _, f := range node.Decls {
		genD, ok := f.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range genD.Specs {
			currType, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			currStruct, ok := currType.Type.(*ast.StructType)
			if !ok {
				continue
			}

			if len(currStruct.Fields.List) == 0 {
				continue
			}

			fields := make([]models.Field, 0, len(currStruct.Fields.List))
			for _, field := range currStruct.Fields.List {
				fieldType, ok := field.Type.(*ast.Ident)
				if !ok {
					continue
				}

				if len(field.Names) == 0 {
					continue
				}

				// TODO: add package into field type
				fields = append(fields, models.Field{
					Name: field.Names[0].Name,
					Type: models.Type{Name: fieldType.Name},
					Tags: parseTag(field.Tag),
				})
			}

			structs[currType.Name.Name] = models.Struct{
				Name:   currType.Name.Name,
				Fields: fields,
			}
		}
	}
	return structs, nil
}

func parseTag(tag *ast.BasicLit) []models.Tag {
	if tag == nil {
		return nil
	}

	value := strings.Trim(tag.Value, "`")
	textTags := strings.Split(value, " ")

	tags := make([]models.Tag, 0, len(textTags))
	for _, textTag := range textTags {
		sepIndex := strings.Index(textTag, ":")
		if sepIndex == -1 {
			continue
		}

		tags = append(tags, models.Tag{
			Name:  textTag[:sepIndex],
			Value: strings.Trim(textTag[sepIndex+1:], "\""),
		})
	}

	return tags
}
