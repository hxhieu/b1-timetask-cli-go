package common

import (
	"fmt"
	"math/rand/v2"
	"time"
)

// Week starts
func getWeekStart(date time.Time) time.Time {
	startAt := time.Monday
	offset := (int(startAt) - int(date.Weekday()) - 7) % 7
	result := date.Add(time.Duration(offset*24) * time.Hour)
	return result
}

// From `date` as start of week to next 7 days
func GetWeekRange(date time.Time, weekOffset int) []time.Time {
	// Add week offset
	actualDate := date.AddDate(0, 0, 7*weekOffset)
	result := make([]time.Time, 0)
	weekStart := getWeekStart(actualDate)
	for i := range 7 {
		dateTime := weekStart.AddDate(0, 0, i)
		// Convert to local date at start of day
		localDate := time.Date(
			dateTime.Year(), dateTime.Month(), dateTime.Day(),
			0, 0, 0, 0, dateTime.Location(),
		)
		result = append(result, localDate)
	}
	return result
}

func DateToString(date time.Time) string {
	return date.Format("2006-01-02")
}

func RandomDelay(min, max time.Duration) time.Duration {
	if max <= min {
		return min
	}
	diff := max - min
	return min + time.Duration(rand.Int64N(int64(diff)+1))
}

// ParseGraphDateTime parses a datetime string from Microsoft Graph API using multiple possible formats
func ParseGraphDateTime(dateTimeStr string) (time.Time, error) {
	// Try multiple time formats that Microsoft Graph might use
	timeFormats := []string{
		time.RFC3339,                  // "2006-01-02T15:04:05Z07:00"
		"2006-01-02T15:04:05.0000000", // Microsoft Graph format without timezone
		"2006-01-02T15:04:05",         // Basic ISO format
		time.RFC3339Nano,              // With nanoseconds
	}

	for _, format := range timeFormats {
		if parsedTime, err := time.Parse(format, dateTimeStr); err == nil {
			return parsedTime, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse datetime string: %s", dateTimeStr)
}
