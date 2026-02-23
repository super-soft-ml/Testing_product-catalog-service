package get_product

import "math/big"

// ProductDTO is the read model for a single product with effective price.
type ProductDTO struct {
	ProductID      string
	Name          string
	Description   string
	Category       string
	BasePrice      *big.Rat
	EffectivePrice *big.Rat
	DiscountPercent *int64
	Status         string
}
