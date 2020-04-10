package infrastructure

import (
	"jimiko/domain"
)

type ItemFinder interface {
	FindAll() ([]*domain.Item, error)
}
