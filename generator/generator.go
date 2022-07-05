package generator

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"html/template"
	"os"

	"github.com/underbek/datamapper/converts"
	"github.com/underbek/datamapper/models"
	"golang.org/x/exp/maps"
)

type FieldsPair struct {
	FromName   string
	FromType   string
	ToName     string
	ToType     string
	Conversion string
}

type result struct {
	fields  []FieldsPair
	imports []string
}

const convertorFilePath = "templates/convertor.temp"

//go:embed templates
var templates embed.FS

type Generator struct {
	convertorFactory *converts.Factory
}

func New(convertorFactory *converts.Factory) *Generator {
	return &Generator{
		convertorFactory: convertorFactory,
	}
}

func (g *Generator) CreateConvertor(from, to models.Struct, dest string, packageName string) error {
	content, err := g.generateConvertor(from, to, packageName)

	file, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	_, err = file.Write(content)
	if err != nil {
		return err
	}

	return nil
}

func (g *Generator) generateConvertor(from, to models.Struct, packageName string) ([]byte, error) {
	temp, err := template.ParseFS(templates, convertorFilePath)
	if err != nil {
		return nil, err
	}

	res, err := g.createModelsPair(from, to)
	if err != nil {
		return nil, err
	}

	data := map[string]any{
		"packageName": packageName,
		"fromName":    from.Name,
		"toName":      to.Name,
		"fields":      res.fields,
		"imports":     res.imports,
	}

	buf := bytes.Buffer{}
	err = temp.Execute(&buf, data)
	if err != nil {
		return nil, err
	}

	content, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (g *Generator) createModelsPair(from, to models.Struct) (result, error) {
	var fields []FieldsPair
	imports := make(map[string]struct{})

	fromFields := make(map[string]models.Field)
	for _, field := range from.Fields {
		fromFields[field.Tags[0].Value] = field
	}

	for _, toField := range to.Fields {
		fromField, ok := fromFields[toField.Tags[0].Value]
		if !ok {
			//TODO: warning or error politics
			continue
		}

		conversion, pack := g.convertorFactory.GetConvertorFunctions(fromField.Type, toField.Type, fromField.Name)
		if conversion == "" {
			return result{}, fmt.Errorf(
				"not found convertor function for types %s -> %s by %s field",
				fromField.Type,
				toField.Type,
				fromField.Name,
			)
		}

		if pack != "" {
			imports[pack] = struct{}{}
		}

		fields = append(fields, FieldsPair{
			FromName:   fromField.Name,
			FromType:   fromField.Type,
			ToName:     toField.Name,
			ToType:     toField.Type,
			Conversion: conversion,
		})
	}

	return result{
		fields:  fields,
		imports: maps.Keys(imports),
	}, nil
}
