package get_product

import "context"

// Query is the get-product query.
type Query interface {
	Execute(ctx context.Context, productID string) (*ProductDTO, error)
}
