package list_products

import "math/big"

// ProductDTO is a list item.
type ProductDTO struct {
	ProductID       string
	Name           string
	Description    string
	Category       string
	BasePrice      *big.Rat
	EffectivePrice *big.Rat
	Status         string
}

// ListProductsRequest is the list query request.
type ListProductsRequest struct {
	Category string
	Status   string // active, draft, etc.; empty = active
	Limit    int32
	Offset   int32
}

// ListProductsResult is the paginated result.
type ListProductsResult struct {
	Products []*ProductDTO
	Total    int32
}
