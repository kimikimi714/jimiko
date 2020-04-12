package jimiko

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kimikimi714/jimiko/controller"
	)

// Slack is Slack向けep
func Slack(w http.ResponseWriter, r *http.Request) {
	var d controller.SlackRequestBody
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("failed to parse: %v", r.Body)
		return
	}

	// 地味子にメンション付きで話しかけないと反応しない
	if d.Event.Type != "app_mention" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c := controller.SlackController{}
	err := c.Reply(d)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Dialogflow is Dialogflow向けep
func Dialogflow(w http.ResponseWriter, r *http.Request) {
	var d controller.DialogflowRequestBody
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("failed to parse: %v", r.Body)
		return
	}

	c := controller.DialogflowController{}
	jsonStr, err := c.Reply(d)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte(jsonStr))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("failed to write response body with json: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
