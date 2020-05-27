package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/kimikimi714/jimiko/presenter"
	"github.com/kimikimi714/jimiko/usecase"
)

// SlackRequestBody represents a request from Slack.
type SlackRequestBody struct {
	Type      string `json:"type"`
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	Event     EventData
}

// EventData represents event data from slack.
type EventData struct {
	Type           string `json:"type"`
	UserID         string `json:"user"`
	Text           string `json:"text"`
	Timestamp      string `json:"ts"`
	ChannelID      string `json:"channel"`
	EventTimestamp string `json:"event_ts"`
}

// SlackController represents interface which communicates with Slack.
type SlackController struct{}

// text extracts text excluding bot name.
func (e EventData) text() string {
	text := e.Text
	prefix := os.Getenv("SLACK_BOT_NAME")
	if strings.HasPrefix(text, prefix) {
		return strings.TrimPrefix(text, prefix)
	}
	return text
}

// Reply replies messages with enough / not enough shopping list to Slack.
func (c SlackController) Reply(r SlackRequestBody) error {
	text := r.Event.text()
	ii, _ := usecase.NewItemFilterWithSpreadsheet(os.Getenv("SPREADSHEET_ID"))
	ip := presenter.ItemPresenter{}
	var m string
	var err error
	switch text {
	case "何がある?":
		m, err = ip.ReadAllFullItems(ii)
	case "何がない?":
		m, err = ip.ReadAllLackedItems(ii)
	case "買い物リスト":
		m = "https://docs.google.com/spreadsheets/d/" + os.Getenv("SPREADSHEET_ID")
	default:
		log.Print("text: " + text)
		// FIXME 本当は text を直接 slack 表示させたい
		// text の中にメンションが含まれると無限ループに入ってしまうので
		// 今はログに出して slack には表示させないようにしている
		m = "何していいかわかりません。ログを見てください。"
	}
	if err != nil {
		log.Printf("failed to get items: %v", err)
		m = "買い物リストがうまく取得できませんでした"
	}

	jsonStr, err := createSlackMessage(m)
	if err != nil {
		log.Fatalf("failed to create a message: %v", err)
		return err
	}
	log.Print(jsonStr)
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
