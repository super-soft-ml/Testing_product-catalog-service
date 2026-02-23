package clock

import "time"

// Clock abstracts time for testing.
type Clock interface {
	Now() time.Time
}

// RealClock uses system time.
type RealClock struct{}

func (RealClock) Now() time.Time { return time.Now() }
