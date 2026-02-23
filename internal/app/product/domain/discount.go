package domain

import "time"

// Discount is a percentage-based discount valid between start and end date.
// Pure value object.
type Discount struct {
	percent    int64     // e.g. 20 for 20%
	startDate  time.Time
	endDate    time.Time
}

// NewDiscount creates a discount. percent is 0-100.
func NewDiscount(percent int64, startDate, endDate time.Time) *Discount {
	return &Discount{
		percent:   percent,
		startDate: startDate,
		endDate:   endDate,
	}
}

// Percentage returns the discount percentage (0-100).
func (d *Discount) Percentage() int64 {
	if d == nil {
		return 0
	}
	return d.percent
}

// StartDate returns the start of the discount period.
func (d *Discount) StartDate() time.Time {
	if d == nil {
		return time.Time{}
	}
	return d.startDate
}

// EndDate returns the end of the discount period.
func (d *Discount) EndDate() time.Time {
	if d == nil {
		return time.Time{}
	}
	return d.endDate
}

// IsValidAt returns true if the discount is active at the given time.
func (d *Discount) IsValidAt(t time.Time) bool {
	if d == nil {
		return false
	}
	return !t.Before(d.startDate) && !t.After(d.endDate)
}
