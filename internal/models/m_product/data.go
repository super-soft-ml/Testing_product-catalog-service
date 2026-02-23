package m_product

import (
	"time"

	"cloud.google.com/go/spanner"
)

// Product is the DB row representation for products table.
type Product struct {
	ProductID            string
	Name                 string
	Description          string
	Category             string
	BasePriceNumerator   int64
	BasePriceDenominator int64
	DiscountPercent      *float64
	DiscountStartDate    *time.Time
	DiscountEndDate      *time.Time
	Status               string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	ArchivedAt           *time.Time
}

// ToSpannerMutation returns an Insert mutation for the row.
func (p *Product) ToInsertMut() *spanner.Mutation {
	return spanner.Insert(Table,
		ProductID, Name, Description, Category,
		BasePriceNumerator, BasePriceDenominator,
		DiscountPercent, DiscountStartDate, DiscountEndDate,
		Status, CreatedAt, UpdatedAt, ArchivedAt,
		p.ProductID, p.Name, p.Description, p.Category,
		p.BasePriceNumerator, p.BasePriceDenominator,
		p.DiscountPercent, p.DiscountStartDate, p.DiscountEndDate,
		p.Status, p.CreatedAt, p.UpdatedAt, p.ArchivedAt,
	)
}

// UpdateMut builds an Update mutation for the given key and column values.
// updates must not include ProductID; it is added as the key.
func UpdateMut(productID string, updates map[string]interface{}) *spanner.Mutation {
	// Build ordered columns and values: key first, then UpdatedAt, then rest
	cols := []string{ProductID}
	vals := []interface{}{productID}
	order := []string{Name, Description, Category, BasePriceNumerator, BasePriceDenominator,
		DiscountPercent, DiscountStartDate, DiscountEndDate, Status, CreatedAt, UpdatedAt, ArchivedAt}
	for _, k := range order {
		if v, ok := updates[k]; ok {
			cols = append(cols, k)
			vals = append(vals, v)
		}
	}
	return spanner.Update(Table, cols, vals)
}
