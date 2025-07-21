package utils

import "time"

func TimeGetMillis() int64 {
	return time.Now().UnixMilli()
}

func TimeGetMillisFromTime(t time.Time) int64 {
	return t.UnixMilli()
}

func EmailDateHeader(t time.Time) string {
	return t.Format("2006-01-02 15:04:05 MST")
}
