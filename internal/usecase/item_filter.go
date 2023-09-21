package usecase

import (
	"github.com/kimikimi714/jimiko/internal/domain"
)

// ItemFilter filters items from shopping list
type ItemFilter interface {
	PickUpLackedItems() ([]*domain.Item, error)
	PickUpFullItems() ([]*domain.Item, error)
	PickUpItem(name string) (*domain.Item, error)
}
