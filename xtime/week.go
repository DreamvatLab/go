package xtime

import (
	"time"
)

// GetISOWeekStartEnd calculates the start and end dates of a specific ISO week in a given year
// Parameters:
//   - year: The target year
//   - week: The week number (1-53)
//
// Returns:
//   - time.Time: The start date (Monday) of the specified week
//   - time.Time: The end date (Sunday) of the specified week
func GetISOWeekStartEnd(year, week int) (time.Time, time.Time) {
	// Get the first day of the year
	firstDay := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	// Calculate the weekday of the first day
	weekday := int(firstDay.Weekday())
	if weekday == 0 {
		weekday = 7 // Set Sunday to 7
	}

	// Calculate the first day of the given week (ISO weeks start on Monday)
	startDate := firstDay.AddDate(0, 0, (week-1)*7-(weekday-1))
	endDate := startDate.AddDate(0, 0, 6) // Calculate Sunday

	return startDate, endDate
}

// GetWeekStartEnd calculates the start and end dates of a specific week in a given year,
// with a customizable week start offset from the ISO week start (Monday).
// Parameters:
//   - year: The target year
//   - week: The week number (1-53)
//   - weekStartOffset: The number of days to offset from the ISO week start (Monday).
//     For example:
//   - 0: Week starts on Monday (ISO standard)
//   - 1: Week starts on Tuesday
//   - -1: Week starts on Sunday
//   - 6: Week starts on Sunday (alternative standard)
//
// Returns:
//   - time.Time: The start date of the specified week
//   - time.Time: The end date of the specified week (6 days after start)
//
// Example:
//   - For week 1 of 2024 with offset 0 (Monday start):
//     start: 2024-01-01 (Monday)
//     end: 2024-01-07 (Sunday)
//   - For week 1 of 2024 with offset 6 (Sunday start):
//     start: 2023-12-31 (Sunday)
//     end: 2024-01-06 (Saturday)
func GetWeekStartEnd(year, week, weekStartOffset int) (time.Time, time.Time) {
	start, end := GetISOWeekStartEnd(year, week)
	if weekStartOffset == 0 {
		return start, end
	}

	// Adjust the start date by the offset
	start = start.AddDate(0, 0, weekStartOffset)
	// The end date should be 6 days after the new start date
	end = start.AddDate(0, 0, 6)

	return start, end
}

// GetISOWeeksInYear returns the total number of weeks in a given year according to ISO 8601
// Parameters:
//   - year: The target year
//
// Returns:
//   - int: The total number of weeks in the year (52 or 53)
func GetISOWeeksInYear(year int) int {
	// Get December 28th, which is always in the last week of the year according to ISO 8601
	dec28 := time.Date(year, 12, 28, 0, 0, 0, 0, time.UTC)
	_, week := dec28.ISOWeek()
	return week
}

// AddISOWeeks adds or subtracts ISO 8601 weeks to/from a given ISO year and week.
func AddISOWeeks(year, week, deltaWeeks int) (newYear, newWeek int) {
	startDate := isoWeekStartMonday(year, week)
	targetDate := startDate.AddDate(0, 0, deltaWeeks*7)
	return targetDate.ISOWeek()
}

// isoWeekStartMonday returns the Monday of the specified ISO week and ISO year.
func isoWeekStartMonday(isoYear, isoWeek int) time.Time {
	// January 4th is always in ISO week 1
	referenceDate := time.Date(isoYear, 1, 4, 0, 0, 0, 0, time.UTC)

	// Adjust Go's Weekday (Sunday = 0) so that Monday = 1, Sunday = 7
	weekday := int(referenceDate.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	// Find the Monday of ISO week 1
	firstISOWeekMonday := referenceDate.AddDate(0, 0, -weekday+1)

	// Add (isoWeek - 1) weeks to get the desired ISO week Monday
	targetWeekMonday := firstISOWeekMonday.AddDate(0, 0, (isoWeek-1)*7)

	return targetWeekMonday
}
