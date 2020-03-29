package jimiko

import (
	"fmt"
	"google.golang.org/api/option"
	"log"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/api/sheets/v4"
)

// CheckFood is 食品を取り出す
func CheckFood(exists bool) string {
	ctx := context.Background()

	srv, err := sheets.NewService(ctx, option.WithScopes(sheets.SpreadsheetsReadonlyScope))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
		return ""
	}

	// 参照したいスプレッドシートのIDを環境変数からとってくる
	spreadsheetID := os.Getenv("SPREADSHEET_ID")
	readRange := "食品!A2:B"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		log.Println("No data found.")
		return ""
	} else {
		log.Println("checking")
		for _, row := range resp.Values {
			if exists == true && row[0] == "ある" {
				log.Printf("checked: %s", row[1])
				return fmt.Sprintf("%s", row[1])
			} else if exists == false && row[0] == "なし" {
				log.Printf("checked: %s", row[1])
				return fmt.Sprintf("%s", row[1])
			}
		}
	}
	return ""
}
