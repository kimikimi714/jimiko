package infrastructure

import (
	"github.com/kimikimi714/jimiko/internal/domain"
)

// ItemFinder searches items in shopping list
type ItemFinder interface {
	FindAll() ([]*domain.Item, error)
}
