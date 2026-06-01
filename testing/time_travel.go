package testing

import (
	"sync"
	"time"
)

// TimeTravel provides time manipulation for testing.
type TimeTravel struct {
	now     time.Time
	running bool
	mu      sync.Mutex
}

// NewTimeTravel creates a new TimeTravel instance.
func NewTimeTravel() *TimeTravel {
	return &TimeTravel{
		now: time.Now(),
	}
}

// Now returns the current "fake" time.
func (t *TimeTravel) Now() time.Time {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.running {
		return t.now
	}
	return time.Now()
}

// SetNow sets the current fake time.
func (t *TimeTravel) SetNow(now time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.now = now
	t.running = true
}

// Travel moves the fake time by the given duration.
func (t *TimeTravel) Travel(d time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.now = t.now.Add(d)
	t.running = true
}

// TravelToDate moves the fake time to a specific date.
func (t *TimeTravel) TravelToDate(year int, month time.Month, day, hour, min, sec int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.now = time.Date(year, month, day, hour, min, sec, 0, time.Local)
	t.running = true
}

// Freeze freezes the current time.
func (t *TimeTravel) Freeze() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.now = time.Now()
	t.running = true
}

// Unfreeze stops time travel and returns to real time.
func (t *TimeTravel) Unfreeze() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.running = false
}

// IsFrozen returns whether time travel is active.
func (t *TimeTravel) IsFrozen() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.running
}

// Add adds duration to the current fake time.
func (t *TimeTravel) Add(d time.Duration) {
	t.Travel(d)
}

// Sub returns the difference between two fake times.
func (t *TimeTravel) Sub(other time.Time) time.Duration {
	return t.Now().Sub(other)
}

// Since returns the duration since the given time (using fake time).
func (t *TimeTravel) Since(then time.Time) time.Duration {
	return t.Now().Sub(then)
}

// Until returns the duration until the given time (using fake time).
func (t *TimeTravel) Until(then time.Time) time.Duration {
	return then.Sub(t.Now())
}

// Global time travel instance for convenience
var globalTimeTravel *TimeTravel
var timeTravelOnce sync.Once

// GetTimeTravel returns the global TimeTravel instance.
func GetTimeTravel() *TimeTravel {
	timeTravelOnce.Do(func() {
		globalTimeTravel = NewTimeTravel()
	})
	return globalTimeTravel
}

// FreezeTime freezes the global time.
func FreezeTime() *TimeTravel {
	tt := GetTimeTravel()
	tt.Freeze()
	return tt
}

// TravelTime travels by the given duration using the global time.
func TravelTime(d time.Duration) *TimeTravel {
	tt := GetTimeTravel()
	tt.Travel(d)
	return tt
}

// UnfreezeTime stops global time travel.
func UnfreezeTime() {
	tt := GetTimeTravel()
	tt.Unfreeze()
}