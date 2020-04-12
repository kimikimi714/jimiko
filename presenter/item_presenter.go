package presenter

import (
	"github.com/kimikimi714/jimiko/domain"
	"github.com/kimikimi714/jimiko/usecase"
)

type ItemPresenter struct {
}

func (p ItemPresenter) ReadAllLackedItems(interactor usecase.ItemInteractor) (string, error) {
	is, err := interactor.PickUpLackedItems()
	if err != nil {
		return "", err
	}
	return concatAllItems(is) + "がありません。", nil
}

func (p ItemPresenter) ReadAllFullItems(interactor usecase.ItemInteractor) (string, error) {
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
