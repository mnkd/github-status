package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	// "io/ioutil"
	"net/http"
)

type Slack struct {
	Channel   string
	Username  string
	IconEmoji string
	EndPoint  string
}

// Refet to
// Attaching content and links to messages | Slack
// https://api.slack.com/docs/message-attachments

type Payload struct {
	Channel     string       `json:"channel"`
	Username    string       `json:"username"`
	Text        string       `json:"text"`
	IconEmoji   string       `json:"icon_emoji"`
	LinkNames   bool         `json:"link_names"`
	Attachments []Attachment `json:"attachments"`
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type Attachment struct {
	Fallback string  `json:"fallback"`
	Text     string  `json:"text"`
	Color    string  `json:"color"`
	Fields   []Field `json:"fields"`
}

func NewSlack(config SlackConfig) Slack {
	slack := Slack{
		Channel:   config.Channel,
		Username:  config.Username,
		IconEmoji: config.IconEmoji,
		EndPoint:  config.EndPoint,
	}
	return slack
}

func (slack *Slack) PostPayload(payload Payload) error {
	return slack.post(payload)
}

func (slack *Slack) PostMessage(message string) error {
	var payload Payload
	payload.Channel = slack.Channel
	payload.Username = slack.Username
	payload.IconEmoji = slack.IconEmoji
	payload.LinkNames = true
	payload.Text = message
	return slack.post(payload)
}

func (slack *Slack) post(payload Payload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Create request
	body := bytes.NewBuffer(data)
	request, err := http.NewRequest("POST", slack.EndPoint, body)
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		return err
	}

	// Invoke Request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	// Read Response Body
	// responseBody, _ := ioutil.ReadAll(response.Body)

	// Display Results
	fmt.Println("response Status : ", response.Status)
	// fmt.Println("response Headers : ", resp.Header)
	// fmt.Println("response Body : ", string(respBody))

	return nil
}
