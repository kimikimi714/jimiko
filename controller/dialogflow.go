package controller

import (
	"log"
	"os"

	"jimiko/presenter"
	"jimiko/usecase"
)

type DialogflowRequestBody struct {
	QueryResult QueryResult `json:"queryResult"`
}

type QueryResult struct {
	QueryText  string                 `json:"queryText"`
	Parameters map[string]interface{} `json:"parameters"`
}

type DialogflowController struct {}

func (c *DialogflowController) Reply(r DialogflowRequestBody) (jsonStr string, err error) {
	exists := r.QueryResult.exists()
	ii, _ := usecase.NewItemInteractorWithSpreadsheet(os.Getenv("SPREADSHEET_ID"))
	ip := presenter.ItemPresenter{}
	if exists {
		m, _ := ip.ReadAllFullItems(ii)
		jsonStr = createDialogFlowMessage(m)
	} else {
		m, _ := ip.ReadAllLackedItems(ii)
		jsonStr = createDialogFlowMessage(m)
	}
	log.Print(jsonStr)
	return jsonStr, nil
}

func (e QueryResult) exists() bool {
	params := e.Parameters
	return params["exists"] == "ある"
}

// createDialogFlowMessage creates a message to post to slack
func createDialogFlowMessage(s string) string {
	str := `{
  "payload": {
    "google": {
      "expectUserResponse": true,
      "richResponse": {
        "items": [
          {
            "simpleResponse": {
              "textToSpeech": "` + s + `",
              "displayText": "` + s + `"
            }
          }
        ]
      }
    }
  }
}`

	return str
}
