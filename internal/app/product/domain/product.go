package domain

import "time"

// ProductStatus represents lifecycle state of a product.
type ProductStatus string

const (
	ProductStatusDraft     ProductStatus = "draft"
	ProductStatusActive    ProductStatus = "active"
	ProductStatusInactive  ProductStatus = "inactive"
	ProductStatusArchived  ProductStatus = "archived"
)

// Product is the product aggregate. All business rules are enforced here.
// No context, no DB, no proto - pure domain.
type Product struct {
	id          string
	name        string
	description string
	category    string
	basePrice   *Money
	discount    *Discount
	status      ProductStatus
	archivedAt  *time.Time
	changes     *ChangeTracker
	events      []DomainEvent
}

// NewProduct creates a new product (draft). Used when creating.
func NewProduct(id, name, description, category string, basePrice *Money, now time.Time) (*Product, error) {
	if name == "" || category == "" || basePrice == nil {
		return nil, ErrInvalidInput
	}
	if basePrice.Cmp(NewMoney(0, 1)) < 0 {
		return nil, ErrInvalidPrice
	}
	p := &Product{
		id:          id,
		name:        name,
		description: description,
		category:    category,
		basePrice:   basePrice,
		status:      ProductStatusDraft,
		changes:     NewChangeTracker(),
		events:      nil,
	}
	// Mark all as dirty for insert
	for _, f := range []string{FieldName, FieldDescription, FieldCategory, FieldBasePrice, FieldStatus} {
		p.changes.MarkDirty(f)
	}
	p.events = append(p.events, &ProductCreatedEvent{ProductID: id})
	return p, nil
}

// ReconstituteProduct rebuilds aggregate from persistence (for updates).
// Does not emit events; used when loading from DB.
func ReconstituteProduct(id, name, description, category string, basePrice *Money, discount *Discount, status ProductStatus, archivedAt *time.Time) *Product {
	p := &Product{
		id:          id,
		name:        name,
		description: description,
		category:    category,
		basePrice:   basePrice,
		discount:    discount,
		status:      status,
		archivedAt:  archivedAt,
		changes:     NewChangeTracker(),
		events:      nil,
	}
	return p
}

// ID returns the product ID.
func (p *Product) ID() string { return p.id }

// Name returns the product name.
func (p *Product) Name() string { return p.name }

// Description returns the product description.
func (p *Product) Description() string { return p.description }

// Category returns the product category.
func (p *Product) Category() string { return p.category }

// BasePrice returns the base price (copy as Money).
func (p *Product) BasePrice() *Money {
	if p.basePrice == nil {
		return nil
	}
	return NewMoneyFromRat(p.basePrice.Rat())
}

// Discount returns the current discount if any.
func (p *Product) Discount() *Discount { return p.discount }

// Status returns the product status.
func (p *Product) Status() ProductStatus { return p.status }

// ArchivedAt returns when the product was archived, or nil.
func (p *Product) ArchivedAt() *time.Time { return p.archivedAt }

// Changes returns the change tracker (for repo to build mutations).
func (p *Product) Changes() *ChangeTracker { return p.changes }

// DomainEvents returns and clears the collected domain events (consume once).
func (p *Product) DomainEvents() []DomainEvent {
	ev := p.events
	p.events = nil
	return ev
}

// Update updates name, description, category. Only draft or active can be updated.
func (p *Product) Update(name, description, category string) error {
	if p.status != ProductStatusDraft && p.status != ProductStatusActive {
		return ErrProductNotActive
	}
	if p.status == ProductStatusArchived {
		return ErrProductAlreadyArchived
	}
	if name != "" {
		p.name = name
		p.changes.MarkDirty(FieldName)
	}
	if description != p.description {
		p.description = description
		p.changes.MarkDirty(FieldDescription)
	}
	if category != "" {
		p.category = category
		p.changes.MarkDirty(FieldCategory)
	}
	p.events = append(p.events, &ProductUpdatedEvent{ProductID: p.id})
	return nil
}

// Activate sets status to active. Only draft can be activated.
func (p *Product) Activate() error {
	if p.status == ProductStatusActive {
		return ErrProductAlreadyActive
	}
	if p.status == ProductStatusArchived {
		return ErrProductAlreadyArchived
	}
	p.status = ProductStatusActive
	p.changes.MarkDirty(FieldStatus)
	p.events = append(p.events, &ProductActivatedEvent{ProductID: p.id})
	return nil
}

// Deactivate sets status to inactive.
func (p *Product) Deactivate() error {
	if p.status == ProductStatusArchived {
		return ErrProductAlreadyArchived
	}
	p.status = ProductStatusInactive
	p.changes.MarkDirty(FieldStatus)
	p.events = append(p.events, &ProductDeactivatedEvent{ProductID: p.id})
	return nil
}

// ApplyDiscount sets the product discount. Only one active discount; only active product.
func (p *Product) ApplyDiscount(discount *Discount, now time.Time) error {
	if p.status != ProductStatusActive {
		return ErrProductNotActive
	}
	if discount == nil || !discount.IsValidAt(now) {
		return ErrInvalidDiscountPeriod
	}
	if p.discount != nil && p.discount.IsValidAt(now) {
		return ErrDiscountAlreadyActive
	}
	p.discount = discount
	p.changes.MarkDirty(FieldDiscount)
	p.events = append(p.events, &DiscountAppliedEvent{ProductID: p.id})
	return nil
}

// RemoveDiscount clears the discount.
func (p *Product) RemoveDiscount() error {
	if p.status == ProductStatusArchived {
		return ErrProductAlreadyArchived
	}
	p.discount = nil
	p.changes.MarkDirty(FieldDiscount)
	p.events = append(p.events, &DiscountRemovedEvent{ProductID: p.id})
	return nil
}

// Archive soft-deletes the product.
func (p *Product) Archive(now time.Time) error {
	if p.status == ProductStatusArchived {
		return ErrProductAlreadyArchived
	}
	p.status = ProductStatusArchived
	p.archivedAt = &now
	p.changes.MarkDirty(FieldStatus)
	p.changes.MarkDirty(FieldArchivedAt)
	p.events = append(p.events, &ProductArchivedEvent{ProductID: p.id})
	return nil
}
