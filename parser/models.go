package parser

import (
	"go/build"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"github.com/underbek/datamapper/logger"
	"github.com/underbek/datamapper/models"
	"github.com/underbek/datamapper/utils"
)

var (
	modelsCache = make(map[string]map[string]models.Struct)
)

func ParseModelsByPackage(lg logger.Logger, source string) (map[string]models.Struct, error) {
	_, err := os.Stat(source)
	if err == nil {
		return ParseModels(lg, source)
	}

	if !os.IsNotExist(err) {
		return nil, err
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	p, err := build.Import(source, wd, build.FindOnly)
	if err != nil {
		return nil, err
	}

	return ParseModels(lg, p.Dir)
}

func ParseModels(lg logger.Logger, source string) (map[string]models.Struct, error) {
	absSourcePath, err := filepath.Abs(source)
	if err != nil {
		return nil, err
	}

	if structs, ok := modelsCache[absSourcePath]; ok {
		return structs, nil
	}

	pkg, err := utils.LoadPackage(lg, source)
	if err != nil {
		return nil, err
	}

	structs := make(map[string]models.Struct)

	names := pkg.Types.Scope().Names()
	for _, name := range names {
		obj := pkg.Types.Scope().Lookup(name)

		fset := pkg.Fset.Position(obj.Pos())
		if !strings.Contains(fset.Filename, absSourcePath) {
			continue
		}

		currType, ok := obj.(*types.TypeName)
		if !ok {
			continue
		}

		if !currType.Exported() {
			continue
		}

		currStruct, ok := currType.Type().Underlying().(*types.Struct)
		if !ok {
			continue
		}

		if currStruct.NumFields() == 0 {
			continue
		}

		fields := make([]models.Field, 0, currStruct.NumFields())
		for i := 0; i < currStruct.NumFields(); i++ {
			field := currStruct.Field(i)
			tts, err := parseType(field.Type())
			if err != nil {
				return nil, err
			}

			if len(tts) != 1 {
				continue
			}

			fields = append(fields, models.Field{
				Name: field.Name(),
				Type: tts[0].Type,
				Tags: parseTag(currStruct.Tag(i)),
			})
		}

		structs[currType.Name()] = models.Struct{
			Type: models.Type{
				Name: currType.Name(),
				Package: models.Package{
					Name: pkg.Name,
					Path: pkg.PkgPath,
				},
				Kind: models.StructType,
			},
			Fields: fields,
		}
	}

	modelsCache[absSourcePath] = structs

	return structs, nil
}
