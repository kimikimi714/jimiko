package infrastructure

import (
	"github.com/kimikimi714/jimiko/domain"
)

type ItemFinder interface {
	FindAll() ([]*domain.Item, error)
}
