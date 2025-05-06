package xtime

import (
	"testing"
	"time"
)

func TestGetISOWeekStartEnd(t *testing.T) {
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
		{
			name: "Week 53 of 2020",
			year: 2020,
			week: 53,
			expected: struct {
				start time.Time
				end   time.Time
			}{
				start: time.Date(2020, 12, 28, 0, 0, 0, 0, time.UTC), // Monday
				end:   time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),   // Sunday
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end := GetISOWeekStartEnd(tt.year, tt.week)
			if !start.Equal(tt.expected.start) {
				t.Errorf("GetISOWeekStartEnd(%d, %d) start = %v, expected %v",
					tt.year, tt.week, start, tt.expected.start)
			}
			if !end.Equal(tt.expected.end) {
				t.Errorf("GetISOWeekStartEnd(%d, %d) end = %v, expected %v",
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

func TestGetISOWeeksInYear(t *testing.T) {
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
			result := GetISOWeeksInYear(tt.year)
			if result != tt.expected {
				t.Errorf("GetISOWeeksInYear(%d) = %d, expected %d",
					tt.year, result, tt.expected)
			}
		})
	}
}

func TestAddISOWeeks(t *testing.T) {
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
			gotYear, gotWeek := AddISOWeeks(tt.year, tt.week, tt.deltaWeeks)
			if gotYear != tt.wantYear || gotWeek != tt.wantWeek {
				t.Errorf("AddISOWeeks(%d, %d, %d) = (%d, %d), want (%d, %d)",
					tt.year, tt.week, tt.deltaWeeks, gotYear, gotWeek, tt.wantYear, tt.wantWeek)
			}
		})
	}
}

func TestIsoWeekStartMonday(t *testing.T) {
	tests := []struct {
		name     string
		year     int
		week     int
		expected time.Time
	}{
		{
			name:     "First week of 2024",
			year:     2024,
			week:     1,
			expected: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Last week of 2023",
			year:     2023,
			week:     52,
			expected: time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "First week of 2020 (53-week year)",
			year:     2020,
			week:     1,
			expected: time.Date(2019, 12, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Last week of 2020 (53-week year)",
			year:     2020,
			week:     53,
			expected: time.Date(2020, 12, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Middle week of 2024",
			year:     2024,
			week:     26,
			expected: time.Date(2024, 6, 24, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "December 28, 2020 (53-week year)",
			year:     2020,
			week:     53,
			expected: time.Date(2020, 12, 28, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isoWeekStartMonday(tt.year, tt.week)
			if !result.Equal(tt.expected) {
				t.Errorf("isoWeekStartMonday(%d, %d) = %v, expected %v",
					tt.year, tt.week, result, tt.expected)
			}

			// Verify if the result is a Monday
			if result.Weekday() != time.Monday {
				t.Errorf("Result date %v is not a Monday", result)
			}

			// Verify if the ISO week number matches
			_, week := result.ISOWeek()
			if week != tt.week {
				t.Errorf("Week number mismatch: got week %d, want week %d", week, tt.week)
			}
		})
	}
}

func TestGetWeekStartEnd(t *testing.T) {
	tests := []struct {
		name            string
		year            int
		week            int
		weekStartOffset int
		expected        struct {
			start time.Time
			end   time.Time
		}
	}{
		{
			name:            "ISO week (Monday start) - First week of 2024",
			year:            2024,
			week:            1,
			weekStartOffset: 0,
			expected: struct {
				start time.Time
				end   time.Time
			}{
				start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), // Monday
				end:   time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC), // Sunday
			},
		},
		{
			name:            "Sunday start (offset 6) - First week of 2024",
			year:            2024,
			week:            1,
			weekStartOffset: 6,
			expected: struct {
				start time.Time
				end   time.Time
			}{
				start: time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC),  // Next Sunday
				end:   time.Date(2024, 1, 13, 0, 0, 0, 0, time.UTC), // Saturday
			},
		},
		{
			name:            "Tuesday start - First week of 2024",
			year:            2024,
			week:            1,
			weekStartOffset: 1,
			expected: struct {
				start time.Time
				end   time.Time
			}{
				start: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), // Tuesday
				end:   time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC), // Monday
			},
		},
		{
			name:            "ISO week (Monday start) - Week 53 of 2020",
			year:            2020,
			week:            53,
			weekStartOffset: 0,
			expected: struct {
				start time.Time
				end   time.Time
			}{
				start: time.Date(2020, 12, 28, 0, 0, 0, 0, time.UTC), // Monday
				end:   time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),   // Sunday
			},
		},
		{
			name:            "Sunday start (offset 6) - Week 53 of 2020",
			year:            2020,
			week:            53,
			weekStartOffset: 6,
			expected: struct {
				start time.Time
				end   time.Time
			}{
				start: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC), // Next Sunday
				end:   time.Date(2021, 1, 9, 0, 0, 0, 0, time.UTC), // Saturday
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end := GetWeekStartEnd(tt.year, tt.week, tt.weekStartOffset)
			if !start.Equal(tt.expected.start) {
				t.Errorf("GetWeekStartEnd(%d, %d, %d) start = %v, expected %v",
					tt.year, tt.week, tt.weekStartOffset, start, tt.expected.start)
			}
			if !end.Equal(tt.expected.end) {
				t.Errorf("GetWeekStartEnd(%d, %d, %d) end = %v, expected %v",
					tt.year, tt.week, tt.weekStartOffset, end, tt.expected.end)
			}

			// Verify if the difference between start and end is 6 days
			if end.Sub(start).Hours() != 144 { // 6 days = 144 hours
				t.Errorf("Difference between start and end should be 6 days, got %v", end.Sub(start))
			}

			// Verify if the dates are in the correct order
			if end.Before(start) {
				t.Errorf("End date %v is before start date %v", end, start)
			}
		})
	}
}
