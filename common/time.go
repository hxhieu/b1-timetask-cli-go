package common

import "time"

// Week starts
func GetWeekStart(date time.Time, weekStartAt ...time.Weekday) time.Time {
	startAt := time.Monday
	if len(weekStartAt) > 0 {
		startAt = weekStartAt[0]
	}
	offset := (int(startAt) - int(date.Weekday()) - 7) % 7
	result := date.Add(time.Duration(offset*24) * time.Hour)
	return result
}

// From `weekStartAt` to next 7 days
func GetWeekRange(date time.Time, weekStartAt ...time.Weekday) []time.Time {
	result := make([]time.Time, 0)
	weekStart := GetWeekStart(date)
	for i := range 7 {
		result = append(result, weekStart.AddDate(0, 0, i))
	}
	return result
}
