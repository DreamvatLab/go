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

func TestGetWeekStartEnd(t *testing.T) {
	tests := []struct {
		name     string
		year     int
		week     int
		expected struct {
			start time.Time
			end   time.Time
		}
	}{
		{
			name: "First week of 2024",
			year: 2024,
			week: 1,
			expected: struct {
				start time.Time
				end   time.Time
			}{
				start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), // Monday
				end:   time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC), // Sunday
			},
		},
		{
			name: "Last week of 2023",
			year: 2023,
			week: 52,
			expected: struct {
				start time.Time
				end   time.Time
			}{
				start: time.Date(2023, 12, 18, 0, 0, 0, 0, time.UTC), // Monday
				end:   time.Date(2023, 12, 24, 0, 0, 0, 0, time.UTC), // Sunday
			},
		},
		{
			name: "Middle week of 2024",
			year: 2024,
			week: 26,
			expected: struct {
				start time.Time
				end   time.Time
			}{
				start: time.Date(2024, 6, 24, 0, 0, 0, 0, time.UTC), // Monday
				end:   time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC), // Sunday
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end := GetWeekStartEnd(tt.year, tt.week)
			if !start.Equal(tt.expected.start) {
				t.Errorf("GetWeekStartEnd(%d, %d) start = %v, expected %v",
					tt.year, tt.week, start, tt.expected.start)
			}
			if !end.Equal(tt.expected.end) {
				t.Errorf("GetWeekStartEnd(%d, %d) end = %v, expected %v",
					tt.year, tt.week, end, tt.expected.end)
			}

			// Verify if the dates are Monday and Sunday respectively
			if start.Weekday() != time.Monday {
				t.Errorf("Start date %v is not a Monday", start)
			}
			if end.Weekday() != time.Sunday {
				t.Errorf("End date %v is not a Sunday", end)
			}

			// Verify if the difference between start and end is 6 days
			if end.Sub(start).Hours() != 144 { // 6 days = 144 hours
				t.Errorf("Difference between start and end should be 6 days, got %v", end.Sub(start))
			}

			// Verify if the dates are in the correct year
			if start.Year() != tt.year && end.Year() != tt.year {
				// Special case: weeks crossing year boundaries may belong to previous or next year
				_, startWeek := start.ISOWeek()
				if startWeek != tt.week {
					t.Errorf("Week number mismatch: got week %d, want week %d", startWeek, tt.week)
				}
			}
		})
	}
}

func TestGetTotalWeeksInYear(t *testing.T) {
	tests := []struct {
		name     string
		year     int
		expected int
	}{
		{
			name:     "2024 (52 weeks)",
			year:     2024,
			expected: 52,
		},
		{
			name:     "2023 (52 weeks)",
			year:     2023,
			expected: 52,
		},
		{
			name:     "2020 (53 weeks)",
			year:     2020,
			expected: 53,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetTotalWeeksInYear(tt.year)
			if result != tt.expected {
				t.Errorf("GetTotalWeeksInYear(%d) = %d, expected %d",
					tt.year, result, tt.expected)
			}
		})
	}
}

func TestAddWeeks(t *testing.T) {
	tests := []struct {
		name       string
		year       int
		week       int
		deltaWeeks int
		wantYear   int
		wantWeek   int
	}{
		{
			name:       "Same year positive delta",
			year:       2024,
			week:       3,
			deltaWeeks: 2,
			wantYear:   2024,
			wantWeek:   5,
		},
		{
			name:       "Same year negative delta",
			year:       2024,
			week:       5,
			deltaWeeks: -2,
			wantYear:   2024,
			wantWeek:   3,
		},
		{
			name:       "Cross year positive",
			year:       2024,
			week:       52,
			deltaWeeks: 2,
			wantYear:   2025,
			wantWeek:   2,
		},
		{
			name:       "Cross year negative",
			year:       2024,
			week:       2,
			deltaWeeks: -2,
			wantYear:   2023,
			wantWeek:   52,
		},
		{
			name:       "Multiple years positive",
			year:       2024,
			week:       52,
			deltaWeeks: 54,
			wantYear:   2026,
			wantWeek:   2,
		},
		{
			name:       "Multiple years negative",
			year:       2024,
			week:       2,
			deltaWeeks: -54,
			wantYear:   2022,
			wantWeek:   52,
		},
		{
			name:       "53-week year - last week",
			year:       2020,
			week:       53,
			deltaWeeks: 0,
			wantYear:   2020,
			wantWeek:   53,
		},
		{
			name:       "53-week year - week 52 to 53",
			year:       2020,
			week:       52,
			deltaWeeks: 1,
			wantYear:   2020,
			wantWeek:   53,
		},
		{
			name:       "53-week year - week 53 to next year",
			year:       2020,
			week:       53,
			deltaWeeks: 1,
			wantYear:   2021,
			wantWeek:   1,
		},
		{
			name:       "53-week year - multiple weeks forward",
			year:       2020,
			week:       52,
			deltaWeeks: 2,
			wantYear:   2021,
			wantWeek:   1,
		},
		{
			name:       "53-week year - multiple weeks backward",
			year:       2021,
			week:       1,
			deltaWeeks: -2,
			wantYear:   2020,
			wantWeek:   52,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotYear, gotWeek := AddWeeks(tt.year, tt.week, tt.deltaWeeks)
			if gotYear != tt.wantYear || gotWeek != tt.wantWeek {
				t.Errorf("AddWeeks(%d, %d, %d) = (%d, %d), want (%d, %d)",
					tt.year, tt.week, tt.deltaWeeks, gotYear, gotWeek, tt.wantYear, tt.wantWeek)
			}
		})
	}
}
