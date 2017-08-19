package main

import (
	"flag"
	"fmt"
	"os"
)

// Exit codes
const (
	ExitCodeOK int = iota
	ExitCodeError
)

var (
	version  string
	revision string
)

var app = App{}

func init() {
	var path string
	var printVersion, good bool

	flag.BoolVar(&printVersion, "v", false, "prints current github-status version")
	flag.StringVar(&path, "c", "", "/path/to/config.json (default: $HOME/.config/github-status/config.json)")
	flag.BoolVar(&good, "good", false, "Post Good Status (default: false)")
	flag.Parse()

	if printVersion {
		fmt.Fprintln(os.Stdout, "Version:", version)
		fmt.Fprintln(os.Stdout, "Revision:", revision)
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
