package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"jimiko/presenter"
	"jimiko/usecase"
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

type SlackController struct {}

// text is prefixを除去してメッセージの本体だけを取り出す
func (e EventData) text() string {
	text := e.Text
	prefix := os.Getenv("SLACK_BOT_NAME")
	if strings.HasPrefix(text, prefix) {
		return strings.TrimPrefix(text, prefix)
	}
	return text
}

// Reply is slack bot にリクエストに応じて返信をさせる
func (c SlackController) Reply(r SlackRequestBody) error {
	text := r.Event.text()
	ii, _ := usecase.NewItemInteractorWithSpreadsheet(os.Getenv("SPREADSHEET_ID"))
	ip := presenter.ItemPresenter{}
	var m string
	var err error
	switch text {
	case "何がある?":
		m, err = ip.ReadAllFullItems(ii)
	case "何がない?" :
		m, err = ip.ReadAllLackedItems(ii)
	default:
		log.Print("text: " + text)
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
