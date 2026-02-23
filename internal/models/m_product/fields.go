package m_product

// Table and column names for products table.
const (
	Table                = "products"
	ProductID            = "product_id"
	Name                 = "name"
	Description          = "description"
	Category             = "category"
	BasePriceNumerator   = "base_price_numerator"
	BasePriceDenominator = "base_price_denominator"
	DiscountPercent      = "discount_percent"
	DiscountStartDate    = "discount_start_date"
	DiscountEndDate      = "discount_end_date"
	Status               = "status"
	CreatedAt            = "created_at"
	UpdatedAt            = "updated_at"
	ArchivedAt           = "archived_at"
)

// TableColumns returns column names for full table read.
func TableColumns() []string {
	return []string{
		ProductID, Name, Description, Category,
		BasePriceNumerator, BasePriceDenominator,
		DiscountPercent, DiscountStartDate, DiscountEndDate,
		Status, CreatedAt, UpdatedAt, ArchivedAt,
	}
}
