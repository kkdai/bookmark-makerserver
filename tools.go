package main

import (
	"strings"
	"time"
)

// ParseDateTime 將日期時間字符串解析為time.Time類型
func ParseDateTime(dateTimeStr string) (time.Time, error) {
	const layout = "January 02, 2006 at 03:04PM"
	return time.Parse(layout, dateTimeStr)
}

// DaysBefore 計算給定時間`numDays`天之前的時間
func DaysBefore(t time.Time, numDays int) time.Time {
	return t.AddDate(0, 0, -numDays)
}

// removeFirstAndLastLine takes a string and removes the first and last lines.
func removeFirstAndLastLine(s string) string {
	// Split the string into lines.
	lines := strings.Split(s, "\n")

	// If there are less than 3 lines, return an empty string because removing the first and last would leave nothing.
	if len(lines) < 3 {
		return ""
	}

	// Join the lines back together, skipping the first and last lines.
	return strings.Join(lines[1:len(lines)-1], "\n")
}
