package xtime

import "time"

// UTCNowUnixMS returns the current UTC unix timestamp in milliseconds
func UTCNowUnixMS() int64 {
	return time.Now().UTC().UnixMilli()
}

// MSToUTCTime converts a UTC unix milliseconds timestamp to time.Time
func MSToUTCTime(ms int64) time.Time {
	return time.UnixMilli(ms).UTC()
}

// GetWeekStartEnd calculates the start and end dates of a specific week in a given year
// Parameters:
//   - year: The target year
//   - week: The week number (1-53)
// Returns:
//   - time.Time: The start date (Monday) of the specified week
//   - time.Time: The end date (Sunday) of the specified week
func GetWeekStartEnd(year, week int) (time.Time, time.Time) {
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

// GetTotalWeeksInYear returns the total number of weeks in a given year according to ISO 8601
// Parameters:
//   - year: The target year
// Returns:
//   - int: The total number of weeks in the year (52 or 53)
func GetTotalWeeksInYear(year int) int {
	// Get December 28th, which is always in the last week of the year according to ISO 8601
	dec28 := time.Date(year, 12, 28, 0, 0, 0, 0, time.UTC)
	_, week := dec28.ISOWeek()
	return week
}

// AddWeeks adds or subtracts ISO 8601 weeks to/from a given ISO year and week.
func AddWeeks(year, week, deltaWeeks int) (newYear, newWeek int) {
	startDate := isoWeekStartDate(year, week)
	targetDate := startDate.AddDate(0, 0, deltaWeeks*7)
	return targetDate.ISOWeek()
}

// isoWeekStartDate returns the Monday of the given ISO year and week.
func isoWeekStartDate(year, week int) time.Time {
	// ISO week 1 always contains January 4th
	jan4 := time.Date(year, 1, 4, 0, 0, 0, 0, time.UTC)

	// Calculate the Monday of the first ISO week
	weekday := int(jan4.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday fix
	}
	monday := jan4.AddDate(0, 0, -weekday+1)

	// Add (week - 1) weeks
	return monday.AddDate(0, 0, (week-1)*7)
}
