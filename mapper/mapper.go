package mapper

import (
	"errors"
	"fmt"

	"github.com/underbek/datamapper/generator"
	"github.com/underbek/datamapper/models"
	"github.com/underbek/datamapper/options"
	"github.com/underbek/datamapper/parser"
	"github.com/underbek/datamapper/utils"
	"golang.org/x/exp/maps"
)

const internalConvertsPackagePath = "github.com/underbek/datamapper/converts"

var (
	ErrNotFoundStruct = errors.New("not found struct error")
	ErrNotFoundTag    = errors.New("not found tag error")
)

func MapModels(opts options.Options) error {
	structs, err := parser.ParseModelsByPackage(opts.FromSource)
	if err != nil {
		return fmt.Errorf("parse models error: %w", err)
	}

	from, ok := structs[opts.FromName]
	if !ok {
		return fmt.Errorf(" %w: source model %s from %s", ErrNotFoundStruct, opts.FromName, opts.FromSource)
	}

	structs, err = parser.ParseModelsByPackage(opts.ToSource)
	if err != nil {
		return fmt.Errorf("parse models error: %w", err)
	}

	to, ok := structs[opts.ToName]
	if !ok {
		return fmt.Errorf("%w: to model %s from %s", ErrNotFoundStruct, opts.ToName, opts.ToSource)
	}

	//TODO: add cf aliases
	aliases := map[string]string{
		from.Package.Path: opts.FromPackageAlias,
		to.Package.Path:   opts.ToPackageAlias,
	}

	setPackageAliasToStruct(&from, aliases)
	setPackageAliasToStruct(&to, aliases)

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

	if len(opts.UserCFSources) != 0 {
		for _, source := range opts.UserCFSources {
			userFuncs, err := parser.ParseConversionFunctionsByPackage(source)
			if err != nil {
				return fmt.Errorf("parse user conversion functions error: %w", err)
			}

			// TODO: set alias to each cf
			res := make(models.Functions)
			for key, cf := range userFuncs {
				setPackageAliasToCfKey(&key, aliases)
				setPackageAliasToCf(&cf, aliases)
				res[key] = cf
			}

			maps.Copy(funcs, res)
		}
	}

	err = generator.CreateConvertor(from, to, opts.Destination, funcs)
	if err != nil {
		return fmt.Errorf("generate convertor error: %w", err)
	}

	return nil
}

func setPackageAliasToStruct(m *models.Struct, aliases map[string]string) {
	m.Package.Alias = aliases[m.Package.Path]
	for i := range m.Fields {
		m.Fields[i].Type.Package.Alias = aliases[m.Fields[i].Type.Package.Path]
	}
}

func setPackageAliasToCfKey(key *models.ConversionFunctionKey, aliases map[string]string) {
	key.FromType.Package.Alias = aliases[key.FromType.Package.Path]
	key.ToType.Package.Alias = aliases[key.ToType.Package.Path]
}

func setPackageAliasToCf(cf *models.ConversionFunction, aliases map[string]string) {
	cf.Package.Alias = aliases[cf.Package.Path]
	cf.FromType.Package.Alias = aliases[cf.FromType.Package.Path]
	cf.ToType.Package.Alias = aliases[cf.ToType.Package.Path]
}
