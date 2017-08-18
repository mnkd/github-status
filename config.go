package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	slack "github.com/mnkd/slackposter"
)

type Config struct {
	SlackWebhook slack.Config `json:"slack_webhook"`
}

func NewConfig(path string) (Config, error) {
	var config Config

	usr, err := user.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not get current user.", err)
		return config, err
	}

	// Decide config.json path
	if len(path) == 0 {
		path = filepath.Join(usr.HomeDir, "/.config/github-status/config.json")
	} else {
		p, err := filepath.Abs(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not return absolute representation of path:", err, path)
			return config, err
		}
		path = p
	}

	str, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not read config.json. ", err)
		return config, err
	}

	if err := json.Unmarshal(str, &config); err != nil {
		fmt.Fprintln(os.Stderr, "Unmarshal config.json error:", err)
		return config, err
	}

	return config, nil
}
