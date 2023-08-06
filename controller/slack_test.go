package controller

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func newHeader(sign, timestamp string) http.Header {
	h := http.Header{}
	h.Set("X-Slack-Signature", sign)
	h.Set("X-Slack-Request-Timestamp", timestamp)
	return h
}

func TestVerify(t *testing.T) {
	tests := []struct {
		name       string
		signature  string
		timestamp  string
		secret     string
		errMessage string
	}{
		{
			name:       "headers missing",
			signature:  "",
			timestamp:  "",
			secret:     "",
			errMessage: "Required headers are missing.",
		},
		{
			name:       "signature header missing",
			signature:  "",
			timestamp:  "11111",
			secret:     "",
			errMessage: "Required headers are missing.",
		},
		{
			name:       "timestamp header missing",
			signature:  "aaa",
			timestamp:  "",
			secret:     "",
			errMessage: "Required headers are missing.",
		},
		{
			name:       "secret is empty",
			signature:  "aaa",
			timestamp:  "dummy",
			secret:     "",
			errMessage: "SLACK_SIGINING_SECRET is empty.",
		},
		{
			name:       "expired timestamp",
			signature:  "aaa",
			timestamp:  "111",
			secret:     "dummy",
			errMessage: "Expired timestamp.",
		},
		{
			name:       "cannot parse timestamp",
			signature:  "aaa",
			timestamp:  "dummy",
			secret:     "dummy",
			errMessage: "Cannot parse X-Slack-Request-Timestamp header.",
		},
		{
			name:       "expired timestamp",
			signature:  "aaa",
			timestamp:  strconv.FormatInt(time.Unix(time.Now().Unix()+60*3, 0).Unix(), 10),
			secret:     "dummy",
			errMessage: "Expired timestamp.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := SlackController{}
			err := c.Verify(newHeader(tt.signature, tt.timestamp), []byte{}, tt.secret)
			if err != nil && !strings.Contains(err.Error(), tt.errMessage) {
				t.Fatalf("failed test: %+v", err)
			}
		})
	}
}

func TestParseRequest(t *testing.T) {
	tests := []struct {
		name string
		json string
		want SlackRequestBody
	}{
		{
			name: "only type",
			json: `{ "type": "test" }`,
			want: SlackRequestBody{Type: "test"},
		},
		{
			name: "url_verification request", // see: https://api.slack.com/events/url_verification
			json: `{ "token": "token", "challenge": "XXXX", "type": "url_verification" }`,
			want: SlackRequestBody{Type: "url_verification", Token: "token", Challenge: "XXXX"},
		},
		{
			name: "app_mention request", // see: https://api.slack.com/events/app_mention#app_mention-event__example-event-payloads__standard-app-mention-when-your-app-is-already-in-channel
			json: `{
				"token": "ZZZZZZWSxiZZZ2yIvs3peJ",
				"team_id": "T123ABC456",
				"api_app_id": "A123ABC456",
				"event": {
					"type": "app_mention",
					"user": "U123ABC456",
					"text": "What is the hour of the pearl, <@U0LAN0Z89>?",
					"ts": "1515449522.000016",
					"channel": "C123ABC456",
					"event_ts": "1515449522000016"
				},
				"type": "event_callback",
				"event_id": "Ev123ABC456",
				"event_time": 1515449522000016,
				"authed_users": [
					"U0LAN0Z89"
				]
			}`,
			want: SlackRequestBody{Type: "event_callback", Token: "ZZZZZZWSxiZZZ2yIvs3peJ", Event: EventData{
				Type:           "app_mention",
				UserID:         "U123ABC456",
				Text:           "What is the hour of the pearl, <@U0LAN0Z89>?",
				Timestamp:      "1515449522.000016",
				ChannelID:      "C123ABC456",
				EventTimestamp: "1515449522000016",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got SlackRequestBody
			if err := json.Unmarshal([]byte(tt.json), &got); err != nil {
				t.Fatalf("failed test: %+v", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Decorded struct mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestText(t *testing.T) {
	tests := []struct {
		name string
		args EventData
		want string
	}{
		{
			name: "1 word",
			args: EventData{Text: "test"},
			want: "test",
		},
		{
			name: "any words",
			args: EventData{Text: "test test"},
			want: "test test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.text(); got != tt.want {
				t.Errorf("text() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_RemoveBotName(t *testing.T) {
	_ = os.Setenv("SLACK_BOT_NAME", "test bot ")
	defer os.Unsetenv("SLACK_BOT_NAME")
	e := EventData{
		Text: "test bot test2",
	}
	want := "test2"
	got := e.text()
	if got != want {
		t.Errorf("text() = %v, want %v", got, want)
	}
}

func TestCreateSlackMessage(t *testing.T) {
	org := "test"
	exp := `{"text":"test"}`
	act, err := createSlackMessage(org)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	if act != exp {
		t.Fatalf("failed test %s", act)
	}
}
