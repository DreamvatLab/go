package xtime

import (
	"testing"
	"time"
)

func TestUTCNowUnixMS(t *testing.T) {
	// Test that UTCNowUnixMS returns a non-zero value
	ms := UTCNowUnixMS()
	if ms <= 0 {
		t.Errorf("UTCNowUnixMS() returned %d, expected a positive value", ms)
	}

	// Test that the value is close to the current time
	now := time.Now().UTC().UnixMilli()
	diff := now - ms
	if diff < 0 {
		diff = -diff
	}
	if diff > 1000 { // Allow 1 second difference
		t.Errorf("UTCNowUnixMS() returned %d, which is too far from current time %d", ms, now)
	}
}

func TestMSToUTCTime(t *testing.T) {
	// Test with a known timestamp
	ms := int64(1640995200000) // 2022-01-01 00:00:00 UTC
	expected := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	result := MSToUTCTime(ms)
	if !result.Equal(expected) {
		t.Errorf("MSToUTCTime(%d) = %v, expected %v", ms, result, expected)
	}
}
