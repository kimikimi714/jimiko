// Package jimiko is jimiko!
package jimiko

import (
	"io"
	"net/http"

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
	c.Response(r, body, w)
}
