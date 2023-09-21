package controller

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestCheckHeaders(t *testing.T) {
	tests := []struct {
		name       string
		signature  string
		timestamp  string
		errMessage string
	}{
		{
			name:       "headers missing",
			signature:  "",
			timestamp:  "",
			errMessage: "Required headers are missing.",
		},
		{
			name:       "signature header missing",
			signature:  "",
			timestamp:  "11111",
			errMessage: "Required headers are missing.",
		},
		{
			name:       "timestamp header missing",
			signature:  "aaa",
			timestamp:  "",
			errMessage: "Required headers are missing.",
		},
		{
			name:       "cannot parse timestamp",
			signature:  "aaa",
			timestamp:  "dummy",
			errMessage: "Cannot parse X-Slack-Request-Timestamp header.",
		},
		{
			name:       "expired 10 min ago",
			signature:  "aaa",
			timestamp:  strconv.FormatInt(time.Now().Unix()-60*10, 10),
			errMessage: "Expired timestamp.",
		},
		{
			name:       "nomal headers",
			signature:  "aaa",
			timestamp:  strconv.FormatInt(time.Now().Unix()+60, 10),
			errMessage: "none",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkHeaders(tt.timestamp, tt.signature)
			if err != nil && !strings.Contains(err.Error(), tt.errMessage) {
				t.Fatalf("SlackController.checkHeaders():%s failed. err: %v.", tt.name, err)
			} else if err == nil && tt.errMessage != "none" {
				t.Fatalf("SlackController.checkHeaders():%s failed.", tt.name)
			}
		})
	}
}

func TestCheckHMAC(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		secret     string
		timestamp  string
		signature  string
		errMessage string
	}{
		{
			name:       "cannot verify",
			body:       "hoge",
			secret:     "dummy",
			timestamp:  "111",
			signature:  "v0=a2114d57b48eac39b9ad189dd8316235a7b4a8d21a10bd27519666489c69b503",
			errMessage: "Cannot verify this request.",
		},
		{
			// see: https://api.slack.com/authentication/verifying-requests-from-slack
			name:       "example request",
			body:       "token=xyzz0WbapA4vBCDEFasx0q6G&team_id=T1DC2JH3J&team_domain=testteamnow&channel_id=G8PSS9T3V&channel_name=foobar&user_id=U2CERLKJA&user_name=roadrunner&command=%2Fwebhook-collect&text=&response_url=https%3A%2F%2Fhooks.slack.com%2Fcommands%2FT1DC2JH3J%2F397700885554%2F96rGlfmibIGlgcZRskXaIFfN&trigger_id=398738663015.47445629121.803a0bc887a14d10d2c447fce8b6703c",
			secret:     "8f742231b10e8888abcd99yyyzzz85a5",
			timestamp:  "1531420618",
			signature:  "v0=a2114d57b48eac39b9ad189dd8316235a7b4a8d21a10bd27519666489c69b503",
			errMessage: "none",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkHMAC(tt.body, tt.secret, tt.timestamp, tt.signature)
			if err != nil && !strings.Contains(err.Error(), tt.errMessage) {
				t.Errorf("SlackController.checkHMAC(): %s failed. error = %v.", tt.name, err)
			} else if err == nil && tt.errMessage != "none" {
				t.Errorf("SlackController.checkHMAC(): %s failed.", tt.name)
			}
		})
	}
}

func TestParseRequest(t *testing.T) {
	tests := []struct {
		name string
		json string
		want slackRequestBody
	}{
		{
			name: "only type",
			json: `{ "type": "test" }`,
			want: slackRequestBody{Type: "test"},
		},
		{
			name: "url_verification request", // see: https://api.slack.com/events/url_verification
			json: `{ "token": "token", "challenge": "XXXX", "type": "url_verification" }`,
			want: slackRequestBody{Type: "url_verification", Token: "token", Challenge: "XXXX"},
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
			want: slackRequestBody{Type: "event_callback", Token: "ZZZZZZWSxiZZZ2yIvs3peJ", Event: slackEventData{
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
			var got slackRequestBody
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
		args slackEventData
		want string
	}{
		{
			name: "1 word",
			args: slackEventData{Text: "test"},
			want: "test",
		},
		{
			name: "any words",
			args: slackEventData{Text: "test test"},
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
	e := slackEventData{
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
