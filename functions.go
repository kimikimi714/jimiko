package jimiko

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/kimikimi714/jimiko/internal/controller"
	"github.com/kimikimi714/jimiko/internal/log"
)

func init() {
	functions.HTTP("Slack", slack)
}

// slack is slack向けep
func slack(w http.ResponseWriter, r *http.Request) {
	c := controller.SlackController{}
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Error("Cannot read request body. err is %s.", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	secret := os.Getenv("SLACK_SIGINING_SECRET")
	if err := c.Verify(r.Header, string(body), secret); err != nil {
		log.Error("SlackController.Verify got error: %s.", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var d controller.SlackRequestBody
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

	if err := c.Reply(d); err != nil {
		log.Error("SlackController.Reply got error: %s.", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
