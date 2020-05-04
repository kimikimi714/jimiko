package controller

import (
	"log"
	"os"

	"github.com/kimikimi714/jimiko/presenter"
	"github.com/kimikimi714/jimiko/usecase"
)

// DialogflowRequestBody represents a request from Dialogflow.
type DialogflowRequestBody struct {
	QueryResult QueryResult `json:"queryResult"`
}

// QueryResult includes data that Dialogflow analyzed a phase and
// split into words.
type QueryResult struct {
	QueryText  string                 `json:"queryText"`
	Parameters map[string]interface{} `json:"parameters"`
}

// DialogflowController represents interface which communicates with Dialogflow.
type DialogflowController struct{}

// Reply replies messages with enough / not enough shopping list to Dialogflow.
func (c *DialogflowController) Reply(r DialogflowRequestBody) (jsonStr string, err error) {
	exists := r.QueryResult.exists()
	ii, _ := usecase.NewItemFilterWithSpreadsheet(os.Getenv("SPREADSHEET_ID"))
	ip := presenter.ItemPresenter{}
	m := ""
	name := r.QueryResult.getItemName()
	if name != "" {
		m, err = ip.ReadItemStatus(name, ii)
	} else if exists {
		m, err = ip.ReadAllFullItems(ii)
	} else {
		m, err = ip.ReadAllLackedItems(ii)
	}
	if err != nil {
		log.Printf("failed to get items: %v", err)
		m = "買い物リストがうまく取得できませんでした"
	}
	jsonStr = createDialogFlowMessage(m)
	log.Print(jsonStr)
	return jsonStr, nil
}

func (e QueryResult) exists() bool {
	params := e.Parameters
	return params["exists"] == "ある"
}

func (e QueryResult) getItemName() string {
	params := e.Parameters
	if params["item"] != nil {
		return params["item"].(string)
	}
	return ""
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
