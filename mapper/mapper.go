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
	//TODO: parse or copy embed sources
	funcs, err := parser.ParseConversionFunctionsByPackage(internalConvertsPackagePath)
	if err != nil {
		return fmt.Errorf("parse internal conversion functions error: %w", err)
	}

	cfAliases := map[string]string{}
	userFuncs := make(models.Functions)

	if len(opts.CFSources) != 0 {
		for _, optSource := range opts.CFSources {
			source, cfAlias := parseSourceOption(optSource)
			res, err := parser.ParseConversionFunctionsByPackage(source)
			if err != nil {
				return fmt.Errorf("parse user conversion functions error: %w", err)
			}

			for key, function := range res {
				cfAliases[function.Package.Path] = cfAlias
				userFuncs[key] = function
			}
		}
	}

	for _, opt := range opts.Options {
		fromSource, fromAlias := parseSourceOption(opt.FromSource)

		structs, err := parser.ParseModelsByPackage(fromSource)
		if err != nil {
			return fmt.Errorf("parse models error: %w", err)
		}

		fromName, isFromPointer := parseModelName(opt.FromName)
		from, ok := structs[fromName]
		if !ok {
			return fmt.Errorf(" %w: source model %s from %s", ErrNotFoundStruct, opt.FromName, opt.FromSource)
		}
		from.Type.Pointer = isFromPointer

		toSource, toAlias := parseSourceOption(opt.ToSource)
		structs, err = parser.ParseModelsByPackage(toSource)
		if err != nil {
			return fmt.Errorf("parse models error: %w", err)
		}

		toName, isToPointer := parseModelName(opt.ToName)
		to, ok := structs[toName]
		if !ok {
			return fmt.Errorf("%w: to model %s from %s", ErrNotFoundStruct, opt.ToName, opt.ToSource)
		}
		to.Type.Pointer = isToPointer

		aliases := map[string]string{
			from.Type.Package.Path: fromAlias,
			to.Type.Package.Path:   toAlias,
		}

		from.Fields = utils.FilterFields(opt.FromTag, from.Fields)
		if len(from.Fields) == 0 {
			return fmt.Errorf(
				"%w: source model %s does not contain tag %s",
				ErrNotFoundTag,
				opt.FromName,
				opt.FromTag,
			)
		}

		to.Fields = utils.FilterFields(opt.ToTag, to.Fields)
		if len(to.Fields) == 0 {
			return fmt.Errorf("%w: to model %s does not contain tag %s", ErrNotFoundTag, opt.ToName, opt.ToTag)
		}

		for cfPath, alias := range cfAliases {
			aliases[cfPath] = alias
		}

		// set aliases
		setPackageAliasToStruct(&from, aliases)
		setPackageAliasToStruct(&to, aliases)
		userFuncs = setPackageAliasToFunctions(userFuncs, aliases)
		for key, function := range userFuncs {
			funcs[key] = function
		}

		err = os.MkdirAll(path.Dir(opt.Destination), os.ModePerm)
		if err != nil {
			return fmt.Errorf("create destination dir %s error: %w", path.Dir(opt.Destination), err)
		}

		pkg, err := parser.ParseDestinationPackage(opt.Destination)
		if err != nil {
			return fmt.Errorf("parse destination package %s error: %w", opt.Destination, err)
		}

		var convertors []string
		pkgs, convertor, err := generator.GenerateConvertor(from, to, pkg, funcs)
		if err != nil {
			return fmt.Errorf("generate convertor error: %w", err)
		}
		convertors = append(convertors, convertor)

		if opt.Invert {
			invertPkgs, convertor, err := generator.GenerateConvertor(to, from, pkg, funcs)
			if err != nil {
				return fmt.Errorf("generate convertor error: %w", err)
			}
			convertors = append(convertors, convertor)

			for key := range invertPkgs {
				pkgs[key] = struct{}{}
			}
		}

		err = generator.CreateConvertorSource(pkg, pkgs, convertors, opt.Destination)
		if err != nil {
			return fmt.Errorf("create convertor source error: %w", err)
		}
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

func parseModelName(modelName string) (string, bool) {
	if strings.HasPrefix(modelName, "*") {
		return strings.TrimPrefix(modelName, "*"), true
	}

	return modelName, false
}
