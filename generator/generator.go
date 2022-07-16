package generator

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"html/template"
	"os"
	"strings"

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
	WithError  bool
}

type result struct {
	fields  []FieldsPair
	imports []string
}

const convertorFilePath = "templates/convertor.temp"

//go:embed templates
var templates embed.FS

func CreateConvertor(from, to models.Struct, dest string, functions models.Functions) error {
	content, err := generateConvertor(from, to, dest, functions)
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

func generateConvertor(from, to models.Struct, dest string, functions models.Functions) ([]byte, error) {
	temp, err := template.ParseFS(templates, convertorFilePath)
	if err != nil {
		return nil, err
	}

	pkg, err := utils.LoadPackage(dest)
	if err != nil {
		return nil, err
	}

	res, err := createModelsPair(from, to, pkg.PkgPath, functions)
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

	pkgName := pkg.Name
	if pkgName == "" {
		names := strings.Split(pkg.PkgPath, "/")
		if len(names) == 0 {
			return nil, fmt.Errorf("incorrect parsed package path from destination %s", dest)
		}
		pkgName = names[len(names)-1]
	}

	data := map[string]any{
		"packageName":   pkgName,
		"fromName":      fromName,
		"toName":        toName,
		"convertorName": convertorName,
		"fields":        res.fields,
		"imports":       filterImports(pkg.PkgPath, res.imports),
		"withError":     isReturnError(res.fields),
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

func createModelsPair(from, to models.Struct, pkgPath string, functions models.Functions) (result, error) {
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

		conversion, pack, withError, err := getConversionFunction(fromField.Type, toField.Type, fromField.Name, pkgPath, functions)
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
			WithError:  withError,
		})
	}

	return result{
		fields:  fields,
		imports: maps.Keys(imports),
	}, nil
}

func getConversionFunction(fromType, toType models.Type, fromFieldName, pkgPath string, functions models.Functions,
) (ConvertorType, ImportType, bool, error) {

	// TODO: check package
	if fromType.Name == toType.Name {
		return fmt.Sprintf("from.%s", fromFieldName), "", false, nil
	}

	cf, ok := functions[models.ConversionFunctionKey{
		FromType: fromType,
		ToType:   toType,
	}]

	if !ok {
		return "", "", false, fmt.Errorf(
			"not found convertor function for types %s -> %s by %s field",
			fromType,
			toType,
			fromFieldName,
		)
	}

	typeParams := getTypeParams(cf, fromType, toType)

	if cf.PackagePath == pkgPath {
		conversion := fmt.Sprintf("%s%s(from.%s)", cf.Name, typeParams, fromFieldName)
		return conversion, cf.PackagePath, cf.WithError, nil
	}

	conversion := fmt.Sprintf("%s.%s%s(from.%s)", cf.PackageName, cf.Name, typeParams, fromFieldName)

	return conversion, cf.PackagePath, cf.WithError, nil
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

func isReturnError(fields []FieldsPair) bool {
	for _, field := range fields {
		if field.WithError {
			return true
		}
	}

	return false
}
