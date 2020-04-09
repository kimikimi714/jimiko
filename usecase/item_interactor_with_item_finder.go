package usecase

import (
	"jimiko/domain"
	"jimiko/infrastructure"
)

type ItemInteractorWithItemFinder struct {
	itemFinder infrastructure.ItemFinder
}

func NewItemInteractorWithSpreadsheet(spreadsheetId string) (*ItemInteractorWithItemFinder, error) {
	f, err := infrastructure.NewItemFinderWithSpreadsheet(spreadsheetId)
	if err != nil {
		return nil, err
	}
	return &ItemInteractorWithItemFinder{
		itemFinder: f,
	}, nil
}

func (p *ItemInteractorWithItemFinder) PickUpLackedItems() ([]*domain.Item, error) {
	is, err := p.itemFinder.FindAll()
	if err != nil {
		return nil, err
	}
	res := []*domain.Item{}
	for _, item := range is {
		if item.Amount == 0 {
			res = append(res, item)
		}
	}
	return res, nil
}

func (p *ItemInteractorWithItemFinder) PickUpFullItems() ([]*domain.Item, error) {
	is, err := p.itemFinder.FindAll()
	if err != nil {
		return nil, err
	}
	res := []*domain.Item{}
	for _, item := range is {
		if item.Amount == 1 {
			res = append(res, item)
		}
	}
	return res, nil
}
