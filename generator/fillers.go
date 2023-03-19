package generator

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/underbek/datamapper/models"
	"golang.org/x/tools/imports"
)

const (
	convertorSourceFilePath            = "templates/convertor_source.temp"
	convertorFilePath                  = "templates/convertor.temp"
	sliceConvertorFilePath             = "templates/slice_convertor.temp"
	errorConversionFilePath            = "templates/error_conversion.temp"
	pointerCheckFilePath               = "templates/pointer_check.temp"
	pointerConversionFilePath          = "templates/pointer_conversion.temp"
	pointerToPointerConversionFilePath = "templates/pointer_to_pointer_conversion.temp"
	sliceConversionFilePath            = "templates/slice_conversion.temp"
	convertErrorFilePath               = "templates/convert_error.temp"
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

func fillConvertorsSource(pkg models.Package, packages models.Packages, convertors []string) ([]byte, error) {
	imps := make([]string, 0, len(packages))
	for pkg := range packages {
		imps = append(imps, pkg.Import())
	}

	data := map[string]any{
		"packageName": pkg.Name,
		"imports":     filterAndSortImports(pkg.Import(), imps),
		"convertors":  convertors,
	}

	body, err := fillTemplate[[]byte](convertorSourceFilePath, data)
	if err != nil {
		return nil, err
	}

	content, err := imports.Process(pkg.Path, body, nil)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func fillConvertor(res result) (string, error) {
	data := map[string]any{
		"fromName":      res.fromName,
		"toName":        res.toName,
		"fromTag":       res.fromTag,
		"toTag":         res.toTag,
		"convertorName": res.convertorName,
		"fields":        res.fields,
		"withError":     res.withError,
		"conversions":   res.conversions,
		"resName":       strings.Replace(res.toName, "*", "&", 1),
	}

	return fillTemplate[string](convertorFilePath, data)
}

func fillSliceConvertor(res sliceResult) (string, error) {
	data := map[string]any{
		"fromName":      res.fromName,
		"toName":        res.toName,
		"convertorName": res.convertorName,
		"withError":     res.withError,
		"conversion":    res.conversion,
	}

	return fillTemplate[string](sliceConvertorFilePath, data)
}

func getPointerCheck(fromFullName, toModelName, err string, isError bool) (string, error) {
	data := map[string]any{
		"fromFullName": fromFullName,
		"resValue":     nilOrDefault(toModelName),
		"error":        err,
		"isError":      isError,
	}

	return fillTemplate[string](pointerCheckFilePath, data)
}

func nilOrDefault(fullName string) string {
	if strings.HasPrefix(fullName, "*") {
		return "nil"
	}

	return fmt.Sprintf("%s{}", fullName)
}

func getErrorConversion(fromFieldFullName, toModelName, conversionFunction, err string) (string, error) {
	data := map[string]any{
		"resValue":           nilOrDefault(toModelName),
		"fromFieldFullName":  fromFieldFullName,
		"conversionFunction": conversionFunction,
		"error":              err,
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
	conversionFunction, err string, isError bool) (string, error) {

	data := map[string]any{
		"fromFieldResName":   fromFieldResName,
		"fromFieldFullName":  fromFieldFullName,
		"resValue":           nilOrDefault(toModelName),
		"toFullFieldType":    toFullFieldType,
		"conversionFunction": conversionFunction,
		"error":              err,
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

func getConvertError(fromTypeName, fromFieldName, toTypeName, toFieldName string) (string, error) {
	data := map[string]any{
		"fromTypeName":  fromTypeName,
		"fromFieldName": fromFieldName,
		"toTypeName":    toTypeName,
		"toFieldName":   toFieldName,
	}

	return fillTemplate[string](convertErrorFilePath, data)
}
