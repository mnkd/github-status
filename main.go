package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	ExitCodeOK int = iota
	ExitCodeError
)

var (
	Version  string
	Revision string
)

var app = App{}

func init() {
	var path string
	var version, good bool

	flag.BoolVar(&version, "v", false, "prints current github-status version")
	flag.StringVar(&path, "c", "", "/path/to/config.json (default: $HOME/.config/github-status/config.json)")
	flag.BoolVar(&good, "good", false, "Post Good Status (default: false)")
	flag.Parse()

	if version {
		fmt.Fprintln(os.Stdout, "Version:", Version)
		fmt.Fprintln(os.Stdout, "Revision:", Revision)
		os.Exit(ExitCodeOK)
	}

	config, err := NewConfig(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(ExitCodeError)
	}

	app.Good = good
	app.SlackConfig = config.SlackWebhook
}

func main() {
	os.Exit(app.Run())
}
