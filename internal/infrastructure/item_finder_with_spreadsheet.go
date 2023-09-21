package infrastructure

import (
	"context"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"github.com/kimikimi714/jimiko/internal/domain"
	"github.com/kimikimi714/jimiko/internal/log"
)

// ItemFinderWithSpreadsheet searches items in shopping list spreadsheet.
type ItemFinderWithSpreadsheet struct {
	svr *sheets.Service
	id  string
}

// NewItemFinderWithSpreadsheet creates ItemFinderWithSpreadsheet instance.
func NewItemFinderWithSpreadsheet(spreadsheetID string) (*ItemFinderWithSpreadsheet, error) {
	ctx := context.Background()

	s, err := sheets.NewService(ctx, option.WithScopes(sheets.SpreadsheetsReadonlyScope))
	if err != nil {
		log.Error("Unable to retrieve Sheets client: %s", err)
		return nil, err
	}

	return &ItemFinderWithSpreadsheet{
		svr: s,
		id:  spreadsheetID,
	}, nil
}

// FindAll finds all items in shopping list from spreadsheet.
func (f *ItemFinderWithSpreadsheet) FindAll() ([]*domain.Item, error) {
	ss, err := f.svr.Spreadsheets.Get(f.id).Do()
	if err != nil {
		log.Error("Could not find spreadsheet %s", f.id)
		return nil, err
	}
	var is []*domain.Item
	for _, s := range ss.Sheets {
		resp, err := f.svr.Spreadsheets.Values.Get(f.id, s.Properties.Title+"!A:B").Do()
		if err != nil {
			log.Error("Could not read the sheet %v", s.Properties.Title)
			return nil, err
		}
		is = append(is, fetchAllItemsFrom(s.Properties.Title, resp)...)
	}
	return is, nil
}

func fetchAllItemsFrom(c string, r *sheets.ValueRange) []*domain.Item {
	var res []*domain.Item
	for i, row := range r.Values {
		if i == 0 {
			continue
		}
		n, ok := row[1].(string)
		if !ok {
			continue
		}
		switch row[0] {
		case "あり":
			res = append(res, &domain.Item{
				Category: domain.Category(c),
				Amount:   1,
				Name:     n,
			})
		case "なし":
			res = append(res, &domain.Item{
				Category: domain.Category(c),
				Amount:   0,
				Name:     n,
			})
		default:
			continue
		}

	}
	return res
}
