package infrastructure

import (
	"context"
	"jimiko/domain"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type ItemFinderWithSpreadsheet struct {
	svr *sheets.Service
	id  string
}

func NewItemFinderWithSpreadsheet(spreadsheetId string) (*ItemFinderWithSpreadsheet, error) {
	ctx := context.Background()

	s, err := sheets.NewService(ctx, option.WithScopes(sheets.SpreadsheetsReadonlyScope))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
		return nil, err
	}

	return &ItemFinderWithSpreadsheet{
		svr: s,
		id:  spreadsheetId,
	}, nil
}

func (f *ItemFinderWithSpreadsheet) FindAll() ([]*domain.Item, error) {
	ss, err := f.svr.Spreadsheets.Get(f.id).Do()
	if err != nil {
		log.Fatalf("Could not find spreadsheet %v", f.id)
		return nil, err
	}
	var is []*domain.Item
	for _, s := range ss.Sheets {
		resp, err := f.svr.Spreadsheets.Values.Get(f.id, s.Properties.Title+"!A:B").Do()
		if err != nil {
			log.Fatalf("Could not read the sheet %v", s.Properties.Title)
			return nil, err
		}
		is = append(is, fetchAllItemsFrom(s.Properties.Title, resp)...)
	}
	return is, nil
}

func fetchAllItemsFrom(c string, r *sheets.ValueRange) []*domain.Item {
	res := []*domain.Item{}
	for i, row := range r.Values {
		if i == 0 {
			continue
		}
		a := -1
		switch row[0] {
		case "あり":
			a = 1
		case "なし":
			a = 0
		default:
			continue
		}
		n, ok := row[1].(string)
		if !ok {
			continue
		}
		res = append(res, &domain.Item{
			Category: domain.Category(c),
			Amount:   a,
			Name:     n,
		})
	}
	return res
}
