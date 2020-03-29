package jimiko

import (
	"testing"
)

func TestParseText(t *testing.T) {
	e := EventData{
		Text: "test",
	}
	exp := "test"
	act := e.parseText()
	if act != exp {
		t.Fatalf("failed test %s", act)
	}

	e = EventData{
		Text: "\u003c@UEG9LPTND\u003e test2",
	}
	exp = "test2"
	act = e.parseText()
	if act != exp {
		t.Fatalf("failed test2 %s", act)
	}

	e = EventData{
		Text: "\u003c@UEG9LPTND test3",
	}
	exp = "\u003c@UEG9LPTND test3"
	act = e.parseText()
	if act != exp {
		t.Fatalf("failed test3 %s", act)
	}
}

func TestCreateMessage(t *testing.T) {
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
