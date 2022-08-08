package options

import (
	"github.com/jessevdk/go-flags"
)

//nolint:lll
type Options struct {
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

func ParseOptions() (Options, error) {
	var options Options
	_, err := flags.Parse(&options)
	return options, err
}
