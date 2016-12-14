package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Slack []SlackConfig `json:"slack"`
}

type SlackConfig struct {
	Channel   string `json:"channel"`
	Username  string `json:"username"`
	IconEmoji string `json:"icon_emoji"`
	EndPoint  string `json:"end_point"`
}

func NewConfig() (Config, error) {
	var config Config

	str, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Could not read config.json. ", err)
		return config, err
	}

	if err := json.Unmarshal(str, &config); err != nil {
		fmt.Println("JSON Unmarshal Error:", err)
		return config, err
	}

	return config, nil
}
