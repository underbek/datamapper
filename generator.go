package main

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"html/template"
	"os"
)

type FieldsPair struct {
	FromName   string
	FromType   string
	ToName     string
	ToType     string
	Conversion string
}

const convertorFilePath = "templates/convertor.temp"

//go:embed templates
var templates embed.FS

func createConvertor(from, to Struct, dest string, packageName string) error {
	temp, err := template.ParseFS(templates, convertorFilePath)
	if err != nil {
		return err
	}

	fields := createModelsPair(from, to)

	data := map[string]any{
		"packageName": packageName,
		"fromName":    from.Name,
		"toName":      to.Name,
		"fields":      fields,
	}

	buf := bytes.Buffer{}
	err = temp.Execute(&buf, data)
	if err != nil {
		return err
	}

	content, err := format.Source(buf.Bytes())
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

func createModelsPair(from, to Struct) []FieldsPair {
	var res []FieldsPair

	fromFields := make(map[string]Field)
	for _, field := range from.Fields {
		fromFields[field.Tags[0].Value] = field
	}

	for _, toField := range to.Fields {
		fromField, ok := fromFields[toField.Tags[0].Value]
		if !ok {
			//TODO: warning or error politics
			continue
		}

		//TODO: find convert function if to and from types are different
		conversion := fmt.Sprintf("from.%s", fromField.Name)

		res = append(res, FieldsPair{
			FromName:   fromField.Name,
			FromType:   fromField.Type,
			ToName:     toField.Name,
			ToType:     toField.Type,
			Conversion: conversion,
		})
	}

	return res
}
