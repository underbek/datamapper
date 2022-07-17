package generator

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"os"
	"strings"
	"text/template"

	"github.com/underbek/datamapper/models"
	"github.com/underbek/datamapper/utils"
	"golang.org/x/exp/maps"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ConvertorType = string
type ImportType = string

type FieldsPair struct {
	FromName       string
	FromType       string
	ToName         string
	ToType         string
	Assignment     string
	Conversions    []string
	WithError      bool
	PointerToValue bool
}

type result struct {
	fields        []FieldsPair
	imports       []string
	conversations []string
}

const (
	convertorFilePath           = "templates/convertor.temp"
	errorConversationFilePath   = "templates/error_conversion.temp"
	pointerCheckFilePath        = "templates/pointer_check.temp"
	pointerConversationFilePath = "templates/pointer_conversion.temp"
)

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
	if from.PackagePath != pkg.PkgPath {
		convertorName += cases.Title(language.Und, cases.NoLower).String(from.PackageName)
	}
	convertorName += from.Name
	convertorName += "To"
	if to.PackagePath != pkg.PkgPath {
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

	fromName := getFullStructName(from, pkg.PkgPath)
	toName := getFullStructName(to, pkg.PkgPath)

	return createConvertor(pkgName, fromName, toName, convertorName, pkg.PkgPath, res)
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

		pair, packs, err := getFieldsPair(fromField, toField, from, to, pkgPath, functions)
		if err != nil {
			return result{}, err
		}

		for _, pack := range packs {
			if pack != "" {
				imports[pack] = struct{}{}
			}
		}

		fields = append(fields, pair)
	}

	return result{
		fields:        fields,
		imports:       maps.Keys(imports),
		conversations: fillConversations(fields),
	}, nil
}

func getFieldsPair(from, to models.Field, fromModel, toModel models.Struct, pkgPath string, functions models.Functions,
) (FieldsPair, []ImportType, error) {

	// TODO: check package
	if from.Type.Name == to.Type.Name {
		res, pkg, err := getFieldsPairBySameTypes(from, to, fromModel.Name, toModel.Name)
		if err != nil {
			return FieldsPair{}, nil, err
		}
		return res, []ImportType{pkg}, nil
	}

	key := models.ConversionFunctionKey{
		FromType: from.Type,
		ToType:   to.Type,
	}

	//TODO: Use conversion functions with pointers
	key.FromType.Pointer = false
	key.ToType.Pointer = false

	cf, ok := functions[key]

	if !ok {
		return FieldsPair{}, nil, fmt.Errorf(
			"not found convertor function for types %s -> %s by %s field",
			from.Type.Name,
			to.Type.Name,
			from.Name,
		)
	}

	typeParams := getTypeParams(cf, from.Type, to.Type)

	res := FieldsPair{
		FromName:  from.Name,
		FromType:  from.Type.Name,
		ToName:    to.Name,
		ToType:    to.Type.Name,
		WithError: cf.WithError,
	}

	ptr := ""
	imports := []ImportType{cf.PackagePath}

	if from.Type.Pointer {
		pointerCheck, err := getPointerCheck(from, to,
			getFullStructName(fromModel, pkgPath),
			getFullStructName(toModel, pkgPath),
		)
		if err != nil {
			return FieldsPair{}, nil, err
		}

		res.Conversions = append(res.Conversions, pointerCheck)
		ptr = "*"
		imports = append(imports, "fmt")
		res.PointerToValue = true
	}

	conversation := fmt.Sprintf("%s.%s%s(%sfrom.%s)", cf.PackageName, cf.Name, typeParams, ptr, from.Name)
	if cf.PackagePath == pkgPath {
		conversation = fmt.Sprintf("%s%s(%sfrom.%s)", cf.Name, typeParams, ptr, from.Name)
	}

	if to.Type.Pointer {
		pointerConversion, err := getPointerConversion(from.Name, conversation)
		if err != nil {
			return FieldsPair{}, nil, err
		}

		res.Conversions = append(res.Conversions, pointerConversion)
		conversation = fmt.Sprintf("&from%s", from.Name)
	}

	if !cf.WithError {
		res.Assignment = conversation
		return res, imports, nil
	}

	errorConversation, err := getErrorConversion(from.Name, getFullStructName(toModel, pkgPath), conversation)
	if err != nil {
		return FieldsPair{}, nil, err
	}

	res.Assignment = fmt.Sprintf("from%s", from.Name)
	res.Conversions = append(res.Conversions, errorConversation)

	return res, imports, nil
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

func fillConversations(fields []FieldsPair) []string {
	var res []string
	for _, field := range fields {
		res = append(res, field.Conversions...)
	}

	return res
}

func isReturnError(fields []FieldsPair) bool {
	for _, field := range fields {
		if field.WithError || field.PointerToValue {
			return true
		}
	}

	return false
}

func getFieldsPairBySameTypes(from, to models.Field, fromName, toName string) (FieldsPair, ImportType, error) {
	res := FieldsPair{
		FromName: from.Name,
		FromType: from.Type.Name,
		ToName:   to.Name,
		ToType:   to.Type.Name,
	}

	if from.Type.Pointer == to.Type.Pointer {
		res.Assignment = fmt.Sprintf("from.%s", from.Name)
		return res, "", nil
	}

	if to.Type.Pointer {
		res.Assignment = fmt.Sprintf("&from.%s", from.Name)
		return res, "", nil
	}

	res.PointerToValue = true
	res.Assignment = fmt.Sprintf("*from.%s", from.Name)

	conversion, err := getPointerCheck(from, to, fromName, toName)
	if err != nil {
		return FieldsPair{}, "", err
	}

	res.Conversions = append(res.Conversions, conversion)

	return res, "fmt", nil
}

func createConvertor(pkgName, fromName, toName, convertorName, pkgPath string, res result) ([]byte, error) {
	temp, err := template.ParseFS(templates, convertorFilePath)
	if err != nil {
		return nil, err
	}

	data := map[string]any{
		"packageName":   pkgName,
		"fromName":      fromName,
		"toName":        toName,
		"convertorName": convertorName,
		"fields":        res.fields,
		"imports":       filterImports(pkgPath, res.imports),
		"withError":     isReturnError(res.fields),
		"conversations": res.conversations,
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

func getPointerCheck(from, to models.Field, fromName, toName string) (string, error) {
	temp, err := template.ParseFS(templates, pointerCheckFilePath)
	if err != nil {
		return "", err
	}

	data := map[string]any{
		"fromModelName": fromName,
		"toModelName":   toName,
		"fromFieldName": from.Name,
		"toFieldName":   to.Name,
	}

	buf := bytes.Buffer{}
	err = temp.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func getErrorConversion(fromFiledName, toModelName, conversionFunction string) (string, error) {
	temp, err := template.ParseFS(templates, errorConversationFilePath)
	if err != nil {
		return "", err
	}

	data := map[string]any{
		"toModelName":        toModelName,
		"fromFieldName":      fromFiledName,
		"conversionFunction": conversionFunction,
	}

	buf := bytes.Buffer{}
	err = temp.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func getPointerConversion(fromFieldName string, conversionFunction string) (string, error) {
	temp, err := template.ParseFS(templates, pointerConversationFilePath)
	if err != nil {
		return "", err
	}

	data := map[string]any{
		"fromFieldName":      fromFieldName,
		"conversionFunction": conversionFunction,
	}

	buf := bytes.Buffer{}
	err = temp.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func getFullStructName(model models.Struct, pkgPath string) string {
	if model.PackagePath != pkgPath {
		return fmt.Sprintf("%s.%s", model.PackageName, model.Name)
	}

	return model.Name
}
