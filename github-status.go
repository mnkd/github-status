package main

import (
	"time"

	slack "github.com/mnkd/slackposter"
)

func jstDate(dateString string) time.Time {
	// 2016-10-04T13:40:42Z
	utc, _ := time.Parse(time.RFC3339, dateString)
	jst := utc.In(time.FixedZone("JST", 3600*9))
	return jst
}

// GitHubStatus is represented GitHib Site Status.
type GitHubStatus struct {
	Status      string `json:"status"`
	LastUpdated string `json:"last_updated"`
}

// IsGood returns whether or not the github status is good.
func (github *GitHubStatus) IsGood() bool {
	return github.Status == "good"
}

// BuildPayload returns a slack message payload from github status.
func (github *GitHubStatus) BuildPayload(config slack.Config) slack.Payload {
	date := jstDate(github.LastUpdated)
	dateString := date.Format("2006-01-02 15:04")

	var payload slack.Payload
	payload.Channel = config.Channel
	payload.Username = "GitHub Status | " + github.Status
	payload.IconEmoji = config.IconEmoji
	payload.LinkNames = true

	statusField := slack.Field{
		Title: "Status",
		Value: github.Status,
		Short: true,
	}

	dateField := slack.Field{
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

	attachment := slack.Attachment{
		Fallback: mention + "GitHub Status: " + github.Status + " - https://status.github.com",
		Text:     mention + "<https://status.github.com/|GitHub Status> : " + github.Status,
		Color:    color,
		Fields:   []slack.Field{statusField, dateField},
	}

	payload.Attachments = []slack.Attachment{attachment}

	return payload
}
