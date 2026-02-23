package m_outbox

// Table and column names for outbox_events table.
const (
	Table        = "outbox_events"
	EventID      = "event_id"
	EventType    = "event_type"
	AggregateID  = "aggregate_id"
	Payload      = "payload"
	Status       = "status"
	CreatedAt    = "created_at"
	ProcessedAt  = "processed_at"
)

const StatusPending = "pending"
