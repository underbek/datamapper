package main

import (
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Source      string `short:"s" long:"source" description:"Source file path" required:"true"`
	Destination string `short:"d" long:"destination" description:"Destination file path" required:"true"`
	Package     string `short:"p" long:"package" description:"Destination package name" required:"true"`
	FromName    string `long:"from" description:"Model from name" required:"true"`
	ToName      string `long:"to" description:"Model to name" required:"true"`
	FromTag     string `long:"from-tag" description:"Model from tag" required:"true"`
	ToTag       string `long:"to-tag" description:"Model to tag" required:"true"`
}

func parseOptions() (Options, error) {
	var options Options
	_, err := flags.Parse(&options)
	return options, err
}
