package generator

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"html/template"
	"os"

	"github.com/underbek/datamapper/models"
	"golang.org/x/exp/maps"
)

type ConvertorType = string
type ImportType = string

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
	functions models.Functions
}

func New(functions models.Functions) *Generator {
	return &Generator{
		functions: functions,
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

		conversion, pack, err := g.getConvertorFunctions(fromField.Type, toField.Type, fromField.Name)
		if err != nil {
			return result{}, err
		}

		if pack != "" {
			imports[pack] = struct{}{}
		}

		fields = append(fields, FieldsPair{
			FromName:   fromField.Name,
			FromType:   fromField.Type.Name,
			ToName:     toField.Name,
			ToType:     toField.Type.Name,
			Conversion: conversion,
		})
	}

	return result{
		fields:  fields,
		imports: maps.Keys(imports),
	}, nil
}

func (g *Generator) getConvertorFunctions(fromType, toType models.Type, fromFieldName string) (ConvertorType, ImportType, error) {
	// TODO: check package
	if fromType.Name == toType.Name {
		return fmt.Sprintf("from.%s", fromFieldName), "", nil
	}

	cf, ok := g.functions[models.ConversionFunctionKey{
		FromType: fromType,
		ToType:   toType,
	}]

	if !ok {
		return "", "", fmt.Errorf(
			"not found convertor function for types %s -> %s by %s field",
			fromType,
			toType,
			fromFieldName,
		)
	}

	typeParams := getTypeParams(cf, fromType, toType)
	conversion := fmt.Sprintf("%s.%s%s(from.%s)", cf.PackageName, cf.Name, typeParams, fromFieldName)

	return conversion, cf.Import, nil
}

func getTypeParams(cf models.ConversionFunction, fromType, toType models.Type) string {
	// TODO: add packages to type params
	switch cf.TypeParam {
	case models.ToTypeParam:
		return fmt.Sprintf("[%s]", toType.Name)
	case models.FromToTypeParam:
		return fmt.Sprintf("[%s,%s]", fromType.Name, toType.Name)
	default:
		return ""
	}
}
