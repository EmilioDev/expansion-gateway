// file: /dto/timeout_tracker.go
package sessions

import (
	"sync"
	"time"
)

// TimeoutTracker is a tiny concurrency-safe model that stores last-activity
// and a timeout duration. It does NOT close sessions; it only reports expiration.
type TimeoutTracker struct {
	mu           sync.RWMutex
	lastActivity time.Time
	timeout      time.Duration
}

// NewTimeoutTracker creates a new tracker with the provided timeout.
func NewTimeoutTracker(timeout time.Duration) *TimeoutTracker {
	return &TimeoutTracker{
		lastActivity: time.Now(),
		timeout:      timeout,
	}
}

// Refresh records activity now.
func (t *TimeoutTracker) Refresh() {
	t.mu.Lock()
	t.lastActivity = time.Now()
	t.mu.Unlock()
}

// Expired reports whether now - lastActivity > timeout.
func (t *TimeoutTracker) Expired() bool {
	t.mu.RLock()
	last := t.lastActivity
	timeout := t.timeout
	t.mu.RUnlock()

	if timeout <= 0 {
		return false
	}
	return time.Since(last) > timeout
}

// SetTimeout changes the timeout duration.
func (t *TimeoutTracker) SetTimeout(d time.Duration) {
	t.mu.Lock()
	t.timeout = d
	t.mu.Unlock()
}

// TimeSinceLastActivity returns the duration since the last activity.
func (t *TimeoutTracker) TimeSinceLastActivity() time.Duration {
	t.mu.RLock()
	last := t.lastActivity
	t.mu.RUnlock()
	return time.Since(last)
}
