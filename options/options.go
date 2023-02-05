package options

import (
	"os"
	"strings"

	"github.com/creasty/defaults"
	"github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v3"
)

//nolint:lll
type Config struct {
	ConfigPath string `short:"c" long:"config" description:"Yaml config path" required:"false"`
	Flags
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
	Inverse       bool     `short:"i" long:"inverse" description:"Create direct and inverse conversions" required:"false"`
}

type Model struct {
	Name   string `yaml:"name"`
	Tag    string `yaml:"tag" default:"map"`
	Source string `yaml:"source" default:"."`
	Alias  string `yaml:"alias"`
}

type Option struct {
	From        Model  `yaml:"from"`
	To          Model  `yaml:"to"`
	Inverse     bool   `yaml:"inverse"`
	Destination string `yaml:"destination"`
}

type Options struct {
	ConversionFunctions []ConversionFunction `yaml:"conversion-functions"`
	Options             []Option             `yaml:"options"`
}

type ConversionFunction struct {
	Source string `yaml:"source"`
	Alias  string `yaml:"alias"`
}

func (m *Model) UnmarshalYAML(unmarshal func(interface{}) error) error {
	defaults.Set(m)

	type model Model
	if err := unmarshal((*model)(m)); err != nil {
		return err
	}

	return nil
}

func parseConfig(path string) (Options, error) {
	var opts Options
	file, err := os.Open(path)
	if err != nil {
		return Options{}, err
	}

	err = yaml.NewDecoder(file).Decode(&opts)
	if err != nil {
		return Options{}, err
	}

	return opts, nil
}

func parseFlags(params Flags) (Options, error) {
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
				Inverse: params.Inverse,
			},
		},
	}, nil
}

func ParseOptions() (Options, error) {
	var config Config
	_, err := flags.NewParser(&config, flags.HelpFlag|flags.PassDoubleDash).Parse()
	if config.ConfigPath != "" {
		return parseConfig(config.ConfigPath)
	}

	if err != nil {
		return Options{}, err
	}

	return parseFlags(config.Flags)
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
