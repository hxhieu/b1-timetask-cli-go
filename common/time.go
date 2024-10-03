package common

import "time"

// Week starts
func GetWeekStart(date time.Time) time.Time {
	startAt := time.Monday
	offset := (int(startAt) - int(date.Weekday()) - 7) % 7
	result := date.Add(time.Duration(offset*24) * time.Hour)
	return result
}

// From `date` to next 7 days
func GetWeekRange(date time.Time) []time.Time {
	result := make([]time.Time, 0)
	weekStart := GetWeekStart(date)
	for i := range 7 {
		result = append(result, weekStart.AddDate(0, 0, i))
	}
	return result
}

func DateToString(date time.Time) string {
	return date.Format("2006-01-02")
}
