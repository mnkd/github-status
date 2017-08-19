package main

import (
	"encoding/json"
	"io/ioutil"
	"os/user"
	"path/filepath"

	slack "github.com/mnkd/slackposter"
	"github.com/pkg/errors"
)

// Config is the configuration.
type Config struct {
	SlackWebhook slack.Config `json:"slack_webhook"`
}

// NewConfig returns a new Config instance with config.json path.
func NewConfig(path string) (Config, error) {
	usr, err := user.Current()
	if err != nil {
		return Config{}, errors.Wrap(err, "failed user.Current()")
	}

	// Decide config.json path
	if len(path) == 0 {
		path = filepath.Join(usr.HomeDir, "/.config/github-status/config.json")
	} else {
		p, absErr := filepath.Abs(path)
		if absErr != nil {
			return Config{}, errors.Wrapf(absErr, "failed filepath.Abs: %v", path)
		}
		path = p
	}

	str, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, errors.Wrap(err, "read config.json")
	}

	var config Config
	if err := json.Unmarshal(str, &config); err != nil {
		return config, errors.Wrap(err, "failed to unmarshal config.json")
	}

	return config, nil
}
