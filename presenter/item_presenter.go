package presenter

import (
	"github.com/kimikimi714/jimiko/domain"
	"github.com/kimikimi714/jimiko/usecase"
)

// ItemPresenter creates some sentences which take into account
// item statuses in shopping list.
type ItemPresenter struct {}

// ReadAllLackedItems processes into a phrase that means what items are not enough
// in shopping list.
func (p ItemPresenter) ReadAllLackedItems(interactor usecase.ItemFilter) (string, error) {
	is, err := interactor.PickUpLackedItems()
	if err != nil {
		return "", err
	}
	return concatAllItems(is) + "がありません。", nil
}

// ReadAllFullItems processes into a phrase that means what items are enough
// in shopping list.
func (p ItemPresenter) ReadAllFullItems(interactor usecase.ItemFilter) (string, error) {
	is, err := interactor.PickUpFullItems()
	if err != nil {
		return "", err
	}
	return concatAllItems(is) + "があります。", nil
}

func concatAllItems(is []*domain.Item) string {
	res := ""
	for _, item := range is {
		res += item.Name + "、"
	}
	return res
}
