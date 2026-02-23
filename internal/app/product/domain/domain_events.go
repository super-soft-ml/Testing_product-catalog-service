package domain

// DomainEvent is a marker for events emitted by the domain.
type DomainEvent interface {
	EventType() string
}

// ProductCreatedEvent is emitted when a product is created.
type ProductCreatedEvent struct {
	ProductID string
}

func (e *ProductCreatedEvent) EventType() string { return "product.created" }

// ProductUpdatedEvent is emitted when product details change.
type ProductUpdatedEvent struct {
	ProductID string
}

func (e *ProductUpdatedEvent) EventType() string { return "product.updated" }

// ProductActivatedEvent is emitted when a product is activated.
type ProductActivatedEvent struct {
	ProductID string
}

func (e *ProductActivatedEvent) EventType() string { return "product.activated" }

// ProductDeactivatedEvent is emitted when a product is deactivated.
type ProductDeactivatedEvent struct {
	ProductID string
}

func (e *ProductDeactivatedEvent) EventType() string { return "product.deactivated" }

// DiscountAppliedEvent is emitted when a discount is applied.
type DiscountAppliedEvent struct {
	ProductID string
}

func (e *DiscountAppliedEvent) EventType() string { return "discount.applied" }

// DiscountRemovedEvent is emitted when a discount is removed.
type DiscountRemovedEvent struct {
	ProductID string
}

func (e *DiscountRemovedEvent) EventType() string { return "discount.removed" }

// ProductArchivedEvent is emitted when a product is archived (soft delete).
type ProductArchivedEvent struct {
	ProductID string
}

func (e *ProductArchivedEvent) EventType() string { return "product.archived" }
