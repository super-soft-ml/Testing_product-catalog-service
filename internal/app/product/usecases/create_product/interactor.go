package create_product

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"product-catalog-service/internal/app/product/domain"
	"product-catalog-service/internal/pkg/clock"
	"product-catalog-service/internal/pkg/committer"
)

// Request for CreateProduct.
type Request struct {
	Name           string
	Description    string
	Category       string
	BasePriceNum   int64
	BasePriceDenom int64
}

// Interactor creates a new product (draft).
type Interactor struct {
	repo      ProductRepo
	outbox    OutboxRepo
	committer *committer.Committer
	clock     clock.Clock
}

// ProductRepo is the minimal interface for insert.
type ProductRepo interface {
	InsertMut(p *domain.Product) interface{}
}

// OutboxRepo adds outbox mutations.
type OutboxRepo interface {
	InsertMut(eventID, eventType, aggregateID, payload string, now time.Time) interface{}
}

// NewInteractor creates a CreateProduct interactor.
func NewInteractor(repo ProductRepo, outbox OutboxRepo, committer *committer.Committer, clock clock.Clock) *Interactor {
	return &Interactor{repo: repo, outbox: outbox, committer: committer, clock: clock}
}

// Execute creates the product and applies the plan.
func (it *Interactor) Execute(ctx context.Context, req Request) (string, error) {
	basePrice := domain.NewMoney(req.BasePriceNum, req.BasePriceDenom)
	if basePrice == nil {
		basePrice = domain.NewMoney(0, 1)
	}
	productID := uuid.New().String()
	now := it.clock.Now()
	product, err := domain.NewProduct(productID, req.Name, req.Description, req.Category, basePrice, now)
	if err != nil {
		return "", err
	}
	plan := committer.NewPlan()
	if mut := it.repo.InsertMut(product); mut != nil {
		plan.Add(mut)
	}
	for _, ev := range product.DomainEvents() {
		eventID := uuid.New().String()
		payload, _ := json.Marshal(map[string]string{"product_id": productID})
		if m := it.outbox.InsertMut(eventID, ev.EventType(), productID, string(payload), now); m != nil {
			plan.Add(m)
		}
	}
	if err := it.committer.Apply(ctx, plan); err != nil {
		return "", err
	}
	return productID, nil
}
