package domain

// ChangeTracker tracks which fields were modified for targeted DB updates.
type ChangeTracker struct {
	dirtyFields map[string]bool
}

// NewChangeTracker creates a new change tracker.
func NewChangeTracker() *ChangeTracker {
	return &ChangeTracker{dirtyFields: make(map[string]bool)}
}

// MarkDirty marks a field as modified.
func (c *ChangeTracker) MarkDirty(field string) {
	if c == nil {
		return
	}
	c.dirtyFields[field] = true
}

// Dirty returns true if the field was modified.
func (c *ChangeTracker) Dirty(field string) bool {
	if c == nil {
		return false
	}
	return c.dirtyFields[field]
}

// Reset clears all dirty flags (e.g. after persist).
func (c *ChangeTracker) Reset() {
	if c == nil {
		return
	}
	c.dirtyFields = make(map[string]bool)
}
