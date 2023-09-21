package controller

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kimikimi714/jimiko/internal/log"
	"github.com/kimikimi714/jimiko/internal/presenter"
	"github.com/kimikimi714/jimiko/internal/usecase"
)

// slackRequestBody represents a request from Slack.
type slackRequestBody struct {
	Type      string         `json:"type"`
	Token     string         `json:"token,omitempty"`
	Challenge string         `json:"challenge,omitempty"`
	Event     slackEventData `json:"event,omitempty"`
}

// slackEventData represents event data from slack.
type slackEventData struct {
	Type           string `json:"type"`
	UserID         string `json:"user"`
	Text           string `json:"text"`
	Timestamp      string `json:"ts"`
	ChannelID      string `json:"channel"`
	EventTimestamp string `json:"event_ts"`
}

// SlackController represents interface which communicates with Slack.
type SlackController struct{}

func (c SlackController) Response(r *http.Request, body []byte, w http.ResponseWriter) {
	secret := os.Getenv("SLACK_SIGINING_SECRET")
	if err := c.verify(r.Header, string(body), secret); err != nil {
		log.Error("SlackController.Verify got error: %s.", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var d slackRequestBody
	if err := json.Unmarshal(body, &d); err != nil {
		log.Error("Failed to parse request body: %s.", string(body))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if d.Type == "url_verification" {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, d.Challenge)
		return
	}

	if d.Type != "event_callback" || d.Event.Type != "app_mention" {
		log.Warn("Not accepted event type:  %s / %s.", d.Type, d.Event.Type)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := c.reply(d); err != nil {
		log.Error("SlackController.Reply got error: %s.", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

// removeText extracts removeText excluding bot name.
func (e slackEventData) removeText(remove string) string {
	original := e.Text
	return strings.ReplaceAll(original, remove, "")
}

func (c SlackController) verify(headers http.Header, body, secret string) error {
	timestamp := headers.Get("X-Slack-Request-Timestamp")
	signature := headers.Get("X-Slack-Signature")
	if err := checkHeaders(timestamp, signature); err != nil {
		return err
	}
	if secret == "" {
		return fmt.Errorf("SLACK_SIGINING_SECRET is empty.")
	}
	if err := checkHMAC(body, secret, timestamp, signature); err != nil {
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
		log.Warn("sig: %v, calc: %v", []byte(signature[3:]), expectedMAC)
		return fmt.Errorf("Cannot verify this request.")
	}
	return nil
}

// reply replies messages with enough / not enough shopping list to Slack.
func (c SlackController) reply(r slackRequestBody) error {
	remove := os.Getenv("SLACK_BOT_NAME")
	text := r.Event.removeText(remove)
	ii, _ := usecase.NewItemFilterWithSpreadsheet(os.Getenv("SPREADSHEET_ID"))
	ip := presenter.ItemPresenter{}
	var m string
	var err error
	switch {
	case strings.Contains(text, "ある"):
		m, err = ip.ReadAllFullItems(ii)
	case strings.Contains(text, "ない"):
		m, err = ip.ReadAllLackedItems(ii)
	case strings.Contains(text, "リスト"):
		m = "https://docs.google.com/spreadsheets/d/" + os.Getenv("SPREADSHEET_ID")
	default:
		log.Warn("text: " + text)
		m = "「" + text + "」だと何していいかわかりませんでした :cry:"
	}
	if err != nil {
		log.Warn("failed to get items: %v", err)
		m = "買い物リストがうまく取得できませんでした"
	}

	jsonStr, err := createSlackMessage(m)
	if err != nil {
		log.Error("failed to create a message: %v", err)
		return err
	}
	log.Info(jsonStr)
	err = postMessage(jsonStr)
	if err != nil {
		log.Error("failed to post a message to slack: %v", err)
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
