package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/mnkd/slackposter"
)

type Config struct {
	SlackWebhook slack.Config `json:"slack_webhook"`
}

func NewConfig() (Config, error) {
	var config Config

	usr, err := user.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not get current user.", err)
		return config, err
	}

	path := filepath.Join(usr.HomeDir, "/.config/slackposter/config.json")
	str, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not read config.json. ", err)
		return config, err
	}

	if err := json.Unmarshal(str, &config); err != nil {
		fmt.Fprintln(os.Stderr, "JSON Unmarshal Error:", err)
		return config, err
	}

	return config, nil
}
