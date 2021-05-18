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

func NewTime() *SystemTime {
	if defaultTime != nil {
		return defaultTime
	}

	return &SystemTime{
		mu: sync.Mutex{},
	}
}

func (t *SystemTime) Time() time.Time {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.now != nil {
		return *t.now
	}

	return time.Now().UTC()
}

func (t *SystemTime) Freeze() {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now().UTC()
	t.now = &now
}

func (t *SystemTime) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.now = nil
}

func Time() time.Time {
	return defaultTime.Time()
}

func Freeze() {
	defaultTime.Freeze()
}

func Reset() {
	defaultTime.Reset()
}
