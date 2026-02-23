package domain

import "math/big"

// Money represents a precise decimal amount (e.g. price) using numerator/denominator.
// Pure value object; no external dependencies.
type Money struct {
	num, denom *big.Int
}

// NewMoney creates Money from numerator and denominator (e.g. 1999, 100 = 19.99).
func NewMoney(numerator, denominator int64) *Money {
	if denominator == 0 {
		denominator = 1
	}
	return &Money{
		num:   big.NewInt(numerator),
		denom: big.NewInt(denominator),
	}
}

// NewMoneyFromRat creates Money from a big.Rat (copy).
func NewMoneyFromRat(r *big.Rat) *Money {
	if r == nil {
		return NewMoney(0, 1)
	}
	return &Money{
		num:   new(big.Int).Set(r.Num()),
		denom: new(big.Int).Set(r.Denom()),
	}
}

// Rat returns a copy of the amount as *big.Rat.
func (m *Money) Rat() *big.Rat {
	if m == nil {
		return new(big.Rat)
	}
	return new(big.Rat).SetFrac(m.num, m.denom)
}

// Num returns the numerator (copy).
func (m *Money) Num() *big.Int {
	if m == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Set(m.num)
}

// Denom returns the denominator (copy).
func (m *Money) Denom() *big.Int {
	if m == nil {
		return big.NewInt(1)
	}
	return new(big.Int).Set(m.denom)
}

// Add returns m + other (new Money).
func (m *Money) Add(other *Money) *Money {
	if m == nil || other == nil {
		return m
	}
	r := m.Rat().Add(m.Rat(), other.Rat())
	return NewMoneyFromRat(r)
}

// Sub returns m - other (new Money).
func (m *Money) Sub(other *Money) *Money {
	if m == nil || other == nil {
		return m
	}
	r := new(big.Rat).Sub(m.Rat(), other.Rat())
	return NewMoneyFromRat(r)
}

// Mul returns m * factor (new Money). factor is e.g. big.NewRat(20, 100) for 20%.
func (m *Money) Mul(factor *big.Rat) *Money {
	if m == nil || factor == nil {
		return m
	}
	r := new(big.Rat).Mul(m.Rat(), factor)
	return NewMoneyFromRat(r)
}

// Cmp returns -1 if m < other, 0 if equal, 1 if m > other.
func (m *Money) Cmp(other *Money) int {
	if m == nil || other == nil {
		return 0
	}
	return m.Rat().Cmp(other.Rat())
}
