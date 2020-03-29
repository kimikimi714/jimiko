package jimiko

type DialogflowRequestBody struct {
	QueryResult QueryResult `json:"queryResult"`
}

type QueryResult struct {
	QueryText  string                 `json:"queryText"`
	Parameters map[string]interface{} `json:"parameters"`
}

// ReplyMention replies a message
func Reply(e QueryResult) (string, error) {
	exists := e.exists()
	food := CheckFood(exists)
	var jsonStr string
	var message string
	if exists {
		message = food + "はあるよ"
	} else {
		message = food + "はないよ"
	}
	jsonStr = createDialogFlowMessage(message)
	return jsonStr, nil
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
