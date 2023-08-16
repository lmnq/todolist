package utils

import (
	"time"
)

const (
	dateFormat = "2006-01-02"
)

func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse(dateFormat, timeStr)
}

func FormatTime(t time.Time) string {
	return t.Format(dateFormat)
}
