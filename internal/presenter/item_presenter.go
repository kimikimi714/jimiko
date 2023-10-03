// Package presenter implements functions to add specific information to retrieved items and return strings.
package presenter

import (
	"github.com/kimikimi714/jimiko/internal/domain"
	"github.com/kimikimi714/jimiko/internal/usecase"
)

// ItemPresenter creates some sentences which take into account
// item statuses in shopping list.
type ItemPresenter struct{}

// ReadAllLackedItems processes into a phrase that means what items are not enough
// in shopping list.
func (p ItemPresenter) ReadAllLackedItems(filter usecase.ItemFilter) (string, error) {
	is, err := filter.PickUpLackedItems()
	if err != nil {
		return "", err
	}
	return concatAllItems(is) + "がありません。", nil
}

// ReadAllFullItems processes into a phrase that means what items are enough
// in shopping list.
func (p ItemPresenter) ReadAllFullItems(filter usecase.ItemFilter) (string, error) {
	is, err := filter.PickUpFullItems()
	if err != nil {
		return "", err
	}
	return concatAllItems(is) + "があります。", nil
}

// ReadItemStatus processes into a phrase that means we have the item or not.
func (p ItemPresenter) ReadItemStatus(name string, filter usecase.ItemFilter) (string, error) {
	i, err := filter.PickUpItem(name)
	if err != nil {
		return "", err
	}
	if i.Amount > 0 {
		return i.Name + "はあるよ。", nil
	}
	return i.Name + "はないよ。", nil
}

func concatAllItems(is []*domain.Item) string {
	res := ""
	for _, item := range is {
		res += item.Name + "、"
	}
	return res
}
