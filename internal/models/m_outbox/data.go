package m_outbox

import (
	"time"

	"cloud.google.com/go/spanner"
)

// OutboxEvent is the DB row for outbox_events.
type OutboxEvent struct {
	EventID     string
	EventType   string
	AggregateID string
	Payload     string
	Status      string
	CreatedAt   time.Time
	ProcessedAt *time.Time
}

// ToInsertMut returns an Insert mutation.
func (e *OutboxEvent) ToInsertMut() *spanner.Mutation {
	return spanner.Insert(Table,
		EventID, EventType, AggregateID, Payload, Status, CreatedAt, ProcessedAt,
		e.EventID, e.EventType, e.AggregateID, e.Payload, e.Status, e.CreatedAt, e.ProcessedAt,
	)
}
