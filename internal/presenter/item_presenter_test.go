package presenter

import (
	"errors"
	"testing"

	"github.com/kimikimi714/jimiko/internal/domain"
)

func TestItemPresenter_ReadAllFullItems(t *testing.T) {
	var p ItemPresenter
	mock := dummyItemFilter{
		is: []*domain.Item{
			{
				Category: domain.Foods,
				Amount:   1,
				Name:     "food1",
			},
			{
				Category: domain.Foods,
				Amount:   1,
				Name:     "food2",
			},
			{
				Category: domain.HouseholdGoods,
				Amount:   1,
				Name:     "goods1",
			},
		},
	}
	exp := "food1、food2、goods1、があります。"
	act, _ := p.ReadAllFullItems(&mock)
	if act != exp {
		t.Fatalf("act is not expected format: %s", act)
	}

	mock = dummyItemFilter{
		err: errors.New("test"),
	}
	act, _ = p.ReadAllFullItems(&mock)
	if act != "" {
		t.Fatalf("act is not empty string: %s", act)
	}
}

func TestItemPresenter_ReadAllLackedItems(t *testing.T) {
	var p ItemPresenter
	mock := dummyItemFilter{
		is: []*domain.Item{
			{
				Category: domain.Foods,
				Amount:   0,
				Name:     "food1",
			},
			{
				Category: domain.Foods,
				Amount:   0,
				Name:     "food2",
			},
			{
				Category: domain.HouseholdGoods,
				Amount:   0,
				Name:     "goods1",
			},
		},
	}
	exp := "food1、food2、goods1、がありません。"
	act, _ := p.ReadAllLackedItems(&mock)
	if act != exp {
		t.Fatalf("act is not expected format: %s", act)
	}

	mock = dummyItemFilter{
		err: errors.New("test"),
	}
	act, _ = p.ReadAllLackedItems(&mock)
	if act != "" {
		t.Fatalf("act is not empty string: %s", act)
	}
}

func TestItemPresenter_ReadItemStatus(t *testing.T) {
	var p ItemPresenter
	mock := dummyItemFilter{
		is: []*domain.Item{
			{
				Category: domain.Foods,
				Amount:   1,
				Name:     "food1",
			},
			{
				Category: domain.Foods,
				Amount:   0,
				Name:     "food2",
			},
		},
	}
	exp := "food1はあるよ。"
	act, _ := p.ReadItemStatus("food1", &mock)
	if act != exp {
		t.Fatalf("act is not expected format: %s", act)
	}
	exp = "food2はないよ。"
	act, _ = p.ReadItemStatus("food2", &mock)
	if act != exp {
		t.Fatalf("act is not expected format: %s", act)
	}
	mock = dummyItemFilter{
		err: errors.New("test"),
	}
	act, _ = p.ReadItemStatus("", &mock)
	if act != "" {
		t.Fatalf("act is not empty string: %s", act)
	}
}

type dummyItemFilter struct {
	is  []*domain.Item
	err error
}

func (p *dummyItemFilter) PickUpLackedItems() ([]*domain.Item, error) {
	return p.is, p.err
}

func (p *dummyItemFilter) PickUpFullItems() ([]*domain.Item, error) {
	return p.is, p.err
}

func (p *dummyItemFilter) PickUpItem(name string) (*domain.Item, error) {
	for _, i := range p.is {
		if i.Name == name {
			return i, nil
		}
	}
	return nil, p.err
}
