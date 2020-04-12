package usecase

import (
	"github.com/kimikimi714/jimiko/domain"
)

// ItemFilter filters items from shopping list
type ItemFilter interface {
	PickUpLackedItems() ([]*domain.Item, error)
	PickUpFullItems() ([]*domain.Item, error)
}
