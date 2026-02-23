package activate_product

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"product-catalog-service/internal/app/product/domain"
	"product-catalog-service/internal/pkg/clock"
	"product-catalog-service/internal/pkg/committer"
)

// Request for ActivateProduct.
type Request struct {
	ProductID string
}

type Interactor struct {
	repo      ProductRepo
	outbox    OutboxRepo
	committer *committer.Committer
	clock     clock.Clock
}

type ProductRepo interface {
	Load(ctx context.Context, productID string) (*domain.Product, error)
	UpdateMut(p *domain.Product) interface{}
}

type OutboxRepo interface {
	InsertMut(eventID, eventType, aggregateID, payload string, now time.Time) interface{}
}

func NewInteractor(repo ProductRepo, outbox OutboxRepo, committer *committer.Committer, clock clock.Clock) *Interactor {
	return &Interactor{repo: repo, outbox: outbox, committer: committer, clock: clock}
}

func (it *Interactor) Execute(ctx context.Context, req Request) error {
	product, err := it.repo.Load(ctx, req.ProductID)
	if err != nil {
		return err
	}
	if err := product.Activate(); err != nil {
		return err
	}
	now := it.clock.Now()
	plan := committer.NewPlan()
	if mut := it.repo.UpdateMut(product); mut != nil {
		plan.Add(mut)
	}
	for _, ev := range product.DomainEvents() {
		eventID := uuid.New().String()
		payload, _ := json.Marshal(map[string]string{"product_id": req.ProductID})
		if m := it.outbox.InsertMut(eventID, ev.EventType(), req.ProductID, string(payload), now); m != nil {
			plan.Add(m)
		}
	}
	return it.committer.Apply(ctx, plan)
}
