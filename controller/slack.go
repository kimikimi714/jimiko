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

type SlackController struct {
	r SlackRequestBody
}

func NewSlackController(r SlackRequestBody) *SlackController {
	return &SlackController{r: r}
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

func (c SlackController) Reply() error {
	text := c.r.Event.parseText()
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
	jsonStr, err := createSlackMessage(m)
	if err != nil {
		log.Fatalf("failed to create a message: %v", err)
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
