package options

import (
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

type Option struct {
	Destination string
	FromName    string
	FromTag     string
	FromSource  string
	ToName      string
	ToTag       string
	ToSource    string
	Invert      bool
}

type Options struct {
	CFSources []string
	Options   []Option
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

	return Options{
		CFSources: params.UserCFSources,
		Options: []Option{
			{
				Destination: params.Destination,
				FromName:    params.FromName,
				FromTag:     params.FromTag,
				FromSource:  params.FromSource,
				ToName:      params.ToName,
				ToTag:       params.ToTag,
				ToSource:    params.ToSource,
				Invert:      params.Invert,
			},
		},
	}, nil
}

func ParseOptions() (Options, error) {
	return parseOptions()
}
