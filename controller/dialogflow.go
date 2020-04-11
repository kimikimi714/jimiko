package controller

import (
	"log"
	"net/http"
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

type DialogflowController struct {
	r DialogflowRequestBody
	w http.ResponseWriter
}

func NewDialogflowController(r DialogflowRequestBody, w http.ResponseWriter) *DialogflowController {
	return &DialogflowController{r: r, w: w}
}

func (c *DialogflowController) Reply() error {
	exists := c.r.QueryResult.exists()
	ii, _ := usecase.NewItemInteractorWithSpreadsheet(os.Getenv("SPREADSHEET_ID"))
	ip := presenter.ItemPresenter{}
	jsonStr := ""
	if exists {
		m, _ := ip.ReadAllFullItems(ii)
		jsonStr = createDialogFlowMessage(m)
	} else {
		m, _ := ip.ReadAllLackedItems(ii)
		jsonStr = createDialogFlowMessage(m)
	}
	log.Print(jsonStr)
	_, err := c.w.Write([]byte(jsonStr))
	if err != nil {
		return err
	}
	c.w.Header().Set("Content-Type", "application/json")

	return nil
}

// parseText is prefixを除去してメッセージの本体だけを取り出す
func (e QueryResult) exists() bool {
	params := e.Parameters
	if params["exists"] == "ある" {
		return true
	}
	return false
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
