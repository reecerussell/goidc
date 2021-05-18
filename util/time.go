package util

import (
	"sync"
	"time"
)

var defaultTime *SystemTime

func init() {
	defaultTime = NewTime()
}

// SystemTime is a wrapper around time.Time, used to aid
// testing by being able to freeze time.
type SystemTime struct {
	mu  sync.Mutex
	now *time.Time
}

// NewTime returns a new, unfrozen instance of SystemTime.
func NewTime() *SystemTime {
	if defaultTime != nil {
		return defaultTime
	}

	return &SystemTime{
		mu: sync.Mutex{},
	}
}

// Time returns the current UTC time and is equivalent to
// time.Now().UTC(). However, if t is frozen, the date
// at which it was frozen will be returned.
func (t *SystemTime) Time() time.Time {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.now != nil {
		return *t.now
	}

	return time.Now().UTC()
}

// Freeze sets t's time to the current UTC time. This
// time will be returned from Time(), until Reset() is called.
func (t *SystemTime) Freeze() {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now().UTC()
	t.now = &now
}

// Reset unfreezes the instance of SystemTime.
func (t *SystemTime) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.now = nil
}

// Time returns the current UTC time and is equivalent to
// time.Now().UTC(). However, if the default SystemTime instance,
// is frozen, the date at which it was frozen will be returned.
func Time() time.Time {
	return defaultTime.Time()
}

// Freeze sets default SystemTime instance to the current UTC time. This
// time will be returned from Time(), until Reset() is called.
func Freeze() {
	defaultTime.Freeze()
}

// Reset unfreezes the default SystemTime.
func Reset() {
	defaultTime.Reset()
}
