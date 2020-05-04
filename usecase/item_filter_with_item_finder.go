package usecase

import (
	"github.com/kimikimi714/jimiko/domain"
	"github.com/kimikimi714/jimiko/infrastructure"
)

// ItemFilterWithItemFinder filters items ItemFinder found.
type ItemFilterWithItemFinder struct {
	itemFinder infrastructure.ItemFinder
}

// NewItemFilterWithSpreadsheet creates ItemFilterWithItemFinder instance.
func NewItemFilterWithSpreadsheet(spreadsheetID string) (*ItemFilterWithItemFinder, error) {
	f, err := infrastructure.NewItemFinderWithSpreadsheet(spreadsheetID)
	if err != nil {
		return nil, err
	}
	return &ItemFilterWithItemFinder{
		itemFinder: f,
	}, nil
}

// PickUpLackedItems picks up not enough items from shopping list.
func (p *ItemFilterWithItemFinder) PickUpLackedItems() ([]*domain.Item, error) {
	is, err := p.itemFinder.FindAll()
	if err != nil {
		return nil, err
	}
	var res []*domain.Item
	for _, i := range is {
		if i.Amount == 0 {
			res = append(res, i)
		}
	}
	return res, nil
}

// PickUpFullItems picks up enough items from shopping list.
func (p *ItemFilterWithItemFinder) PickUpFullItems() ([]*domain.Item, error) {
	is, err := p.itemFinder.FindAll()
	if err != nil {
		return nil, err
	}
	var res []*domain.Item
	for _, i := range is {
		if i.Amount > 0 {
			res = append(res, i)
		}
	}
	return res, nil
}

// PickUpItem picks up an item by item name.
func (p *ItemFilterWithItemFinder) PickUpItem(name string) (*domain.Item, error) {
	is, err := p.itemFinder.FindAll()
	if err != nil {
		return nil, err
	}
	for _, i := range is {
		if i.Name == name {
			return i, nil
		}
	}
	return nil, nil
}
