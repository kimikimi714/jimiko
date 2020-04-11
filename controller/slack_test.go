package controller

import (
	"os"
	"testing"
)

func TestText(t *testing.T) {
	e := EventData{
		Text: "test",
	}
	exp := "test"
	act := e.text()
	if act != exp {
		t.Fatalf("failed test %s", act)
	}

	_ = os.Setenv("SLACK_BOT_NAME", "test bot ")
	defer os.Unsetenv("SLACK_BOT_NAME")
	e = EventData{
		Text: "test bot test2",
	}
	exp = "test2"
	act = e.text()
	if act != exp {
		t.Fatalf("failed test2 %s", act)
	}

	e = EventData{
		Text: "test test3",
	}
	exp = "test test3"
	act = e.text()
	if act != exp {
		t.Fatalf("failed test3 %s", act)
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
