package options

import (
	"strings"

	"github.com/jessevdk/go-flags"
)

//nolint:lll
type Config struct {
	Config string `short:"c" long:"config" description:"Yaml config path" required:"false"`
}

//nolint:lll
type Flags struct {
	Destination   string   `short:"d" long:"destination" description:"Destination file path" required:"true"`
	UserCFSources []string `long:"cf" description:"User conversion functions sources/packages. Can add package alias like {package_path}:{alias)" required:"false"`
	FromName      string   `long:"from" description:"Model from name" required:"true"`
	FromTag       string   `long:"from-tag" description:"Model from tag" default:"map" required:"false"`
	FromSource    string   `long:"from-source" description:"From model source/package. Can add package alias like {package_path}:{alias)" default:"." required:"false"`
	ToName        string   `long:"to" description:"Model to name" required:"true"`
	ToTag         string   `long:"to-tag" description:"Model to tag" default:"map" required:"false"`
	ToSource      string   `long:"to-source" description:"To model source/package. Can add package alias like {package_path}:{alias)" default:"." required:"false"`
	Invert        bool     `short:"i" long:"inverse" description:"Create direct and inverse conversions" required:"false"`
}

type Model struct {
	Name   string `yaml:"name"`
	Tag    string `yaml:"tag"`
	Source string `yaml:"source"`
	Alias  string `yaml:"alias"`
}

type Option struct {
	From        Model  `yaml:"from"`
	To          Model  `yaml:"to"`
	Invert      bool   `yaml:"invert"`
	Destination string `yaml:"dest"`
}

type Options struct {
	ConversionFunctions []ConversionFunction `yaml:"cf"`
	Options             []Option             `yaml:"options"`
}

type ConversionFunction struct {
	Source string `yaml:"source"`
	Alias  string `yaml:"alias"`
}

func parseConfig() (string, error) {
	var config Config
	_, err := flags.Parse(&config)
	return config.Config, err
}

func parseOptions() (Options, error) {
	var params Flags
	_, err := flags.Parse(&params)
	if err != nil {
		return Options{}, err
	}

	functions := make([]ConversionFunction, 0, len(params.UserCFSources))
	for _, opt := range params.UserCFSources {
		source, alias := parseSourceOption(opt)
		functions = append(functions, ConversionFunction{
			Source: source,
			Alias:  alias,
		})
	}

	fromSource, fromAlias := parseSourceOption(params.FromSource)
	toSource, toAlias := parseSourceOption(params.ToSource)

	return Options{
		ConversionFunctions: functions,
		Options: []Option{
			{
				Destination: params.Destination,
				From: Model{
					Name:   params.FromName,
					Tag:    params.FromTag,
					Source: fromSource,
					Alias:  fromAlias,
				},
				To: Model{
					Name:   params.ToName,
					Tag:    params.ToTag,
					Source: toSource,
					Alias:  toAlias,
				},
				Invert: params.Invert,
			},
		},
	}, nil
}

func ParseOptions() (Options, error) {
	return parseOptions()
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
