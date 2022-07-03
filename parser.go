package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func parseStructs(source string) (map[string]Struct, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, source, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	structs := make(map[string]Struct)

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

			fields := make([]Field, 0, len(currStruct.Fields.List))
			for _, field := range currStruct.Fields.List {
				fieldType, ok := field.Type.(*ast.Ident)
				if !ok {
					continue
				}

				if len(field.Names) == 0 {
					continue
				}

				fields = append(fields, Field{
					Name: field.Names[0].Name,
					Type: fieldType.Name,
					Tags: parseTag(field.Tag),
				})
			}

			structs[currType.Name.Name] = Struct{
				Name:   currType.Name.Name,
				Fields: fields,
			}
		}
	}
	return structs, nil
}

func parseTag(tag *ast.BasicLit) []Tag {
	if tag == nil {
		return nil
	}

	value := strings.Trim(tag.Value, "`")
	textTags := strings.Split(value, " ")

	tags := make([]Tag, 0, len(textTags))
	for _, textTag := range textTags {
		sepIndex := strings.Index(textTag, ":")
		if sepIndex == -1 {
			continue
		}

		tags = append(tags, Tag{
			Name:  textTag[:sepIndex],
			Value: strings.Trim(textTag[sepIndex+1:], "\""),
		})
	}

	return tags
}
