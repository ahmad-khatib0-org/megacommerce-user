package utils

import "time"

func TimeGetMillis() int64 {
	return time.Now().UnixMilli()
}

func EmailDateHeader(t time.Time) string {
	return t.Format("2006-01-02 15:04:05 MST")
}
