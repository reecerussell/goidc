package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTime_WhereDefaultTimeIsNil_ReturnsNewTime(t *testing.T) {
	dt := defaultTime
	defaultTime = nil

	t.Cleanup(func() {
		defaultTime = dt
	})

	nt := NewTime()
	assert.NotNil(t, nt)
}

func TestNewTime_WhereDefaultTimeIsNotNil_ReturnsDefaultTime(t *testing.T) {
	nt := NewTime()
	assert.Equal(t, defaultTime, nt)
}

func TestFreeze_SetsInternalTime(t *testing.T) {
	t.Cleanup(func() {
		defaultTime.now = nil
	})

	Freeze()

	assert.NotNil(t, defaultTime.now)
}

func TestReset_ClearsInternalTime(t *testing.T) {
	n := time.Now()
	defaultTime.now = &n

	Reset()

	assert.Nil(t, defaultTime.now)
}

func TestTime_WhereTimeIsFrozen_ReturnsFrozenTime(t *testing.T) {
	n := time.Now()
	defaultTime.now = &n

	nt := Time()
	assert.Equal(t, n, nt)
}

func TestTime_WhereTimeIsNotFrozen_ReturnsUtcNow(t *testing.T) {
	defaultTime.now = nil

	nt := Time()
	assert.NotNil(t, nt)
}
