package domain

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMoney_AddSubMul(t *testing.T) {
	a := NewMoney(1999, 100) // 19.99
	b := NewMoney(500, 100)   // 5.00

	sum := a.Add(b)
	assert.Equal(t, "24.99", sum.Rat().FloatString(2))

	diff := a.Sub(b)
	assert.Equal(t, "14.99", diff.Rat().FloatString(2))

	twentyPercent := big.NewRat(20, 100)
	discount := a.Mul(twentyPercent)
	assert.Equal(t, "4.00", discount.Rat().FloatString(2))
}

func TestPricingCalculator_EffectivePrice(t *testing.T) {
	base := NewMoney(5000, 100) // 50.00
	discount := NewDiscount(20, time.Now().Add(-time.Hour), time.Now().Add(time.Hour))

	effective := PricingCalculator_EffectivePrice(base, discount, PercentToRat(20))
	require.NotNil(t, effective)
	// 20% off 50 = 40
	assert.Equal(t, "40.00", effective.Rat().FloatString(2))
}

func TestPricingCalculator_NoDiscount(t *testing.T) {
	base := NewMoney(1999, 100)
	effective := PricingCalculator_EffectivePrice(base, nil, nil)
	require.NotNil(t, effective)
	assert.Equal(t, base.Rat().Cmp(effective.Rat()), 0)
}
