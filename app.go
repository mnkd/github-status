package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	slack "github.com/mnkd/slackposter"
)

// App is the application object.
type App struct {
	Good        bool
	SlackConfig slack.Config
}

func (app *App) fetchStatus() (GitHubStatus, error) {
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

// Run invoke the App.
func (app *App) Run() int {
	status, err := app.fetchStatus()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitCodeError
	}

	if app.Good == false && status.IsGood() {
		fmt.Fprintln(os.Stdout, "GitHub Site Status is Good.")
		return ExitCodeOK
	}

	payload := status.BuildPayload(app.SlackConfig)
	if len(payload.Attachments) == 0 {
		fmt.Fprintln(os.Stderr, "Could not build a message payload.")
		return ExitCodeError
	}

	poster := slack.NewSlack(app.SlackConfig)
	err = poster.PostPayload(payload)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitCodeError
	}

	return ExitCodeOK
}
