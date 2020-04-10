package usecase

import (
	"jimiko/domain"
)

type ItemInteractor interface {
	PickUpLackedItems() ([]*domain.Item, error)
	PickUpFullItems() ([]*domain.Item, error)
}
