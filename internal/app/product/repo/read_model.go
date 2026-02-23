package repo

import (
	"context"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"product-catalog-service/internal/app/product/contracts"
	"product-catalog-service/internal/app/product/domain"
	"product-catalog-service/internal/app/product/domain/services"
	"product-catalog-service/internal/app/product/queries/get_product"
	"product-catalog-service/internal/app/product/queries/list_products"
	m_product "product-catalog-service/internal/models/m_product"
)

// ReadModel implements contracts.ReadModel using Spanner.
type ReadModel struct {
	client *spanner.Client
}

// NewReadModel creates a ReadModel.
func NewReadModel(client *spanner.Client) *ReadModel {
	return &ReadModel{client: client}
}

// GetProduct returns the product with effective price.
func (r *ReadModel) GetProduct(ctx context.Context, productID string) (*get_product.ProductDTO, error) {
	row, err := r.client.Single().ReadRow(ctx, m_product.Table, spanner.Key{productID}, m_product.TableColumns())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}
	var rowData m_product.Product
	if err := row.ToStruct(&rowData); err != nil {
		return nil, err
	}
	return rowToGetDTO(&rowData), nil
}

// ListProducts returns paginated products (default filter: active).
// Total is the total number of matching rows (for pagination), not just the page size.
func (r *ReadModel) ListProducts(ctx context.Context, req list_products.ListProductsRequest) (*list_products.ListProductsResult, error) {
	filterStatus := req.Status
	if filterStatus == "" {
		filterStatus = "active"
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	countStmt := spanner.Statement{
		SQL:    `SELECT COUNT(*) AS total FROM products WHERE status = @status AND archived_at IS NULL`,
		Params: map[string]interface{}{"status": filterStatus},
	}
	if req.Category != "" {
		countStmt.SQL = `SELECT COUNT(*) AS total FROM products WHERE status = @status AND category = @category AND archived_at IS NULL`
		countStmt.Params["category"] = req.Category
	}
	countIter := r.client.Single().Query(ctx, countStmt)
	defer countIter.Stop()
	var total int64
	countRow, err := countIter.Next()
	if err != nil && err != iterator.Done {
		return nil, err
	}
	if err != iterator.Done && countRow != nil {
		if err := countRow.Columns(&total); err != nil {
			return nil, err
		}
	}

	stmt := spanner.Statement{
		SQL: `SELECT product_id, name, description, category, base_price_numerator, base_price_denominator,
		      discount_percent, discount_start_date, discount_end_date, status
		      FROM products WHERE status = @status AND archived_at IS NULL`,
		Params: map[string]interface{}{"status": filterStatus},
	}
	if req.Category != "" {
		stmt.SQL = `SELECT product_id, name, description, category, base_price_numerator, base_price_denominator,
		            discount_percent, discount_start_date, discount_end_date, status
		            FROM products WHERE status = @status AND category = @category AND archived_at IS NULL`
		stmt.Params["category"] = req.Category
	}
	stmt.SQL += " ORDER BY product_id LIMIT @limit OFFSET @offset"
	stmt.Params["limit"] = int64(limit)
	stmt.Params["offset"] = int64(req.Offset)

	iter := r.client.Single().Query(ctx, stmt)
	defer iter.Stop()

	var products []*list_products.ProductDTO
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var p m_product.Product
		if err := row.ToStruct(&p); err != nil {
			return nil, err
		}
		products = append(products, rowToListDTO(&p))
	}
	return &list_products.ListProductsResult{Products: products, Total: int32(total)}, nil
}

func rowToGetDTO(row *m_product.Product) *get_product.ProductDTO {
	basePrice := domain.NewMoney(row.BasePriceNumerator, row.BasePriceDenominator)
	var discount *domain.Discount
	var discountPct *int64
	if row.DiscountPercent != nil && row.DiscountStartDate != nil && row.DiscountEndDate != nil {
		pct := int64(*row.DiscountPercent)
		discountPct = &pct
		discount = domain.NewDiscount(pct, *row.DiscountStartDate, *row.DiscountEndDate)
	}
	effectivePrice := basePrice
	if discount != nil {
		effectivePrice = services.PricingCalculator_EffectivePrice(basePrice, discount, domain.PercentToRat(*discountPct))
	}
	return &get_product.ProductDTO{
		ProductID:       row.ProductID,
		Name:            row.Name,
		Description:     row.Description,
		Category:        row.Category,
		BasePrice:       basePrice.Rat(),
		EffectivePrice:  effectivePrice.Rat(),
		DiscountPercent: discountPct,
		Status:          row.Status,
	}
}

func rowToListDTO(row *m_product.Product) *list_products.ProductDTO {
	basePrice := domain.NewMoney(row.BasePriceNumerator, row.BasePriceDenominator)
	var discount *domain.Discount
	var pct int64
	if row.DiscountPercent != nil && row.DiscountStartDate != nil && row.DiscountEndDate != nil {
		pct = int64(*row.DiscountPercent)
		discount = domain.NewDiscount(pct, *row.DiscountStartDate, *row.DiscountEndDate)
	}
	effectivePrice := services.PricingCalculator_EffectivePrice(basePrice, discount, domain.PercentToRat(pct))
	if effectivePrice == nil {
		effectivePrice = basePrice
	}
	return &list_products.ProductDTO{
		ProductID:       row.ProductID,
		Name:           row.Name,
		Description:    row.Description,
		Category:       row.Category,
		BasePrice:      basePrice.Rat(),
		EffectivePrice: effectivePrice.Rat(),
		Status:         row.Status,
	}
}

var _ contracts.ReadModel = (*ReadModel)(nil)
