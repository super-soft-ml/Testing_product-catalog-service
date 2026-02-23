package repo

import (
	"context"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"product-catalog-service/internal/app/product/contracts"
	"product-catalog-service/internal/app/product/domain"
	"product-catalog-service/internal/models/m_outbox"
	m_product "product-catalog-service/internal/models/m_product"
)

// OutboxRepo writes outbox events as mutations (same transaction as product).
type OutboxRepo struct{}

// InsertMut returns a mutation to insert one outbox event. Payload is JSON string.
func (r *OutboxRepo) InsertMut(eventID, eventType, aggregateID, payload string, now time.Time) *spanner.Mutation {
	e := &m_outbox.OutboxEvent{
		EventID:     eventID,
		EventType:   eventType,
		AggregateID: aggregateID,
		Payload:     payload,
		Status:      m_outbox.StatusPending,
		CreatedAt:   now,
		ProcessedAt: nil,
	}
	return e.ToInsertMut()
}

// ProductRepo implements contracts.ProductRepo with Spanner.
type ProductRepo struct {
	client *spanner.Client
}

// NewProductRepo returns a new ProductRepo.
func NewProductRepo(client *spanner.Client) *ProductRepo {
	return &ProductRepo{client: client}
}

// Load reads the product by ID and reconstitutes the aggregate.
func (r *ProductRepo) Load(ctx context.Context, productID string) (*domain.Product, error) {
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
	return rowToDomain(&rowData), nil
}

func rowToDomain(row *m_product.Product) *domain.Product {
	if row == nil {
		return nil
	}
	basePrice := domain.NewMoney(row.BasePriceNumerator, row.BasePriceDenominator)
	var disc *domain.Discount
	if row.DiscountPercent != nil && row.DiscountStartDate != nil && row.DiscountEndDate != nil {
		disc = domain.NewDiscount(int64(*row.DiscountPercent), *row.DiscountStartDate, *row.DiscountEndDate)
	}
	return domain.ReconstituteProduct(
		row.ProductID, row.Name, row.Description, row.Category,
		basePrice, disc, domain.ProductStatus(row.Status), row.ArchivedAt,
	)
}

// InsertMut returns a mutation to insert the product, or nil if invalid.
func (r *ProductRepo) InsertMut(p *domain.Product) interface{} {
	if p == nil {
		return nil
	}
	row := productToRow(p)
	return row.ToInsertMut()
}

// UpdateMut returns a mutation for updated fields only, or nil if nothing dirty.
func (r *ProductRepo) UpdateMut(p *domain.Product) interface{} {
	if p == nil {
		return nil
	}
	updates := make(map[string]interface{})

	if p.Changes().Dirty(domain.FieldName) {
		updates[m_product.Name] = p.Name()
	}
	if p.Changes().Dirty(domain.FieldDescription) {
		updates[m_product.Description] = p.Description()
	}
	if p.Changes().Dirty(domain.FieldCategory) {
		updates[m_product.Category] = p.Category()
	}
	if p.Changes().Dirty(domain.FieldBasePrice) {
		base := p.BasePrice()
		if base != nil {
			updates[m_product.BasePriceNumerator] = base.Num().Int64()
			updates[m_product.BasePriceDenominator] = base.Denom().Int64()
		}
	}
	if p.Changes().Dirty(domain.FieldDiscount) {
		if d := p.Discount(); d != nil {
			pct := float64(d.Percentage())
			updates[m_product.DiscountPercent] = &pct
			st, end := d.StartDate(), d.EndDate()
			updates[m_product.DiscountStartDate] = &st
			updates[m_product.DiscountEndDate] = &end
		} else {
			updates[m_product.DiscountPercent] = nil
			updates[m_product.DiscountStartDate] = nil
			updates[m_product.DiscountEndDate] = nil
		}
	}
	if p.Changes().Dirty(domain.FieldStatus) {
		updates[m_product.Status] = string(p.Status())
	}
	if p.Changes().Dirty(domain.FieldArchivedAt) {
		updates[m_product.ArchivedAt] = p.ArchivedAt()
	}

	if len(updates) == 0 {
		return nil
	}
	updates[m_product.UpdatedAt] = time.Now()
	return m_product.UpdateMut(p.ID(), updates)
}

func productToRow(p *domain.Product) *m_product.Product {
	base := p.BasePrice()
	num, denom := int64(0), int64(1)
	if base != nil {
		num = base.Num().Int64()
		denom = base.Denom().Int64()
		if denom == 0 {
			denom = 1
		}
	}
	row := &m_product.Product{
		ProductID:            p.ID(),
		Name:                 p.Name(),
		Description:          p.Description(),
		Category:             p.Category(),
		BasePriceNumerator:   num,
		BasePriceDenominator: denom,
		Status:               string(p.Status()),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
		ArchivedAt:           p.ArchivedAt(),
	}
	if d := p.Discount(); d != nil {
		pct := float64(d.Percentage())
		st, end := d.StartDate(), d.EndDate()
		row.DiscountPercent = &pct
		row.DiscountStartDate = &st
		row.DiscountEndDate = &end
	}
	return row
}

var _ contracts.ProductRepo = (*ProductRepo)(nil)
