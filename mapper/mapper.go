package mapper

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/underbek/datamapper/generator"
	"github.com/underbek/datamapper/models"
	"github.com/underbek/datamapper/options"
	"github.com/underbek/datamapper/parser"
	"github.com/underbek/datamapper/utils"
)

const internalConvertsPackagePath = "github.com/underbek/datamapper/converts"

var (
	ErrNotFoundStruct = errors.New("not found struct error")
	ErrNotFoundTag    = errors.New("not found tag error")
)

func MapModels(opts options.Options) error {
	fromSource, fromAlias := parseSourceOption(opts.FromSource)

	structs, err := parser.ParseModelsByPackage(fromSource)
	if err != nil {
		return fmt.Errorf("parse models error: %w", err)
	}

	from, ok := structs[opts.FromName]
	if !ok {
		return fmt.Errorf(" %w: source model %s from %s", ErrNotFoundStruct, opts.FromName, opts.FromSource)
	}

	toSource, toAlias := parseSourceOption(opts.ToSource)
	structs, err = parser.ParseModelsByPackage(toSource)
	if err != nil {
		return fmt.Errorf("parse models error: %w", err)
	}

	to, ok := structs[opts.ToName]
	if !ok {
		return fmt.Errorf("%w: to model %s from %s", ErrNotFoundStruct, opts.ToName, opts.ToSource)
	}

	aliases := map[string]string{
		from.Type.Package.Path: fromAlias,
		to.Type.Package.Path:   toAlias,
	}

	from.Fields = utils.FilterFields(opts.FromTag, from.Fields)
	if len(from.Fields) == 0 {
		return fmt.Errorf(
			"%w: source model %s does not contain tag %s",
			ErrNotFoundTag,
			opts.FromName,
			opts.FromTag,
		)
	}

	to.Fields = utils.FilterFields(opts.ToTag, to.Fields)
	if len(to.Fields) == 0 {
		return fmt.Errorf("%w: to model %s does not contain tag %s", ErrNotFoundTag, opts.ToName, opts.ToTag)
	}

	//TODO: parse or copy embed sources
	funcs, err := parser.ParseConversionFunctionsByPackage(internalConvertsPackagePath)
	if err != nil {
		return fmt.Errorf("parse internal conversion functions error: %w", err)
	}

	userFuncs := make(models.Functions)

	if len(opts.UserCFSources) != 0 {
		for _, optSource := range opts.UserCFSources {
			source, cfAlias := parseSourceOption(optSource)
			res, err := parser.ParseConversionFunctionsByPackage(source)
			if err != nil {
				return fmt.Errorf("parse user conversion functions error: %w", err)
			}

			for key, function := range res {
				aliases[function.Package.Path] = cfAlias
				userFuncs[key] = function
			}
		}
	}

	// set aliases
	setPackageAliasToStruct(&from, aliases)
	setPackageAliasToStruct(&to, aliases)
	userFuncs = setPackageAliasToFunctions(userFuncs, aliases)
	for key, function := range userFuncs {
		funcs[key] = function
	}

	err = os.MkdirAll(path.Dir(opts.Destination), os.ModePerm)
	if err != nil {
		return fmt.Errorf("create destination dir %s error: %w", path.Dir(opts.Destination), err)
	}

	pkg, err := parser.ParseDestinationPackage(opts.Destination)
	if err != nil {
		return fmt.Errorf("parse destination package %s error: %w", opts.Destination, err)
	}

	var convertors []string
	pkgs, convertor, err := generator.GenerateConvertor(from, to, pkg, funcs)
	if err != nil {
		return fmt.Errorf("generate convertor error: %w", err)
	}
	convertors = append(convertors, convertor)

	if opts.Invert {
		invertPkgs, convertor, err := generator.GenerateConvertor(to, from, pkg, funcs)
		if err != nil {
			return fmt.Errorf("generate convertor error: %w", err)
		}
		convertors = append(convertors, convertor)

		for key := range invertPkgs {
			pkgs[key] = struct{}{}
		}
	}

	err = generator.CreateConvertorSource(pkg, pkgs, convertors, opts.Destination)
	if err != nil {
		return fmt.Errorf("create convertor source error: %w", err)
	}

	return nil
}

func setPackageAlias(p *models.Package, aliases map[string]string) {
	p.Alias = aliases[p.Path]
}

func setPackageAliasToStruct(m *models.Struct, aliases map[string]string) {
	setPackageAlias(&m.Type.Package, aliases)
	for i := range m.Fields {
		setPackageAlias(&m.Fields[i].Type.Package, aliases)
		switch m.Fields[i].Type.Kind {
		case models.SliceType:
			additional := m.Fields[i].Type.Additional.(models.SliceAdditional)
			setPackageAlias(&additional.InType.Package, aliases)
			m.Fields[i].Type.Additional = additional
		}
	}
}

func setPackageAliasToCfKey(key models.ConversionFunctionKey, aliases map[string]string) models.ConversionFunctionKey {
	setPackageAlias(&key.FromType.Package, aliases)
	setPackageAlias(&key.ToType.Package, aliases)

	return key
}

func setPackageAliasToCf(cf models.ConversionFunction, aliases map[string]string) models.ConversionFunction {
	setPackageAlias(&cf.Package, aliases)
	setPackageAlias(&cf.FromType.Package, aliases)
	setPackageAlias(&cf.ToType.Package, aliases)

	return cf
}

func setPackageAliasToFunctions(funcs models.Functions, aliases map[string]string) models.Functions {
	res := make(models.Functions)
	for key, cf := range funcs {
		res[setPackageAliasToCfKey(key, aliases)] = setPackageAliasToCf(cf, aliases)
	}
	return res
}

func parseSourceOption(optSource string) (string, string) {
	res := strings.Split(optSource, ":")
	if len(res) == 0 {
		return "", ""
	}

	if len(res) == 1 {
		return res[0], ""
	}

	return res[0], res[1]
}
