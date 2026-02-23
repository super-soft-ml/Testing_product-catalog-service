package domain

import "math/big"

// PricingCalculator is a domain service for price calculations.
// Uses big.Rat for precise decimal arithmetic.
func PricingCalculator_EffectivePrice(basePrice *Money, discount *Discount, percentAsRat *big.Rat) *Money {
	if basePrice == nil {
		return NewMoney(0, 1)
	}
	if discount == nil || percentAsRat == nil {
		return NewMoneyFromRat(basePrice.Rat())
	}
	// discountAmount = basePrice * (percent/100)
	discountAmount := basePrice.Mul(percentAsRat)
	// effective = basePrice - discountAmount
	effective := basePrice.Sub(discountAmount)
	return effective
}

// PercentToRat converts percentage (e.g. 20) to *big.Rat (20/100).
func PercentToRat(percent int64) *big.Rat {
	return big.NewRat(percent, 100)
}
