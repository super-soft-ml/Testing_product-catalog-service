package product

import (
	"context"

	"product-catalog-service/internal/app/product/queries/get_product"
	"product-catalog-service/internal/app/product/queries/list_products"
)

// ReadModel provides query-side access (get by ID, list with filters).
type ReadModel interface {
	GetProduct(ctx context.Context, productID string) (*get_product.ProductDTO, error)
	ListProducts(ctx context.Context, req list_products.ListProductsRequest) (*list_products.ListProductsResult, error)
}
