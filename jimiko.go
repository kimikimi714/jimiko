package jimiko

import (
	"encoding/json"
	"log"
	"net/http"

	"jimiko/controller"
)

// Slack is Slack向けep
func Slack(w http.ResponseWriter, r *http.Request) {
	var d controller.SlackRequestBody
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("failed to parse: %v", r.Body)
		return
	}

	c := controller.NewSlackController(d)
	// 地味子にメンション付きで話しかけた場合
	if d.Event.Type == "app_mention" {
		err := c.Reply()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}
