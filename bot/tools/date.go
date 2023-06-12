package tools

import "time"

func FormatTime(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}

func FormatIntTime(Date int) string {
	return time.Unix(int64(Date), 0).Format("2006-01-02 15:04:05")
}
