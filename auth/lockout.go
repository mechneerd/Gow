package auth

import (
	"sync"
	"time"
)

// Lockout tracks failed login attempts and locks users out.
type Lockout struct {
	attempts map[string]int       // key -> attempts
	locked   map[string]time.Time // key -> lock until
	mu       sync.Mutex
	maxAttempts int
	lockoutMinutes int
}

// NewLockout creates a new lockout manager.
func NewLockout(maxAttempts, lockoutMinutes int) *Lockout {
	return &Lockout{
		attempts:       make(map[string]int),
		locked:         make(map[string]time.Time),
		maxAttempts:    maxAttempts,
		lockoutMinutes: lockoutMinutes,
	}
}

// Key can be email, IP, or user ID.
func (l *Lockout) IsLocked(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if until, ok := l.locked[key]; ok {
		if time.Now().Before(until) {
			return true
		}
		delete(l.locked, key)
		delete(l.attempts, key)
	}
	return false
}

func (l *Lockout) RecordFailure(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.attempts[key]++
	if l.attempts[key] >= l.maxAttempts {
		l.locked[key] = time.Now().Add(time.Duration(l.lockoutMinutes) * time.Minute)
	}
}

func (l *Lockout) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.attempts, key)
	delete(l.locked, key)
}

