package product

import (
	"context"

	"product-catalog-service/internal/app/product/domain"
)

// ProductRepo is the repository for product aggregate (writes).
// Returns mutations; does not apply them. Load is for command-side only.
type ProductRepo interface {
	Load(ctx context.Context, productID string) (*domain.Product, error)
	InsertMut(p *domain.Product) interface{}
	UpdateMut(p *domain.Product) interface{}
}
