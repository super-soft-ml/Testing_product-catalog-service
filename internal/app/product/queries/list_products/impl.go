package list_products

import (
	"context"

	"product-catalog-service/internal/app/product/contracts"
)

// QueryImpl implements Query using ReadModel.
type QueryImpl struct {
	ReadModel contracts.ReadModel
}

// Execute delegates to ReadModel.ListProducts.
func (q *QueryImpl) Execute(ctx context.Context, req ListProductsRequest) (*ListProductsResult, error) {
	return q.ReadModel.ListProducts(ctx, req)
}

var _ Query = (*QueryImpl)(nil)
