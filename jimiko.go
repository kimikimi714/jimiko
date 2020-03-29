package jimiko

import (
	"encoding/json"
	"log"
	"net/http"
)

// Dialog is Dialogflowからのリクエストを受け取るエンドポイント
func Dialog(w http.ResponseWriter, r *http.Request) {
	log.Print(r.Body)
	var d DialogflowRequestBody
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("failed to parse: %v", r.Body)
		return
	}
	log.Print(d)

	str, err := Reply(d.QueryResult)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	log.Print(str)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(str))
	w.WriteHeader(http.StatusOK)
}
