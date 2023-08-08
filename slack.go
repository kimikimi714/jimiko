package jimiko

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/kimikimi714/jimiko/controller"
)

// Slack is Slack向けep
func Slack(w http.ResponseWriter, r *http.Request) {
	c := controller.SlackController{}
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	secret := os.Getenv("SLACK_SIGINING_SECRET")
	err = c.Verify(r.Header, string(body), secret)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatalf("err is: %s", err)
	}

	var d controller.SlackRequestBody
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("failed to parse: %v", r.Body)
	}

	// 地味子にメンション付きで話しかけないと反応しない
	if d.Event.Type != "app_mention" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = c.Reply(d)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
