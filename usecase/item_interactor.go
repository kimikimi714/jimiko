package usecase

import (
	"github.com/kimikimi714/jimiko/domain"
)

type ItemInteractor interface {
	PickUpLackedItems() ([]*domain.Item, error)
	PickUpFullItems() ([]*domain.Item, error)
}
