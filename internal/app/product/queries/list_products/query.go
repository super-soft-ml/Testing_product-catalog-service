package list_products

import "context"

// Query is the list-products query.
type Query interface {
	Execute(ctx context.Context, req ListProductsRequest) (*ListProductsResult, error)
}
