package jimiko

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("Slack", Slack)
	functions.HTTP("Dialogflow", Dialogflow)
}
