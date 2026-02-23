package committer

import (
	"context"

	"cloud.google.com/go/spanner"
)

// Plan collects Spanner mutations and applies them in a single read-write transaction.
// This matches the Golden Mutation Pattern: usecases build a plan and apply it once.
type Plan struct {
	mutations []*spanner.Mutation
}

// NewPlan returns a new empty plan.
func NewPlan() *Plan {
	return &Plan{mutations: nil}
}

// Add appends a mutation. Ignores nil.
func (p *Plan) Add(mut interface{}) {
	if mut == nil {
		return
	}
	if m, ok := mut.(*spanner.Mutation); ok {
		p.mutations = append(p.mutations, m)
	}
}

// Mutations returns the collected mutations (read-only).
func (p *Plan) Mutations() []*spanner.Mutation {
	return p.mutations
}

// Committer applies plans to Spanner.
type Committer struct {
	client *spanner.Client
}

// NewCommitter creates a committer for the given Spanner client.
func NewCommitter(client *spanner.Client) *Committer {
	return &Committer{client: client}
}

// Apply runs the plan in a read-write transaction (atomic).
func (c *Committer) Apply(ctx context.Context, plan *Plan) error {
	if plan == nil || len(plan.mutations) == 0 {
		return nil
	}
	_, err := c.client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		return txn.BufferWrite(plan.mutations)
	})
	return err
}
