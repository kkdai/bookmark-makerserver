package main

import "time"

// ParseDateTime 將日期時間字符串解析為time.Time類型
func ParseDateTime(dateTimeStr string) (time.Time, error) {
	const layout = "January 02, 2006 at 03:04PM"
	return time.Parse(layout, dateTimeStr)
}

// DaysBefore 計算給定時間`numDays`天之前的時間
func DaysBefore(t time.Time, numDays int) time.Time {
	return t.AddDate(0, 0, -numDays)
}
