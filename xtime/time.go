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
