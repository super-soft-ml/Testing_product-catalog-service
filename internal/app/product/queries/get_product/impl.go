package get_product

import (
	"context"

	"product-catalog-service/internal/app/product/contracts"
)

// QueryImpl implements Query using ReadModel.
type QueryImpl struct {
	ReadModel contracts.ReadModel
}

// Execute delegates to ReadModel.GetProduct.
func (q *QueryImpl) Execute(ctx context.Context, productID string) (*ProductDTO, error) {
	return q.ReadModel.GetProduct(ctx, productID)
}

var _ Query = (*QueryImpl)(nil)
