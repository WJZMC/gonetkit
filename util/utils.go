package util

import "time"

func GetTimer() int64 {
	return time.Now().UnixNano()
}
