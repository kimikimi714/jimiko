package infrastructure

import (
	"github.com/kimikimi714/jimiko/domain"
)

// ItemFinder searches items in shopping list
type ItemFinder interface {
	FindAll() ([]*domain.Item, error)
}
