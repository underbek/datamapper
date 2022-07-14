package generator

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"html/template"
	"os"

	"github.com/underbek/datamapper/models"
	"github.com/underbek/datamapper/utils"
	"golang.org/x/exp/maps"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

func (g *Generator) CreateConvertor(from, to models.Struct, dest string) error {
	content, err := g.generateConvertor(from, to, dest)
	if err != nil {
		return err
	}

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

func (g *Generator) generateConvertor(from, to models.Struct, dest string) ([]byte, error) {
	temp, err := template.ParseFS(templates, convertorFilePath)
	if err != nil {
		return nil, err
	}

	pkg, err := utils.LoadPackage(dest)
	if err != nil {
		return nil, err
	}

	res, err := g.createModelsPair(from, to, pkg.PkgPath)
	if err != nil {
		return nil, err
	}

	res.imports = append(res.imports, from.PackagePath, to.PackagePath)

	convertorName := "Convert"
	fromName := from.Name
	toName := to.Name
	if from.PackagePath != pkg.PkgPath {
		fromName = fmt.Sprintf("%s.%s", from.PackageName, from.Name)
		convertorName += cases.Title(language.Und, cases.NoLower).String(from.PackageName)
	}
	convertorName += from.Name
	convertorName += "To"
	if to.PackagePath != pkg.PkgPath {
		toName = fmt.Sprintf("%s.%s", to.PackageName, to.Name)
		convertorName += cases.Title(language.Und, cases.NoLower).String(to.PackageName)
	}
	convertorName += to.Name

	data := map[string]any{
		"packageName":   pkg.Name,
		"fromName":      fromName,
		"toName":        toName,
		"convertorName": convertorName,
		"fields":        res.fields,
		"imports":       filterImports(pkg.PkgPath, res.imports),
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

func (g *Generator) createModelsPair(from, to models.Struct, pkgPath string) (result, error) {
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

		conversion, pack, err := g.getConversionFunction(fromField.Type, toField.Type, fromField.Name, pkgPath)
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

func (g *Generator) getConversionFunction(fromType, toType models.Type, fromFieldName, pkgPath string,
) (ConvertorType, ImportType, error) {

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

	if cf.PackagePath == pkgPath {
		conversion := fmt.Sprintf("%s%s(from.%s)", cf.Name, typeParams, fromFieldName)
		return conversion, cf.PackagePath, nil
	}

	conversion := fmt.Sprintf("%s.%s%s(from.%s)", cf.PackageName, cf.Name, typeParams, fromFieldName)

	return conversion, cf.PackagePath, nil
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

func filterImports(currentPkgPath string, imports []string) []string {
	res := make([]string, 0, len(imports))
	for _, imp := range imports {
		if imp != currentPkgPath {
			res = append(res, imp)
		}
	}

	return res
}
