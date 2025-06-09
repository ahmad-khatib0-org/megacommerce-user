package utils

import "time"

func TimeGetMillis() int64 {
	return time.Now().UnixMilli()
}
