package jimiko

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type SlackRequestBody struct {
	Type      string `json:"type"`
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	Event     EventData
}

// EventData is slackから飛んでくるeventを表すEntity
type EventData struct {
	Type           string `json:"type"`
	UserID         string `json:"user"`
	Text           string `json:"text"`
	Timestamp      string `json:"ts"`
	ChannelID      string `json:"channel"`
	EventTimestamp string `json:"event_ts"`
}

// parseText is prefixを除去してメッセージの本体だけを取り出す
func (e EventData) parseText() string {
	text := e.Text
	prefix := os.Getenv("SLACK_BOT_NAME")
	if strings.HasPrefix(text, prefix) {
		return strings.TrimPrefix(text, prefix)
	}
	return text
}

// ReplyMention replies a message
func ReplyMention(e EventData) error {
	text := e.parseText()
	var jsonStr string
	var err error
	if text == "check food" {
		CheckFood(false)
		jsonStr, err = createSlackMessage("log を見てください")
	} else {
		jsonStr, err = createSlackMessage(text)
	}
	if err != nil {
		log.Fatalf("failed to create a message: %v", err)
	}
	err = postMessage(jsonStr)
	if err != nil {
		log.Fatalf("failed to post a message to slack: %v", err)
		return err
	}
	return nil
}

// createSlackMessage creates a message to post to slack
func createSlackMessage(s string) (string, error) {
	message := map[string]interface{}{
		"text": s,
	}

	jsonByte, err := json.Marshal(message)
	if err != nil {
		return "", err
	}
	return string(jsonByte), nil
}

// postMessage posts a message to slack
func postMessage(jsonStr string) (err error) {
	reader := strings.NewReader(jsonStr)
	_, err = http.Post(os.Getenv("SLACK_INCOMING_WEBHOOK"), "application/json", reader)

	if err != nil {
		return err
	}

	return nil
}
