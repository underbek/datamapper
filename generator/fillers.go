package generator

import (
	"bytes"
	"embed"
	"text/template"

	"golang.org/x/tools/imports"
)

const (
	convertorFilePath                  = "templates/convertor.temp"
	errorConversionFilePath            = "templates/error_conversion.temp"
	pointerCheckFilePath               = "templates/pointer_check.temp"
	pointerConversionFilePath          = "templates/pointer_conversion.temp"
	pointerToPointerConversionFilePath = "templates/pointer_to_pointer_conversion.temp"
	sliceConversionFilePath            = "templates/slice_conversion.temp"
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
	imps := make([]string, 0, len(res.packages))
	for pkg := range res.packages {
		imps = append(imps, pkg.Import())
	}

	data := map[string]any{
		"packageName":   res.pkg.Name,
		"fromName":      res.fromName,
		"toName":        res.toName,
		"fromTag":       res.fromTag,
		"toTag":         res.toTag,
		"convertorName": res.convertorName,
		"fields":        res.fields,
		"imports":       filterAndSortImports(res.pkg, imps),
		"withError":     res.withError,
		"conversions":   res.conversions,
	}

	body, err := fillTemplate[[]byte](convertorFilePath, data)
	if err != nil {
		return nil, err
	}

	content, err := imports.Process(res.pkg.Path, body, nil)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func getPointerCheck(fromFieldFullName, fromFieldName, toFieldName, fromName, toName string) (string, error) {
	data := map[string]any{
		"fromModelName":     fromName,
		"fromFieldName":     fromFieldName,
		"toModelName":       toName,
		"fromFieldFullName": fromFieldFullName,
		"toFieldName":       toFieldName,
	}

	return fillTemplate[string](pointerCheckFilePath, data)
}

func getErrorConversion(fromFieldFullName, toModelName, conversionFunction string) (string, error) {
	data := map[string]any{
		"toModelName":        toModelName,
		"fromFieldFullName":  fromFieldFullName,
		"conversionFunction": conversionFunction,
	}

	return fillTemplate[string](errorConversionFilePath, data)
}

func getPointerConversion(fromFieldFullName string, conversionFunction string) (string, error) {
	data := map[string]any{
		"fromFieldFullName":  fromFieldFullName,
		"conversionFunction": conversionFunction,
	}

	return fillTemplate[string](pointerConversionFilePath, data)
}

func getPointerToPointerConversion(fromFieldResName, fromFieldFullName, toModelName, toFullFieldType,
	conversionFunction string, isError bool) (string, error) {

	data := map[string]any{
		"fromFieldResName":   fromFieldResName,
		"fromFieldFullName":  fromFieldFullName,
		"toModelName":        toModelName,
		"toFullFieldType":    toFullFieldType,
		"conversionFunction": conversionFunction,
		"isError":            isError,
	}

	return fillTemplate[string](pointerToPointerConversionFilePath, data)
}

func getSliceConversion(fromFieldName, toItemTypeName, assigment string, conversions []string) (string, error) {
	data := map[string]any{
		"fromFieldName":  fromFieldName,
		"toItemTypeName": toItemTypeName,
		"assigment":      assigment,
		"conversions":    conversions,
	}

	return fillTemplate[string](sliceConversionFilePath, data)
}
