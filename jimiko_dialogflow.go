package jimiko

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kimikimi714/jimiko/controller"
)

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
