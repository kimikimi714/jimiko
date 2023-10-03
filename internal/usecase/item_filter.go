// Package usecase implements functions to filter and retrieve items based on specific use cases.
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
