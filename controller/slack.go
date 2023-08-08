package controller

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kimikimi714/jimiko/presenter"
	"github.com/kimikimi714/jimiko/usecase"
)

// SlackRequestBody represents a request from Slack.
type SlackRequestBody struct {
	Type      string    `json:"type"`
	Token     string    `json:"token,ommitempty"`
	Challenge string    `json:"challenge,ommitempty"`
	Event     EventData `json:"event,ommitempty"`
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

func (c SlackController) Verify(headers http.Header, body, secret string) error {
	timestamp := headers.Get("X-Slack-Request-Timestamp")
	signature := headers.Get("X-Slack-Signature")
	err := checkHeaders(timestamp, signature)
	if err != nil {
		return err
	}
	if secret == "" {
		return fmt.Errorf("SLACK_SIGINING_SECRET is empty.")
	}
	err = checkHMAC(body, secret, timestamp, signature)
	if err != nil {
		return err
	}
	return nil
}

func checkHeaders(timestamp string, signature string) error {
	if timestamp == "" || signature == "" {
		return fmt.Errorf("Required headers are missing.")
	}
	sec, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return fmt.Errorf("Cannot parse X-Slack-Request-Timestamp header. Error: %s", err)
	}

	if time.Now().Unix()-sec > 60*5 {
		return fmt.Errorf("Expired timestamp.")
	}
	return nil
}

func checkHMAC(body, secret, timestamp, signature string) error {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte("v0:" + timestamp + ":" + body))
	expectedMAC := mac.Sum(nil)
	bsignature, err := hex.DecodeString(strings.TrimPrefix(signature, "v0="))
	if err != nil {
		return err
	}

	if !hmac.Equal(bsignature, expectedMAC) {
		log.Printf("sig: %v, calc: %v", []byte(signature[3:]), expectedMAC)
		return fmt.Errorf("Cannot verify this request.")
	}
	return nil
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
