package options

import (
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Destination  string `short:"d" long:"destination" description:"Destination file path" required:"true"`
	UserCFSource string `long:"cf" description:"User conversion functions source" required:"false"`
	FromSource   string `long:"from-source" description:"From source file path" required:"false"`
	ToSource     string `long:"to-source" description:"Source file path" required:"false"`
	FromName     string `long:"from" description:"Model from name" required:"true"`
	ToName       string `long:"to" description:"Model to name" required:"true"`
	FromTag      string `long:"from-tag" description:"Model from tag" required:"true"`
	ToTag        string `long:"to-tag" description:"Model to tag" required:"true"`
}

func ParseOptions() (Options, error) {
	var options Options
	_, err := flags.Parse(&options)
	return options, err
}
