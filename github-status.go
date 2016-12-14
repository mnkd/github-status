package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type App struct {
	Config Config
}

var (
	slackConfigIndex int
	good             bool
)

var app = App{}
var JST = time.FixedZone("JST", 3600*9)

func JSTDate(dateString string) time.Time {
	// 2016-10-04T13:40:42Z
	utc, _ := time.Parse(time.RFC3339, dateString)
	jst := utc.In(JST)
	return jst
}

type GitHubStatus struct {
	Status      string `json:"status"`
	LastUpdated string `json:"last_updated"`
}

func (github *GitHubStatus) payload() Payload {
	date := JSTDate(github.LastUpdated)
	dateString := date.Format("2006-01-02 15:04")
	var slackConfig SlackConfig = app.Config.Slack[slackConfigIndex]

	var payload Payload
	payload.Channel = slackConfig.Channel
	payload.Username = "GitHub Status | " + github.Status
	payload.IconEmoji = slackConfig.IconEmoji
	payload.LinkNames = true

	statusField := Field{
		Title: "Status",
		Value: github.Status,
		Short: true,
	}

	dateField := Field{
		Title: "Date",
		Value: dateString,
		Short: true,
	}

	var color = "good"
	var mention = ""
	switch github.Status {
	case "good":
		color = "good"
		break
	case "minor":
		color = "warning"
		break
	case "major":
		color = "danger"
		mention = "@hear "
		break
	default:
		return payload
	}

	attachment := Attachment{
		Fallback: mention + "GitHub Status: " + github.Status + " - https://status.github.com",
		Text:     mention + "<https://status.github.com/|GitHub Status> : " + github.Status,
		Color:    color,
		Fields:   []Field{statusField, dateField},
	}

	payload.Attachments = []Attachment{attachment}

	return payload
}

func (app *App) status() (GitHubStatus, error) {
	// Prepare HTTP Request

	// https://status.github.com/api
	// GET /api/status.json
	// Returns the current system status--one of good (green), minor (yellow), or major (red)--and timestamp.

	url := "https://status.github.com/api/status.json"
	request, err := http.NewRequest("GET", url, nil)

	// Fetch Request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failure : ", err)
	}

	// Read Response Body
	responseBody, _ := ioutil.ReadAll(response.Body)

	// Decode JSON
	var github GitHubStatus
	if err := json.Unmarshal(responseBody, &github); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return github, err
	}

	return github, nil
}

func (app *App) run() error {
	status, err := app.status()
	if err != nil {
		return err
	}
	if good == false && status.Status == "good" {
		return nil
	}

	payload := status.payload()

	if len(payload.Attachments) == 0 {
		return nil
	}

	slackConfig := app.Config.Slack[slackConfigIndex]
	slack := NewSlack(slackConfig)
	err = slack.PostPayload(payload)
	return err
}

func init() {
	flag.IntVar(&slackConfigIndex, "sci", 0, "Slack Config Index (default: 0)")
	flag.BoolVar(&good, "good", false, "Post Good Status (default: false)")
	flag.Parse()

	config, err := NewConfig()
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	app.Config = config
}

func main() {
	if err := app.run(); err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
}
