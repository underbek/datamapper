package generator

import (
	"bytes"
	"embed"
	"text/template"

	"github.com/underbek/datamapper/models"
	"golang.org/x/tools/imports"
)

const (
	convertorFilePath                  = "templates/convertor.temp"
	errorConversionFilePath            = "templates/error_conversion.temp"
	pointerCheckFilePath               = "templates/pointer_check.temp"
	pointerConversionFilePath          = "templates/pointer_conversion.temp"
	pointerToPointerConversionFilePath = "templates/pointer_to_pointer_conversion.temp"
)

//go:embed templates
var templates embed.FS

func fillTemplate[T []byte | string](tempPath string, data map[string]any) (T, error) {
	temp, err := template.ParseFS(templates, tempPath)

	var res T
	if err != nil {
		return res, err
	}

	buf := bytes.Buffer{}
	err = temp.Execute(&buf, data)
	if err != nil {
		return res, err
	}

	return T(buf.Bytes()), nil
}

func createConvertor(res result) ([]byte, error) {
	data := map[string]any{
		"packageName":   res.pkgName,
		"fromName":      res.fromName,
		"toName":        res.toName,
		"fromTag":       res.fromTag,
		"toTag":         res.toTag,
		"convertorName": res.convertorName,
		"fields":        res.fields,
		"imports":       res.imports,
		"withError":     res.withError,
		"conversions":   res.conversions,
	}

	body, err := fillTemplate[[]byte](convertorFilePath, data)
	if err != nil {
		return nil, err
	}

	content, err := imports.Process(res.pkgPath, body, nil)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func getPointerCheck(from, to models.Field, fromName, toName string) (string, error) {
	data := map[string]any{
		"fromModelName": fromName,
		"toModelName":   toName,
		"fromFieldName": from.Name,
		"toFieldName":   to.Name,
	}

	return fillTemplate[string](pointerCheckFilePath, data)
}

func getErrorConversion(fromFiledName, toModelName, conversionFunction string) (string, error) {
	data := map[string]any{
		"toModelName":        toModelName,
		"fromFieldName":      fromFiledName,
		"conversionFunction": conversionFunction,
	}

	return fillTemplate[string](errorConversionFilePath, data)
}

func getPointerConversion(fromFieldName string, conversionFunction string) (string, error) {
	data := map[string]any{
		"fromFieldName":      fromFieldName,
		"conversionFunction": conversionFunction,
	}

	return fillTemplate[string](pointerConversionFilePath, data)
}

func getPointerToPointerConversion(fromFieldName, toModelName, toFullFieldType, conversionFunction string, isError bool,
) (string, error) {
	data := map[string]any{
		"fromFieldName":      fromFieldName,
		"toModelName":        toModelName,
		"toFullFieldType":    toFullFieldType,
		"conversionFunction": conversionFunction,
		"isError":            isError,
	}

	return fillTemplate[string](pointerToPointerConversionFilePath, data)
}
